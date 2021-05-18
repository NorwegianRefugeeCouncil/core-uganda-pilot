package exceptions

import v1 "github.com/nrc-no/core/api/pkg/apis/meta/v1"

type StatusError struct {
	ErrStatus v1.Status
}

func (e *StatusError) Error() string {
	return e.ErrStatus.Message
}

func NewStatusError(err v1.Status) *StatusError {
	return &StatusError{
		ErrStatus: err,
	}
}
