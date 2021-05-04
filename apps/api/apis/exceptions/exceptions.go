package exceptions

import (
	"fmt"
	"net/http"
)

var (
	ErrNotFound = NewAPIError("ErrNotFound", http.StatusNotFound)
	ErrConflict = NewAPIError("ErrConflict", http.StatusConflict)
)

type ErrorResponse struct {
}

type APIError struct {
	Code   string `json:"code"`
	Status int    `json:"status"`
	Err    error  `json:"error"`
}

func NewAPIError(code string, status int) *APIError {
	return &APIError{
		Code:   code,
		Status: status,
	}
}

func (a *APIError) WithError(err error) *APIError {
	return &APIError{
		Code:   a.Code,
		Status: a.Status,
		Err:    err,
	}
}

func (a *APIError) Error() string {
	return fmt.Sprintf("[%s] API error [%d]: %s", a.Code, a.Status, a.Err.Error())
}

func (a *APIError) Unwrap() error {
	return a.Err
}

func (a *APIError) Is(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Code == a.Code
	}
	return false
}
