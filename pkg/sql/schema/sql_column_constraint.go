package schema

import (
	"fmt"
	"github.com/lib/pq"
)

type NotNullSQLColumnConstraint struct{}

type NullSQLColumnConstraint struct{}

type CheckSQLColumnConstraint struct {
	Expression string `json:"expression" yaml:"expression"`
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
	Schema    string              `json:"schema,omitempty" yaml:"schema,omitempty"`
	Table     string              `json:"table,omitempty" yaml:"table,omitempty"`
	Column    string              `json:"column,omitempty" yaml:"column,omitempty"`
	OnDelete  SQLForeignKeyAction `json:"onDelete,omitempty" yaml:"onDelete,omitempty"`
	OnUpdate  SQLForeignKeyAction `json:"onUpdate,omitempty" yaml:"onUpdate,omitempty"`
	MatchType SQLMatchType        `json:"matchType,omitempty" yaml:"matchType,omitempty"`
}

type SQLColumnConstraint struct {
	Name       string                         `json:"name,omitempty" yaml:"name,omitempty"`
	NotNull    *NotNullSQLColumnConstraint    `json:"notNull,omitempty" yaml:"notNull,omitempty"`
	Null       *NullSQLColumnConstraint       `json:"null,omitempty" yaml:"null,omitempty"`
	Check      *CheckSQLColumnConstraint      `json:"check,omitempty" yaml:"check,omitempty"`
	Unique     *UniqueSQLColumnConstraint     `json:"unique,omitempty" yaml:"unique,omitempty"`
	PrimaryKey *PrimaryKeySQLColumnConstraint `json:"primaryKey,omitempty" yaml:"primaryKey,omitempty"`
	Reference  *ReferenceSQLColumnConstraint  `json:"reference,omitempty" yaml:"reference,omitempty"`
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
