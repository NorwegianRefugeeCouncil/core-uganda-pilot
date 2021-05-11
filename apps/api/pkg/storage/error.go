package storage

import "fmt"

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
