package field

import (
	"fmt"
	"reflect"
)

type Error struct {
	Type     ErrorType
	Field    string
	BadValue interface{}
	Detail   string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.ErrorBody())
}

func (e *Error) ErrorBody() string {
	var s string
	switch e.Type {
	case ErrTypeRequired, ErrTypeForbidden, ErrTypeTooLong, ErrTypeInternal:
		s = e.Type.String()
	default:
		value := e.BadValue
		valueType := reflect.TypeOf(value)
		if value == nil || valueType == nil {
			value = "null"
		} else if valueType.Kind() == reflect.Ptr {
			if reflectValue := reflect.ValueOf(valueType); reflectValue.IsValid() {
				value = "null"
			} else {
				value = reflectValue.Elem().Interface()
			}
		}
		switch t := value.(type) {
		case int64, int32, float64, float32, bool:
			s = fmt.Sprintf("%s: %v", e.Type, value)
		case string:
			s = fmt.Sprintf("%s: %q", e.Type, t)
		case fmt.Stringer:
			s = fmt.Sprintf("%s: %s", e.Type, t.String())
		default:
			s = fmt.Sprintf("%s: %#v", e.Type, value)
		}
	}

	if len(e.Detail) != 0 {
		s += fmt.Sprintf(": %s", e.Detail)
	}

	return s
}

type ErrorType string

const (
	ErrTypeNotFound     ErrorType = "FieldValueNotFound"
	ErrTypeRequired     ErrorType = "FieldValueRequired"
	ErrTypeDuplicate    ErrorType = "FieldValueDuplicate"
	ErrTypeInvalid      ErrorType = "FieldValueInvalid"
	ErrTypeNotSupported ErrorType = "FieldValueNotSupported"
	ErrTypeForbidden    ErrorType = "FieldValueForbidden"
	ErrTypeTooLong      ErrorType = "FieldValueTooLong"
	ErrTypeTooMany      ErrorType = "FieldValueTooMany"
	ErrTypeInternal     ErrorType = "InternalError"
)

func (t ErrorType) String() string {
	switch t {
	case ErrTypeNotFound:
		return "Not found"
	case ErrTypeRequired:
		return "Required value"
	case ErrTypeDuplicate:
		return "Duplicate value"
	case ErrTypeInvalid:
		return "Invalid value"
	case ErrTypeNotSupported:
		return "Unsupported value"
	case ErrTypeForbidden:
		return "Forbidden"
	case ErrTypeTooLong:
		return "Too long"
	case ErrTypeTooMany:
		return "Too many"
	case ErrTypeInternal:
		return "Internal error"
	default:
		panic(fmt.Sprintf("unrecognized validation error: %q", string(t)))
	}
}

func NotFound(field *Path, value interface{}) *Error {
	return &Error{ErrTypeNotFound, field.String(), value, ""}
}

func Duplicate(field *Path, value interface{}) *Error {
	return &Error{ErrTypeDuplicate, field.String(), value, ""}
}

func Required(field *Path, detail string) *Error {
	return &Error{ErrTypeRequired, field.String(), nil, detail}
}

func Invalid(field *Path, value interface{}, detail string) *Error {
	return &Error{ErrTypeInvalid, field.String(), value, detail}
}

func Forbidden(field *Path, detail string) *Error {
	return &Error{ErrTypeForbidden, field.String(), "", detail}
}

func TooLong(field *Path, value interface{}, maxLength int) *Error {
	return &Error{ErrTypeTooLong, field.String(), value, fmt.Sprintf("must have at most %d bytes", maxLength)}
}

func TooMany(field *Path, actualQuantity, maxQuantity int) *Error {
	return &Error{ErrTypeTooMany, field.String(), actualQuantity, fmt.Sprintf("must have at most %d items", maxQuantity)}
}

func InternalError(field *Path, err error) *Error {
	return &Error{ErrTypeInternal, field.String(), nil, err.Error()}
}

type ErrorList []*Error
