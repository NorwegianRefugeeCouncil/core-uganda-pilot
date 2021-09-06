package validation

import (
	uuid "github.com/satori/go.uuid"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"unicode"
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

func (s *Status) Error() string {
	return s.Message
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

func IsValidNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

var InvalidNumericDetail = "Accepted: numbers"

func IsValidAlphaNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

var InvalidAlphaNumericDetail = "Accepted: letters, numbers, spaces"

func IsValidEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

var InvalidEmailDetail = "invalid email"

func IsValidPassword(s string) bool {
	// TODO
	return true
}

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

func IsValidURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

var InvalidURLDetail = "invalid URL"
