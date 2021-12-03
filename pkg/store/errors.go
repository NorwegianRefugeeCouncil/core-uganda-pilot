package store

import (
	"errors"
	"github.com/jackc/pgconn"
	"github.com/mattn/go-sqlite3"
)

func IsUniqueConstraintErr(err error) bool {
	if dbErr, ok := castDbErr(err); ok {
		return dbErr.IsUniqueConstraintErr()
	}
	return false
}

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

type dbErr interface {
	IsUniqueConstraintErr() bool
}

type pgErr struct {
	err *pgconn.PgError
}

func (s *pgErr) IsErrCode(code string) bool {
	return s.err.Code == code
}

func (s *pgErr) IsUniqueConstraintErr() bool {
	return s.IsErrCode("23505")
}

type sqliteErr struct {
	err *sqlite3.Error
}

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
