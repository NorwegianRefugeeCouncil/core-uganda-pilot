package validation

import (
	"fmt"
)

type StatusType string

const (
	Success StatusType = "Success"
	Failure StatusType = "Failure"
)

type Status struct {
	Status  StatusType `json:"status"`
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Errors  ErrorList  `json:"errors"`
}

func (s Status) Error() string {
	return s.Message
}

func (s Status) Unwrap() error {
	return fmt.Errorf("%s", s.Message)
}
