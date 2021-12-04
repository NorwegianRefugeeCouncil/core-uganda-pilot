package schema

import "fmt"

type ErrorCode uint8

const (
	ColumnNotFound = iota
	TableNotFound
	TableConstraintNotFound
)

var (
	ErrColumnNotFound          = &SQLError{Code: ColumnNotFound}
	ErrTableNotFound           = &SQLError{Code: TableNotFound}
	ErrTableConstraintNotFound = &SQLError{Code: TableConstraintNotFound}
)

func newColumnNotFoundErr(columnName string) error {
	return &SQLError{
		Code:    ColumnNotFound,
		Message: fmt.Sprintf("column with name '%s' not found", columnName),
	}
}

func newTableNotFoundErr(tableName string) error {
	return &SQLError{
		Code:    TableNotFound,
		Message: fmt.Sprintf("table with name '%s' not found", tableName),
	}
}
func newTableConstraintNotFoundErr(constraintName string) error {
	return &SQLError{
		Code:    TableConstraintNotFound,
		Message: fmt.Sprintf("constraint with name '%s' not found", constraintName),
	}
}

type SQLError struct {
	Message string
	Code    ErrorCode
}

func (s *SQLError) Error() string {
	return s.Message
}

func (s *SQLError) Is(err error) bool {
	if sqlErr, ok := err.(*SQLError); ok {
		return sqlErr.Code == s.Code
	}
	return false
}
