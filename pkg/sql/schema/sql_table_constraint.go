package schema

import (
	"fmt"
	"github.com/lib/pq"
	"strings"
)

type SQLTableConstraint struct {
	Name       string                        `json:"name,omitempty" yaml:"name,omitempty"`
	Check      *SQLTableConstraintCheck      `json:"check,omitempty" yaml:"check,omitempty"`
	Unique     *SQLTableConstraintUnique     `json:"unique,omitempty" yaml:"unique,omitempty"`
	PrimaryKey *SQLTableConstraintPrimaryKey `json:"primaryKey,omitempty" yaml:"primaryKey,omitempty"`
	ForeignKey *SQLTableConstraintForeignKey `json:"foreignKey,omitempty" yaml:"foreignKey,omitempty"`
}

func NewCheckTableConstraint(name string, expression string) SQLTableConstraint {
	return SQLTableConstraint{
		Name: name,
		Check: &SQLTableConstraintCheck{
			Expression: expression,
		},
	}
}

func NewPrimaryKeyTableConstraint(name string, columnNames ...string) SQLTableConstraint {
	return SQLTableConstraint{
		Name: name,
		PrimaryKey: &SQLTableConstraintPrimaryKey{
			ColumnNames: columnNames,
		},
	}
}

func NewUniqueTableConstraint(name string, columnNames ...string) SQLTableConstraint {
	return SQLTableConstraint{
		Name: name,
		Unique: &SQLTableConstraintUnique{
			ColumnNames: columnNames,
		},
	}
}

func NewForeignKeyTableConstraint(
	name string,
	columnNames []string,
	referencedTable string,
	referencedColumns []string,
	onDelete SQLForeignKeyAction,
	onUpdate SQLForeignKeyAction,
) SQLTableConstraint {
	return SQLTableConstraint{
		Name: name,
		ForeignKey: &SQLTableConstraintForeignKey{
			ColumnNames:       columnNames,
			ReferencesTable:   referencedTable,
			ReferencesColumns: referencedColumns,
			MatchType:         "",
			OnDelete:          onDelete,
			OnUpdate:          onUpdate,
		},
	}
}

func (c SQLTableConstraint) DDL() DDL {
	ddl := DDL{}
	ddl.WriteF("constraint %s", pq.QuoteIdentifier(c.Name))
	if c.Unique != nil {
		ddl = ddl.WriteString(" ").Merge(c.Unique)
	}
	if c.Check != nil {
		ddl = ddl.WriteString(" ").Merge(c.Check)
	}
	if c.PrimaryKey != nil {
		ddl = ddl.WriteString(" ").Merge(c.PrimaryKey)
	}
	if c.ForeignKey != nil {
		ddl = ddl.WriteString(" ").Merge(c.ForeignKey)
	}
	return ddl
}

type SQLTableConstraintCheck struct {
	Expression string `json:"expression,omitempty" yaml:"expression,omitempty"`
}

func (c SQLTableConstraintCheck) DDL() DDL {
	return NewDDL("check (?)", c.Expression)
}

type SQLTableConstraintUnique struct {
	ColumnNames []string `json:"columnNames,omitempty" yaml:"columnNames,omitempty"`
}

func (u SQLTableConstraintUnique) DDL() DDL {
	var columnNames []string
	for _, columnName := range u.ColumnNames {
		columnNames = append(columnNames, pq.QuoteIdentifier(columnName))
	}
	return NewDDL(fmt.Sprintf("unique (%s)", strings.Join(columnNames, ", ")))
}

type SQLTableConstraintPrimaryKey struct {
	ColumnNames []string `json:"columnNames,omitempty" yaml:"columnNames,omitempty"`
}

func (k SQLTableConstraintPrimaryKey) DDL() DDL {
	var columnNames []string
	for _, columnName := range k.ColumnNames {
		columnNames = append(columnNames, pq.QuoteIdentifier(columnName))
	}
	return NewDDL(fmt.Sprintf("primary key (%s)", strings.Join(columnNames, ",")))
}

type SQLForeignKeyAction string

const (
	ActionSetNull    SQLForeignKeyAction = "ActionSetNull"
	ActionSetDefault SQLForeignKeyAction = "ActionSetDefault"
	ActionRestrict   SQLForeignKeyAction = "ActionRestrict"
	ActionNoAction   SQLForeignKeyAction = "ActionNoAction"
	ActionCascade    SQLForeignKeyAction = "ActionCascade"
)

type SQLTableConstraintForeignKey struct {
	ColumnNames       []string            `json:"columnNames,omitempty" yaml:"columnNames,omitempty"`
	ReferencesTable   string              `json:"referencesTable,omitempty" yaml:"referencesTable,omitempty"`
	ReferencesSchema  string              `json:"referencesSchema,omitempty" yaml:"referencesSchema,omitempty"`
	ReferencesColumns []string            `json:"referencesColumns,omitempty" yaml:"referencesColumns,omitempty"`
	MatchType         SQLMatchType        `json:"matchType,omitempty" yaml:"matchType,omitempty"`
	OnDelete          SQLForeignKeyAction `json:"onDelete,omitempty" yaml:"onDelete,omitempty"`
	OnUpdate          SQLForeignKeyAction `json:"onUpdate,omitempty" yaml:"onUpdate,omitempty"`
}

func (k SQLTableConstraintForeignKey) DDL() DDL {
	ddl := DDL{}

	var colNames []string
	for _, colName := range k.ColumnNames {
		colNames = append(colNames, pq.QuoteIdentifier(colName))
	}

	var refColNames []string
	for _, colName := range k.ReferencesColumns {
		refColNames = append(refColNames, pq.QuoteIdentifier(colName))
	}

	ddl = ddl.
		WriteString("foreign key (").
		Write(strings.Join(colNames, ",")).
		WriteF(") references %s.%s (", pq.QuoteIdentifier(k.ReferencesSchema), pq.QuoteIdentifier(k.ReferencesTable)).
		Write(strings.Join(refColNames, ",")).
		Write(")")

	onUpdateActionDDL := getFKActionDDL(k.OnUpdate)
	if len(onUpdateActionDDL) > 0 {
		ddl = ddl.WriteF(" on update %s", onUpdateActionDDL)
	}

	onDeleteActionDDL := getFKActionDDL(k.OnDelete)
	if len(onDeleteActionDDL) > 0 {
		ddl = ddl.WriteF(" on delete %s", onDeleteActionDDL)
	}

	return ddl
}
