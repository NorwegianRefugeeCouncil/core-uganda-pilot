package api

type Error struct {
	Message string
	Code    ErrorCode
}

type ErrorCode uint8

const (
	ErrCodeInvalidRevision ErrorCode = iota
	ErrCodeInvalidPrevision
	ErrCodeInvalidRecordID
	ErrCodeDuplicateField
	ErrCodeMissingRevision
	ErrCodeInvalidTable
	ErrCodeInvalidColumnType
	ErrCodeEmptyTableColumns
	ErrCodeDuplicateColumnName
	ErrCodeInvalidColumnName
	ErrCodeRecordNotFound
	ErrCodeFieldNotFound
	ErrCodeInternalError
	ErrCodeTableAlreadyExists
	ErrCodeUnsupportedDialect
	ErrCodeInvalidTimestamp
	ErrCodeMissingId
)

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) ErrorCode() ErrorCode {
	return e.Code
}

func (e *Error) Is(other error) bool {
	if other == nil {
		return false
	}
	if e == other {
		return true
	}
	if o, ok := other.(*Error); ok {
		return e.Code == o.Code
	}
	return false
}

func IsError(err error, code ErrorCode) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*Error); ok {
		return e.Code == code
	}
	return false
}

func NewError(code ErrorCode, message string) *Error {
	return &Error{
		Message: message,
		Code:    code,
	}
}

var (
	ErrInvalidRevision     = NewError(ErrCodeInvalidRevision, "invalid revision")
	ErrInvalidPrevision    = NewError(ErrCodeInvalidPrevision, "invalid previous revision")
	ErrInvalidRecordID     = NewError(ErrCodeInvalidRecordID, "invalid record id")
	ErrDuplicateField      = NewError(ErrCodeDuplicateField, "duplicate field")
	ErrMissingRevision     = NewError(ErrCodeMissingRevision, "missing revision")
	ErrInvalidTableName    = NewError(ErrCodeInvalidTable, "invalid table name")
	ErrEmptyColumns        = NewError(ErrCodeEmptyTableColumns, "empty columns")
	ErrDuplicateColumnName = NewError(ErrCodeDuplicateColumnName, "duplicate column name")
	ErrInvalidColumnName   = NewError(ErrCodeInvalidColumnName, "invalid column name")
	ErrRecordNotFound      = NewError(ErrCodeRecordNotFound, "record not found")
	ErrFieldNotFound       = NewError(ErrCodeFieldNotFound, "field not found")
	ErrInvalidColumnType   = NewError(ErrCodeInvalidColumnType, "invalid column type")
	ErrTableAlreadyExists  = NewError(ErrCodeTableAlreadyExists, "table already exists")
	ErrUnsupportedDialect  = NewError(ErrCodeUnsupportedDialect, "unsupported dialect")
	ErrInvalidValueType    = NewError(ErrCodeInternalError, "invalid value type")
	ErrInvalidTimestamp    = NewError(ErrCodeInvalidTimestamp, "invalid timestamp")
	ErrMissingId           = NewError(ErrCodeMissingId, "missing id")
)

func NewDuplicateColumnNameErr(name string) *Error {
	return NewError(ErrCodeDuplicateColumnName, "duplicate column name: "+name)
}

func NewInvalidColumnNameErr(name string) *Error {
	return NewError(ErrCodeInvalidColumnName, "invalid column name: "+name)
}

func NewInvalidColumnTypeErr(name string) *Error {
	return NewError(ErrCodeInvalidColumnType, "invalid column type: "+name)
}

func NewTableAlreadyExistsErr(name string) *Error {
	return NewError(ErrCodeTableAlreadyExists, "table already exists: "+name)
}
