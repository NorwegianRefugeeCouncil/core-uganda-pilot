package validation

import (
	"net/mail"
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

var alphaRegexp = regexp.MustCompile(`^[\w\-()]+( [\w\-()]+)*$`)

func IsValidAlpha(s string) bool {
	return alphaRegexp.MatchString(s)
}

var InvalidAlphaDetail = "Accepted: letters, spaces"

func IsValidNumeric(s string) bool {
	if _, err := strconv.Atoi(s); err != nil {
		return false
	}
	return true
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

var InvalidEmailDetail = "Invalid email"

func IsValidPassword(s string) bool {
	// TODO
	return true
}

func IsValidPhone(s string) bool {
	// TODO
	return true
}

func IsValidUUID(s string) bool {
	// TODO
	return true
}
