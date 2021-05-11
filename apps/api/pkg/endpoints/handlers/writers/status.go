package writers

import (
	"encoding/json"
	"fmt"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/negotiation"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/storage"
	"github.com/sirupsen/logrus"

	"net/http"
)

type statusError interface {
	Status() metav1.Status
}

func ErrorNegotiated(err error, s runtime.NegotiatedSerializer, gv schema.GroupVersion, w http.ResponseWriter, req *http.Request) int {
	status := ErrorToApiStatus(err)
	code := int(status.Code)

	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return code
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	WriteObjectNegotiated(s, negotiation.DefaultEndpointRestrictions, gv, w, req, code, status)
	return code
}

func ErrorToApiStatus(err error) *metav1.Status {
	switch t := err.(type) {
	case statusError:
		status := t.Status()
		if len(status.Status) == 0 {
			status.Status = metav1.StatusFailure
		}
		switch status.Status {
		case metav1.StatusSuccess:
			if status.Code == 0 {
				status.Code = http.StatusOK
			}
		case metav1.StatusFailure:
			if status.Code == 0 {
				status.Code = http.StatusInternalServerError
			}
		default:
			if status.Code == 0 {
				status.Code = http.StatusInternalServerError
			}
		}
		status.Kind = "Status"
		status.APIVersion = "meta/v1"
		return &status
	default:
		status := http.StatusInternalServerError
		return &metav1.Status{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "meta/v1",
				Kind:       "Status",
			},
			Status:  metav1.StatusFailure,
			Message: err.Error(),
			Reason:  metav1.StatusReasonUnknown,
			Code:    int32(status),
		}
	}
}

// WriteObjectNegotiated renders an object in the content type negotiated by the client.
func WriteObjectNegotiated(s runtime.NegotiatedSerializer, restrictions negotiation.EndpointRestrictions, gv schema.GroupVersion, w http.ResponseWriter, req *http.Request, statusCode int, object runtime.Object) {
	//stream, ok := object.(rest.ResourceStreamer)
	//if ok {
	//  requestInfo, _ := request.RequestInfoFrom(req.Context())
	//  metrics.RecordLongRunning(req, requestInfo, metrics.APIServerComponent, func() {
	//    StreamObject(statusCode, gv, s, stream, w, req)
	//  })
	//  return
	//}

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

// SerializeObject renders an object in the content type negotiated by the client using the provided encoder.
// The context is optional and can be nil. This method will perform optional content compression if requested by
// a client and the feature gate for APIResponseCompression is enabled.
func SerializeObject(mediaType string, encoder runtime.Encoder, hw http.ResponseWriter, req *http.Request, statusCode int, object runtime.Object) {
	//trace := utiltrace.New("SerializeObject",
	//  utiltrace.Field{"method", req.Method},
	//  utiltrace.Field{"url", req.URL.Path},
	//  utiltrace.Field{"protocol", req.Proto},
	//  utiltrace.Field{"mediaType", mediaType},
	//  utiltrace.Field{"encoder", encoder.Identifier()})
	//defer trace.LogIfLong(5 * time.Second)
	//
	//w := &deferredResponseWriter{
	//  mediaType:       mediaType,
	//  statusCode:      statusCode,
	//  contentEncoding: negotiateContentEncoding(req),
	//  hw:              hw,
	//  trace:           trace,
	//}

	err := encoder.Encode(object, hw)
	if err == nil {
		// err = w.Close()
		//if err != nil {
		//  // we cannot write an error to the writer anymore as the Encode call was successful.
		//  logrus.Errorf("apiserver was unable to close cleanly the response writer: %v", err)
		//  // utilruntime.HandleError(fmt.Errorf("apiserver was unable to close cleanly the response writer: %v", err))
		//}
		return
	}

	logrus.Errorf("apiserver was unable to write a JSON response: %v", err)
	// make a best effort to write the object if a failure is detected
	// utilruntime.HandleError(fmt.Errorf("apiserver was unable to write a JSON response: %v", err))
	status := ErrorToAPIStatus(err)
	candidateStatusCode := int(status.Code)
	// if the current status code is successful, allow the error's status code to overwrite it
	if statusCode >= http.StatusOK && statusCode < http.StatusBadRequest {
		hw.WriteHeader(candidateStatusCode)
		// w.statusCode = candidateStatusCode
	}
	output, err := runtime.Encode(encoder, status)
	if err != nil {
		hw.Header().Set("Content-Type", "text/plain")
		// hw.Header["Content-Type"] ="text/plain"
		// w.mediaType = "text/plain"
		output = []byte(fmt.Sprintf("%s: %s", status.Reason, status.Message))
	}
	if _, err := hw.Write(output); err != nil {
		logrus.Errorf("apiserver was unable to write a fallback JSON response: %v", err)
	}
	// w.Close()
}

// ErrorToAPIStatus converts an error to an metav1.Status object.
func ErrorToAPIStatus(err error) *metav1.Status {
	switch t := err.(type) {
	case statusError:
		status := t.Status()
		if len(status.Status) == 0 {
			status.Status = metav1.StatusFailure
		}
		switch status.Status {
		case metav1.StatusSuccess:
			if status.Code == 0 {
				status.Code = http.StatusOK
			}
		case metav1.StatusFailure:
			if status.Code == 0 {
				status.Code = http.StatusInternalServerError
			}
		default:
			logrus.Errorf("apiserver received an error with wrong status field : %#+v", err)
			// runtime.HandleError(fmt.Errorf("apiserver received an error with wrong status field : %#+v", err))
			if status.Code == 0 {
				status.Code = http.StatusInternalServerError
			}
		}
		status.Kind = "Status"
		status.APIVersion = "v1"
		//TODO: check for invalid responses
		return &status
	default:
		status := http.StatusInternalServerError
		switch {
		//TODO: replace me with NewConflictErr
		case storage.IsConflict(err):
			status = http.StatusConflict
		}
		// Log errors that were not converted to an error status
		// by REST storage - these typically indicate programmer
		// error by not using pkg/api/errors, or unexpected failure
		// cases.
		logrus.Errorf("apiserver received an error that is not an metav1.Status: %#+v: %v", err, err)
		//runtime.HandleError(fmt.Errorf("apiserver received an error that is not an metav1.Status: %#+v: %v", err, err))
		return &metav1.Status{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Status",
				APIVersion: "v1",
			},
			Status:  metav1.StatusFailure,
			Code:    int32(status),
			Reason:  metav1.StatusReasonUnknown,
			Message: err.Error(),
		}
	}
}
