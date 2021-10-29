package schema

import (
	"fmt"
	"github.com/lib/pq"
)

type NotNullSQLColumnConstraint struct{}

type NullSQLColumnConstraint struct{}

type CheckSQLColumnConstraint struct {
	Expression string
}

type UniqueSQLColumnConstraint struct{}

type PrimaryKeySQLColumnConstraint struct{}

type SQLMatchType string

const (
	MatchFull    SQLMatchType = "Full"
	MatchPartial SQLMatchType = "Partial"
	MatchSimple  SQLMatchType = "Simple"
)

type ReferenceSQLColumnConstraint struct {
	Schema    string
	Table     string
	Column    string
	OnDelete  SQLForeignKeyAction
	OnUpdate  SQLForeignKeyAction
	MatchType SQLMatchType
}

type SQLColumnConstraint struct {
	Name       string
	NotNull    *NotNullSQLColumnConstraint
	Null       *NullSQLColumnConstraint
	Check      *CheckSQLColumnConstraint
	Unique     *UniqueSQLColumnConstraint
	PrimaryKey *PrimaryKeySQLColumnConstraint
	Reference  *ReferenceSQLColumnConstraint
}

func (c SQLColumnConstraint) DDL() DDL {
	if c.NotNull != nil {
		return NewDDL("not null")
	}
	if c.Null != nil {
		return NewDDL("null")
	}
	if c.PrimaryKey != nil {
		return NewDDL("primary key")
	}
	if c.Unique != nil {
		return NewDDL("unique")
	}
	if c.Reference != nil {
		return NewDDL(fmt.Sprintf("references %s.%s (%s)",
			pq.QuoteIdentifier(c.Reference.Schema),
			pq.QuoteIdentifier(c.Reference.Table),
			pq.QuoteIdentifier(c.Reference.Column),
		))
	}
	return EmptyDDL()
}
