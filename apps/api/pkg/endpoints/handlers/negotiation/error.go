package negotiation

import (
  "fmt"
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
