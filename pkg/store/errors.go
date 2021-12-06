package store

import (
	"errors"
	"github.com/jackc/pgconn"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

const (
	// https://www.postgresql.org/docs/13/errcodes-appendix.html
	errPgUniqueViolation = "23505"
)

// IsUniqueConstraintErr returns whether the error is a Unique Constraint Violation error or not
func IsUniqueConstraintErr(err error) bool {
	if dbErr, ok := castDbErr(err); ok {
		return dbErr.IsUniqueConstraintErr()
	}
	return false
}

// IsNotFoundErr returns whether the error is a Record Not Found error
func IsNotFoundErr(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// castDbErr casts the error as either a PostgreSQL or SQLite error
func castDbErr(err error) (dbErr, bool) {
	castSqliteErr := &sqlite3.Error{}
	if errors.As(err, castSqliteErr) {
		return &sqliteErr{err: castSqliteErr}, true
	}
	postgresErr := &pgconn.PgError{}
	if errors.As(err, &postgresErr) {
		return &pgErr{err: postgresErr}, true
	}
	return nil, false
}

// dbErr is the interface for database errors
type dbErr interface {
	IsUniqueConstraintErr() bool
}

// pgErr is the PostgreSQL implementation of dbErr
type pgErr struct {
	err *pgconn.PgError
}

// hasErrorCode returns whether the PostgreSQL error is of the given code
func (s *pgErr) hasErrorCode(code string) bool {
	return s.err.Code == code
}

// IsUniqueConstraintErr implements dbErr.IsUniqueConstraintErr
func (s *pgErr) IsUniqueConstraintErr() bool {
	return s.hasErrorCode(errPgUniqueViolation)
}

// sqliteErr is the SQLite implementation of dbErr
type sqliteErr struct {
	err *sqlite3.Error
}

// IsUniqueConstraintErr implements dbErr.IsUniqueConstraintErr
func (s *sqliteErr) IsUniqueConstraintErr() bool {
	if s.err.Code != sqlite3.ErrConstraint {
		return false
	}
	if s.err.ExtendedCode == sqlite3.ErrConstraintUnique {
		return true
	}
	if s.err.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
		return true
	}
	return false
}
