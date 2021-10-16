package validation

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/mail"
	"net/url"
	"regexp"
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

func AsStatus(err error) Status {
	var status *Status
	if errors.As(err, &status) {
		return *status
	}
	return Status{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
		Status:  Failure,
	}
}

func IsValidUUID(s string) bool {
	_, err := uuid.FromString(s)
	return err == nil
}

var InvalidUUIDDetail = "invalid UUID"

var alphaRegexp = regexp.MustCompile(`^[\w\-()]+( [\w\-()]+)*$`)

func IsValidAlpha(s string) bool {
	return alphaRegexp.MatchString(s)
}

var InvalidAlphaDetail = "Accepted: letters, spaces"

var InvalidNumericDetail = "Accepted: numbers"

func IsValidEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

var InvalidEmailDetail = "invalid email"

var ugandaPhoneFormat = regexp.MustCompile(`^(\+?256(\s|-)?|0)\d{3}(\s|-)?\d{6}$`)

func IsValidPhone(s string) bool {
	return ugandaPhoneFormat.MatchString(s)
}

var InvalidPhoneDetail = "invalid phone number"

// input[type="date"] will always yield yyyy-mm-dd regardless of locale (and user facing format)
var dateFormat = regexp.MustCompile(`^\d{4}-(0\d|1[012])-([012]\d|3[01])`)

func IsValidDate(s string) bool {
	return dateFormat.MatchString(s)
}

var InvalidDateDetail = "invalid date"

func IsValidTime(s string) bool {
	return dateFormat.MatchString(s)
}

var InvalidTimeDetail = "invalid time"

func IsValidURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

var InvalidURLDetail = "invalid URL"
