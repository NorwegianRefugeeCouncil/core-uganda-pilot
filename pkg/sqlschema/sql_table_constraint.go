package sqlschema

import (
	"fmt"
	"github.com/lib/pq"
	"strings"
)

type SQLTableConstraint struct {
	Name       string
	Check      *SQLTableConstraintCheck
	Unique     *SQLTableConstraintUnique
	PrimaryKey *SQLTableConstraintPrimaryKey
	ForeignKey *SQLTableConstraintForeignKey
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
	Expression string
}

func (c SQLTableConstraintCheck) DDL() DDL {
	return NewDDL("check (?)", c.Expression)
}

type SQLTableConstraintUnique struct {
	ColumnNames []string
}

func (u SQLTableConstraintUnique) DDL() DDL {
	return NewDDL(fmt.Sprintf("unique (%s)", strings.Join(u.ColumnNames, ", ")))
}

type SQLTableConstraintPrimaryKey struct {
	ColumnNames []string
}

func (k SQLTableConstraintPrimaryKey) DDL() DDL {
	var args []interface{}
	for _, name := range k.ColumnNames {
		args = append(args, name)
	}
	return NewDDL(fmt.Sprintf("primary key (%s)", writeParamPlaceholders(len(k.ColumnNames))), args...)
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
	ColumnNames       []string
	ReferencesTable   string
	ReferencesSchema  string
	ReferencesColumns []string
	MatchType         SQLMatchType
	OnDelete          SQLForeignKeyAction
	OnUpdate          SQLForeignKeyAction
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
