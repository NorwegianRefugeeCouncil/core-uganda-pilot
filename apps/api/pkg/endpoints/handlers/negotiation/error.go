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
