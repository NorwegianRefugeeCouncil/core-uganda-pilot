package handlers

import (
	"fmt"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"net/http"
)

type errNotAcceptable struct {
}

func NewNotAcceptableError() error {
	return errNotAcceptable{}
}

func (e errNotAcceptable) Error() string {
	return fmt.Sprintf("media-type not acceptable")
}

func (e errNotAcceptable) Status() metav1.Status {
	return metav1.Status{
		Status:  metav1.StatusFailure,
		Message: e.Error(),
		Reason:  metav1.StatusReasonNotAcceptable,
		Code:    http.StatusNotAcceptable,
	}
}

type errBadRequest struct {
}

func (e errBadRequest) Error() string {
	return ""
}

func (e errBadRequest) Status() metav1.Status {
	return metav1.Status{
		Status:  metav1.StatusFailure,
		Message: e.Error(),
		Reason:  metav1.StatusReasonBadRequest,
		Code:    http.StatusBadRequest,
	}
}
