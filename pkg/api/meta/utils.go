package meta

import (
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/validation"
	"net/http"
)

// NewNotFound returns a new error which indicates that the resource of the kind and the name was not found.
func NewNotFound(qualifiedResource GroupResourcer, uid string) *StatusError {
	r := qualifiedResource.GroupResource()
	return &StatusError{Status{
		Status: StatusFailure,
		Code:   http.StatusNotFound,
		Reason: StatusReasonNotFound,
		Details: &StatusDetails{
			UID:      uid,
			Group:    r.Group,
			Resource: r.Resource,
		},
		Message: fmt.Sprintf("%s %s not found", r.String(), uid),
	}}
}

// NewAlreadyExists returns an error indicating the item requested exists by that identifier.
func NewAlreadyExists(qualifiedResource GroupResource, uid string) *StatusError {
	return &StatusError{Status{
		Status: StatusFailure,
		Code:   http.StatusConflict,
		Reason: StatusReasonAlreadyExists,
		Details: &StatusDetails{
			UID:      uid,
			Group:    qualifiedResource.Group,
			Resource: qualifiedResource.Resource,
		},
		Message: fmt.Sprintf("%s %q already exists", qualifiedResource.String(), uid),
	}}
}

// NewInvalid returns an error indicating the item is invalid and cannot be processed.
func NewInvalid(qualifiedResource GroupResourcer, uid string, errs validation.ErrorList) *StatusError {
	r := qualifiedResource.GroupResource()
	causes := make([]StatusCause, 0, len(errs))
	for i := range errs {
		err := errs[i]
		causes = append(causes, StatusCause{
			Type:    CauseType(err.Type),
			Message: err.ErrorBody(),
			Field:   err.Field,
		})
	}
	return &StatusError{Status{
		Status: StatusFailure,
		Code:   http.StatusUnprocessableEntity,
		Reason: StatusReasonInvalid,
		Details: &StatusDetails{
			Group:    r.Group,
			Resource: r.Resource,
			UID:      uid,
			Causes:   causes,
		},
		Message: fmt.Sprintf("%s %q is invalid: %v", r.String(), uid, errs.ToAggregate()),
	}}
}

// NewBadRequest creates an error that indicates that the request is invalid and can not be processed.
func NewBadRequest(reason string) *StatusError {
	return &StatusError{Status{
		Status:  StatusFailure,
		Code:    http.StatusBadRequest,
		Reason:  StatusReasonBadRequest,
		Message: reason,
	}}
}

func NewInternalServerError(err error) *StatusError {
	return &StatusError{Status{
		Status:  StatusFailure,
		Code:    http.StatusInternalServerError,
		Reason:  StatusReasonInternalError,
		Message: err.Error(),
	}}
}

// ReasonForError returns the HTTP status for a particular error.
// It supports wrapped errors.
func ReasonForError(err error) StatusReason {
	if status := APIStatus(nil); errors.As(err, &status) {
		return status.Status().Reason
	}
	return StatusReasonUnknown
}

func reasonAndCodeForError(err error) (StatusReason, int32) {
	if status := APIStatus(nil); errors.As(err, &status) {
		return status.Status().Reason, status.Status().Code
	}
	return StatusReasonUnknown, 0
}
