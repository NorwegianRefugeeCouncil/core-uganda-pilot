package handlers

import (
  "encoding/hex"
  "fmt"
  "github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/responsewriters"
  "github.com/nrc-no/core/apps/api/pkg/registry/rest"
  "io"
  "io/ioutil"
  "k8s.io/apimachinery/pkg/api/errors"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/apimachinery/pkg/runtime"
  "k8s.io/apimachinery/pkg/runtime/schema"
  "net/http"
  "time"
)

const (
  requestTimeoutUpperBound = 34 * time.Second
)

type ScopeNamer interface {
  Name(req *http.Request) (namespace, name string, err error)
  ObjectName(obj runtime.Object) (namespace, name string, err error)
}

type RequestScope struct {
  Namer               ScopeNamer
  Serializer          runtime.NegotiatedSerializer
  ParameterCoder      runtime.ParameterCodec
  Creater             runtime.ObjectCreater
  Resource            schema.GroupVersionResource
  Kind                schema.GroupVersionKind
  StandardSerializers []runtime.SerializerInfo
  TableConvertor      rest.TableConvertor
  MetaGroupVersion    schema.GroupVersion
  HubGroupVersion     schema.GroupVersion
  MaxRequestBodyBytes int64
  Typer               runtime.ObjectTyper
  Convertor           runtime.ObjectConvertor
}

func (scope *RequestScope) err(err error, w http.ResponseWriter, req *http.Request) {
  responsewriters.ErrorNegotiated(err, scope.Serializer, scope.Kind.GroupVersion(), w, req)
}

// AcceptsGroupVersion returns true if the specified GroupVersion is allowed
// in create and update requests.
func (scope *RequestScope) AcceptsGroupVersion(gv schema.GroupVersion) bool {
  // Fall back to only allowing the singular Kind. This is the typical behavior.
  return gv == scope.Kind.GroupVersion()
}

func (scope *RequestScope) AllowsMediaTypeTransform(mimeType, mimeSubType string, gvk *schema.GroupVersionKind) bool {
  // some handlers like CRDs can't serve all the mime types that PartialObjectMetadata or Table can - if
  // gvk is nil (no conversion) allow StandardSerializers to further restrict the set of mime types.
  if gvk == nil {
    if len(scope.StandardSerializers) == 0 {
      return true
    }
    for _, info := range scope.StandardSerializers {
      if info.MediaTypeType == mimeType && info.MediaTypeSubType == mimeSubType {
        return true
      }
    }
    return false
  }

  // TODO: this is temporary, replace with an abstraction calculated at endpoint installation time
  if gvk.GroupVersion() == metav1.SchemeGroupVersion {
    switch gvk.Kind {
    case "Table":
      return scope.TableConvertor != nil &&
        mimeType == "application" &&
        (mimeSubType == "json" || mimeSubType == "yaml")
    case "PartialObjectMetadata", "PartialObjectMetadataList":
      // TODO: should delineate between lists and non-list endpoints
      return true
    default:
      return false
    }
  }
  return false
}

func (scope *RequestScope) AllowsServerVersion(version string) bool {
  return version == scope.MetaGroupVersion.Version
}

func (scope *RequestScope) AllowsStreamSchema(s string) bool {
  return s == "watch"
}

func limitedReadBody(req *http.Request, limit int64) ([]byte, error) {
  defer req.Body.Close()
  if limit <= 0 {
    return ioutil.ReadAll(req.Body)
  }
  lr := &io.LimitedReader{
    R: req.Body,
    N: limit + 1,
  }
  data, err := ioutil.ReadAll(lr)
  if err != nil {
    return nil, err
  }
  if lr.N <= 0 {
    return nil, errors.NewRequestEntityTooLargeError(fmt.Sprintf("limit is %d", limit))
  }
  return data, nil
}

// transformDecodeError adds additional information into a bad-request api error when a decode fails.
func transformDecodeError(typer runtime.ObjectTyper, baseErr error, into runtime.Object, gvk *schema.GroupVersionKind, body []byte) error {
  objGVKs, _, err := typer.ObjectKinds(into)
  if err != nil {
    return errors.NewBadRequest(err.Error())
  }
  objGVK := objGVKs[0]
  if gvk != nil && len(gvk.Kind) > 0 {
    return errors.NewBadRequest(fmt.Sprintf("%s in version %q cannot be handled as a %s: %v", gvk.Kind, gvk.Version, objGVK.Kind, baseErr))
  }
  summary := summarizeData(body, 30)
  return errors.NewBadRequest(fmt.Sprintf("the object provided is unrecognized (must be of type %s): %v (%s)", objGVK.Kind, baseErr, summary))
}

func summarizeData(data []byte, maxLength int) string {
  switch {
  case len(data) == 0:
    return "<empty>"
  case data[0] == '{':
    if len(data) > maxLength {
      return string(data[:maxLength]) + " ..."
    }
    return string(data)
  default:
    if len(data) > maxLength {
      return hex.EncodeToString(data[:maxLength]) + " ..."
    }
    return hex.EncodeToString(data)
  }
}
