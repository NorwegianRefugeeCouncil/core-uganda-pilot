package storage

import (
	"fmt"
	"github.com/nrc-no/core/apps/api/pkg/util/exceptions"
)

const (
	ErrCodeKeyNotFound int = iota + 1
	ErrCodeKeyExists
	ErrCodeResourceVersionConflicts
	ErrCodeInvalidObj
	ErrCodeUnreachable
)

var errCodeToMessage = map[int]string{
	ErrCodeKeyNotFound:              "key not found",
	ErrCodeKeyExists:                "key exists",
	ErrCodeResourceVersionConflicts: "resource version conflicts",
	ErrCodeInvalidObj:               "invalid object",
	ErrCodeUnreachable:              "server unreachable",
}

type StorageError struct {
	Code               int
	Key                string
	ResourceVersion    int64
	AdditionalErrorMsg string
}

func (e *StorageError) Error() string {
	return fmt.Sprintf("StorageError: %s, Code: %d, Key: %s, ResourceVersion: %d, AdditionalErrorMsg: %s",
		errCodeToMessage[e.Code], e.Code, e.Key, e.ResourceVersion, e.AdditionalErrorMsg)
}

// IsConflict returns true if and only if err is a write conflict.
func IsConflict(err error) bool {
	return isErrCode(err, ErrCodeResourceVersionConflicts)
}

func isErrCode(err error, code int) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*StorageError); ok {
		return e.Code == code
	}
	return false
}

// InvalidError is generated when an error caused by invalid API object occurs
// in the storage package.
type InvalidError struct {
	Errs exceptions.ErrorList
}

func (e InvalidError) Error() string {
	return e.Errs.ToAggregate().Error()
}

// IsInvalidError returns true if and only if err is an InvalidError.
func IsInvalidError(err error) bool {
	_, ok := err.(InvalidError)
	return ok
}

func NewInvalidError(errors exceptions.ErrorList) InvalidError {
	return InvalidError{errors}
}
