package negotiation

import (
	"fmt"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"net/http"
	"strings"
)

// errNotAcceptable indicates Accept negotiation has failed
type errNotAcceptable struct {
	accepted []string
}

// NewNotAcceptableError returns an error of NotAcceptable which contains specified string
func NewNotAcceptableError(accepted []string) error {
	return errNotAcceptable{accepted}
}

func (e errNotAcceptable) Error() string {
	return fmt.Sprintf("only the following media types are accepted: %v", strings.Join(e.accepted, ", "))
}

func (e errNotAcceptable) Status() metav1.Status {
	return metav1.Status{
		Status:  metav1.StatusFailure,
		Code:    http.StatusNotAcceptable,
		Reason:  metav1.StatusReasonNotAcceptable,
		Message: e.Error(),
	}
}

// errNotAcceptableConversion indicates Accept negotiation has failed specifically
// for a conversion to a known type.
type errNotAcceptableConversion struct {
	target   string
	accepted []string
}

// NewNotAcceptableConversionError returns an error indicating that the desired
// API transformation to the target group version kind string is not accepted and
// only the listed mime types are allowed. This is temporary while Table does not
// yet support protobuf encoding.
func NewNotAcceptableConversionError(target string, accepted []string) error {
	return errNotAcceptableConversion{target, accepted}
}

func (e errNotAcceptableConversion) Error() string {
	return fmt.Sprintf("only the following media types are accepted when converting to %s: %v", e.target, strings.Join(e.accepted, ", "))
}

func (e errNotAcceptableConversion) Status() metav1.Status {
	return metav1.Status{
		Status:  metav1.StatusFailure,
		Code:    http.StatusNotAcceptable,
		Reason:  metav1.StatusReasonNotAcceptable,
		Message: e.Error(),
	}
}

// errUnsupportedMediaType indicates Content-Type is not recognized
type errUnsupportedMediaType struct {
	accepted []string
}

// NewUnsupportedMediaTypeError returns an error of UnsupportedMediaType which contains specified string
func NewUnsupportedMediaTypeError(accepted []string) error {
	return errUnsupportedMediaType{accepted}
}

func (e errUnsupportedMediaType) Error() string {
	return fmt.Sprintf("the body of the request was in an unknown format - accepted media types include: %v", strings.Join(e.accepted, ", "))
}

func (e errUnsupportedMediaType) Status() metav1.Status {
	return metav1.Status{
		Status:  metav1.StatusFailure,
		Code:    http.StatusUnsupportedMediaType,
		Reason:  metav1.StatusReasonUnsupportedMediaType,
		Message: e.Error(),
	}
}
