package exceptions

import (
  "errors"
  metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
  "github.com/nrc-no/core/apps/api/pkg/runtime"
  "net/http"
  "reflect"
)

// IsNotFound returns true if the specified error was created by NewNotFound.
// It supports wrapped errors.
func IsNotFound(err error) bool {
  return ReasonForError(err) == metav1.StatusReasonNotFound
}

// IsAlreadyExists determines if the err is an error which indicates that a specified resource already exists.
// It supports wrapped errors.
func IsAlreadyExists(err error) bool {
  return ReasonForError(err) == metav1.StatusReasonAlreadyExists
}

// IsConflict determines if the err is an error which indicates the provided update conflicts.
// It supports wrapped errors.
func IsConflict(err error) bool {
  return ReasonForError(err) == metav1.StatusReasonConflict
}

// IsInvalid determines if the err is an error which indicates the provided resource is not valid.
// It supports wrapped errors.
func IsInvalid(err error) bool {
  return ReasonForError(err) == metav1.StatusReasonInvalid
}

// IsGone is true if the error indicates the requested resource is no longer available.
// It supports wrapped errors.
func IsGone(err error) bool {
  return ReasonForError(err) == metav1.StatusReasonGone
}

// IsResourceExpired is true if the error indicates the resource has expired and the current action is
// no longer possible.
// It supports wrapped errors.
func IsResourceExpired(err error) bool {
  return ReasonForError(err) == metav1.StatusReasonExpired
}

// IsNotAcceptable determines if err is an error which indicates that the request failed due to an invalid Accept header
// It supports wrapped errors.
func IsNotAcceptable(err error) bool {
  return ReasonForError(err) == metav1.StatusReasonNotAcceptable
}

// IsUnsupportedMediaType determines if err is an error which indicates that the request failed due to an invalid Content-Type header
// It supports wrapped errors.
func IsUnsupportedMediaType(err error) bool {
  return ReasonForError(err) == metav1.StatusReasonUnsupportedMediaType
}

// IsMethodNotSupported determines if the err is an error which indicates the provided action could not
// be performed because it is not supported by the server.
// It supports wrapped errors.
func IsMethodNotSupported(err error) bool {
  return ReasonForError(err) == metav1.StatusReasonMethodNotAllowed
}

// IsServiceUnavailable is true if the error indicates the underlying service is no longer available.
// It supports wrapped errors.
func IsServiceUnavailable(err error) bool {
  return ReasonForError(err) == metav1.StatusReasonServiceUnavailable
}

// IsBadRequest determines if err is an error which indicates that the request is invalid.
// It supports wrapped errors.
func IsBadRequest(err error) bool {
  return ReasonForError(err) == metav1.StatusReasonBadRequest
}

// IsUnauthorized determines if err is an error which indicates that the request is unauthorized and
// requires authentication by the user.
// It supports wrapped errors.
func IsUnauthorized(err error) bool {
  return ReasonForError(err) == metav1.StatusReasonUnauthorized
}

// IsForbidden determines if err is an error which indicates that the request is forbidden and cannot
// be completed as requested.
// It supports wrapped errors.
func IsForbidden(err error) bool {
  return ReasonForError(err) == metav1.StatusReasonForbidden
}

// IsTimeout determines if err is an error which indicates that request times out due to long
// processing.
// It supports wrapped errors.
func IsTimeout(err error) bool {
  return ReasonForError(err) == metav1.StatusReasonTimeout
}

// IsServerTimeout determines if err is an error which indicates that the request needs to be retried
// by the client.
// It supports wrapped errors.
func IsServerTimeout(err error) bool {
  return ReasonForError(err) == metav1.StatusReasonServerTimeout
}

// IsInternalError determines if err is an error which indicates an internal server error.
// It supports wrapped errors.
func IsInternalError(err error) bool {
  return ReasonForError(err) == metav1.StatusReasonInternalError
}

// IsTooManyRequests determines if err is an error which indicates that there are too many requests
// that the server cannot handle.
// It supports wrapped errors.
func IsTooManyRequests(err error) bool {
  if ReasonForError(err) == metav1.StatusReasonTooManyRequests {
    return true
  }
  if status := APIStatus(nil); errors.As(err, &status) {
    return status.Status().Code == http.StatusTooManyRequests
  }
  return false
}

// IsRequestEntityTooLargeError determines if err is an error which indicates
// the request entity is too large.
// It supports wrapped errors.
func IsRequestEntityTooLargeError(err error) bool {
  if ReasonForError(err) == metav1.StatusReasonRequestEntityTooLarge {
    return true
  }
  if status := APIStatus(nil); errors.As(err, &status) {
    return status.Status().Code == http.StatusRequestEntityTooLarge
  }
  return false
}

// IsUnexpectedServerError returns true if the server response was not in the expected API format,
// and may be the result of another HTTP actor.
// It supports wrapped errors.
func IsUnexpectedServerError(err error) bool {
  if status := APIStatus(nil); errors.As(err, &status) && status.Status().Details != nil {
    for _, cause := range status.Status().Details.Causes {
      if cause.Type == metav1.CauseTypeUnexpectedServerResponse {
        return true
      }
    }
  }
  return false
}

// IsUnexpectedObjectError determines if err is due to an unexpected object from the master.
// It supports wrapped errors.
func IsUnexpectedObjectError(err error) bool {
  uoe := &UnexpectedObjectError{}
  return err != nil && errors.As(err, &uoe)
}

// ReasonForError returns the HTTP status for a particular error.
// It supports wrapped errors.
func ReasonForError(err error) metav1.StatusReason {
  if status := APIStatus(nil); errors.As(err, &status) {
    return status.Status().Reason
  }
  return metav1.StatusReasonUnknown
}

// FromObject generates an StatusError from an metav1.Status, if that is the type of obj; otherwise,
// returns an UnexpecteObjectError.
func FromObject(obj runtime.Object) error {
  switch t := obj.(type) {
  case *metav1.Status:
    return &StatusError{ErrStatus: *t}
  case runtime.Unstructured:
    var status metav1.Status
    obj := t.UnstructuredContent()
    if !reflect.DeepEqual(obj["kind"], "Status") {
      break
    }
    if err := runtime.DefaultUnstructuredConverter.FromUnstructured(t.UnstructuredContent(), &status); err != nil {
      return err
    }
    if status.APIVersion != "v1" && status.APIVersion != "meta.k8s.io/v1" {
      break
    }
    return &StatusError{ErrStatus: status}
  }
  return &UnexpectedObjectError{obj}
}

// HasStatusCause returns true if the provided error has a details cause
// with the provided type name.
func HasStatusCause(err error, name metav1.CauseType) bool {
  _, ok := StatusCause(err, name)
  return ok
}

// StatusCause returns the named cause from the provided error if it exists and
// the error is of the type APIStatus. Otherwise it returns false.
func StatusCause(err error, name metav1.CauseType) (metav1.StatusCause, bool) {
  apierr, ok := err.(APIStatus)
  if !ok || apierr == nil || apierr.Status().Details == nil {
    return metav1.StatusCause{}, false
  }
  for _, cause := range apierr.Status().Details.Causes {
    if cause.Type == name {
      return cause, true
    }
  }
  return metav1.StatusCause{}, false
}
