package responsewriters

import (
  "compress/gzip"
  "encoding/json"
  "fmt"
  "github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/negotiation"
  "github.com/sirupsen/logrus"
  "io"
  "k8s.io/apimachinery/pkg/runtime"
  "k8s.io/apimachinery/pkg/runtime/schema"
  "net/http"
  "strconv"
  "strings"
  "sync"
)

func ErrorNegotiated(err error, s runtime.NegotiatedSerializer, gv schema.GroupVersion, w http.ResponseWriter, req *http.Request) int {
  status := ErrorToAPIStatus(err)
  code := int(status.Code)
  // when writing an error, check to see if the status indicates a retry after period
  if status.Details != nil && status.Details.RetryAfterSeconds > 0 {
    delay := strconv.Itoa(int(status.Details.RetryAfterSeconds))
    w.Header().Set("Retry-After", delay)
  }

  if code == http.StatusNoContent {
    w.WriteHeader(code)
    return code
  }

  WriteObjectNegotiated(s, negotiation.DefaultEndpointRestrictions, gv, w, req, code, status)
  return code
}

// WriteObjectNegotiated renders an object in the content type negotiated by the client.
func WriteObjectNegotiated(s runtime.NegotiatedSerializer, restrictions negotiation.EndpointRestrictions, gv schema.GroupVersion, w http.ResponseWriter, req *http.Request, statusCode int, object runtime.Object) {
  _, serializer, err := negotiation.NegotiateOutputMediaType(req, s, restrictions)
  if err != nil {
    // if original statusCode was not successful we need to return the original error
    // we cannot hide it behind negotiation problems
    if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
      WriteRawJSON(int(statusCode), object, w)
      return
    }
    status := ErrorToAPIStatus(err)
    WriteRawJSON(int(status.Code), status, w)
    return
  }

  //if ae := request.AuditEventFrom(req.Context()); ae != nil {
  //  audit.LogResponseObject(ae, object, gv, s)
  //}

  encoder := s.EncoderForVersion(serializer.Serializer, gv)
  SerializeObject(serializer.MediaType, encoder, w, req, statusCode, object)
}

// SerializeObject renders an object in the content type negotiated by the client using the provided encoder.
// The context is optional and can be nil. This method will perform optional content compression if requested by
// a client and the feature gate for APIResponseCompression is enabled.
func SerializeObject(mediaType string, encoder runtime.Encoder, hw http.ResponseWriter, req *http.Request, statusCode int, object runtime.Object) {
  w := &deferredResponseWriter{
    mediaType:       mediaType,
    statusCode:      statusCode,
    contentEncoding: negotiateContentEncoding(req),
    hw:              hw,
  }

  err := encoder.Encode(object, w)
  if err == nil {
    err = w.Close()
    if err != nil {
      logrus.Errorf("apiserver was unable to close cleanly the response writer: %v", err)
    }
    return
  }

  // make a best effort to write the object if a failure is detected
  logrus.Errorf("apiserver was unable to write a JSON response: %v", err)
  status := ErrorToAPIStatus(err)
  candidateStatusCode := int(status.Code)
  // if the current status code is successful, allow the error's status code to overwrite it
  if statusCode >= http.StatusOK && statusCode < http.StatusBadRequest {
    w.statusCode = candidateStatusCode
  }
  output, err := runtime.Encode(encoder, status)
  if err != nil {
    w.mediaType = "text/plain"
    output = []byte(fmt.Sprintf("%s: %s", status.Reason, status.Message))
  }
  if _, err := w.Write(output); err != nil {
    logrus.Errorf("apiserver was unable to write a fallback JSON response: %v", err)
  }
  w.Close()
}

// negotiateContentEncoding returns a supported client-requested content encoding for the
// provided request. It will return the empty string if no supported content encoding was
// found or if response compression is disabled.
func negotiateContentEncoding(req *http.Request) string {
  encoding := req.Header.Get("Accept-Encoding")
  if len(encoding) == 0 {
    return ""
  }
  //if !utilfeature.DefaultFeatureGate.Enabled(features.APIResponseCompression) {
  //  return ""
  //}
  for len(encoding) > 0 {
    var token string
    if next := strings.Index(encoding, ","); next != -1 {
      token = encoding[:next]
      encoding = encoding[next+1:]
    } else {
      token = encoding
      encoding = ""
    }
    switch strings.TrimSpace(token) {
    case "gzip":
      return "gzip"
    }
  }
  return ""
}

var gzipPool = &sync.Pool{
  New: func() interface{} {
    gw, err := gzip.NewWriterLevel(nil, defaultGzipContentEncodingLevel)
    if err != nil {
      panic(err)
    }
    return gw
  },
}

const (
  // defaultGzipContentEncodingLevel is set to 4 which uses less CPU than the default level
  defaultGzipContentEncodingLevel = 4
  // defaultGzipThresholdBytes is compared to the size of the first write from the stream
  // (usually the entire object), and if the size is smaller no gzipping will be performed
  // if the client requests it.
  defaultGzipThresholdBytes = 128 * 1024
)

type deferredResponseWriter struct {
  mediaType       string
  statusCode      int
  contentEncoding string

  hasWritten bool
  hw         http.ResponseWriter
  w          io.Writer
}

func (w *deferredResponseWriter) Write(p []byte) (n int, err error) {
  if w.hasWritten {
    return w.w.Write(p)
  }
  w.hasWritten = true

  hw := w.hw
  header := hw.Header()
  switch {
  case w.contentEncoding == "gzip" && len(p) > defaultGzipThresholdBytes:
    header.Set("Content-Encoding", "gzip")
    header.Add("Vary", "Accept-Encoding")

    gw := gzipPool.Get().(*gzip.Writer)
    gw.Reset(hw)

    w.w = gw
  default:
    w.w = hw
  }

  header.Set("Content-Type", w.mediaType)
  hw.WriteHeader(w.statusCode)
  return w.w.Write(p)
}

func (w *deferredResponseWriter) Close() error {
  if !w.hasWritten {
    return nil
  }
  var err error
  switch t := w.w.(type) {
  case *gzip.Writer:
    err = t.Close()
    t.Reset(nil)
    gzipPool.Put(t)
  }
  return err
}

// WriteRawJSON writes a non-API object in JSON.
func WriteRawJSON(statusCode int, object interface{}, w http.ResponseWriter) {
  output, err := json.MarshalIndent(object, "", "  ")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(statusCode)
  w.Write(output)
}
