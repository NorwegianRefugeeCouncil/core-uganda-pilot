package schema

import (
	"github.com/lib/pq"
)

type SQLTableName string

type SQLTable struct {
	Schema      string
	Name        string
	Fields      []SQLField
	Constraints []SQLTableConstraint
}

func (s SQLTable) WithField(field SQLField) SQLTable {
	s.Fields = append(s.Fields, field)
	return s
}

func (s SQLTable) WithConstraint(constraint SQLTableConstraint) SQLTable {
	s.Constraints = append(s.Constraints, constraint)
	return s
}

func (s SQLTable) WithCheckConstraint(name, expression string) SQLTable {
	s.Constraints = append(s.Constraints, NewCheckTableConstraint(name, expression))
	return s
}

func (s SQLTable) WithUniqueConstraint(name string, columnNames ...string) SQLTable {
	s.Constraints = append(s.Constraints, NewUniqueTableConstraint(name, columnNames...))
	return s
}

func (s SQLTable) WithPrimaryKeyConstraint(name string, columnNames ...string) SQLTable {
	s.Constraints = append(s.Constraints, NewPrimaryKeyTableConstraint(name, columnNames...))
	return s
}

func (s SQLTable) WithForeignKeyConstraint(
	name string,
	columnNames []string,
	referencedTable string,
	referencedColumns []string,
	onDelete SQLForeignKeyAction,
	onUpdate SQLForeignKeyAction,
) SQLTable {
	s.Constraints = append(s.Constraints, NewForeignKeyTableConstraint(name, columnNames, referencedTable, referencedColumns, onDelete, onUpdate))
	return s
}

func NewSQLTable(schema, name string) SQLTable {
	return SQLTable{Schema: schema, Name: name}
}

func (s SQLTable) DDL() DDL {
	ddl := DDL{}

	ddl = ddl.WriteF(`create table %s.%s`,
		pq.QuoteIdentifier(s.Schema),
		pq.QuoteIdentifier(s.Name),
	)

	if len(s.Fields) == 0 && len(s.Constraints) == 0 {
		return ddl.Write(";")
	}

	ddl = ddl.Write("\n(\n")

	var statements []DDLGenerator

	for _, field := range s.Fields {
		statements = append(statements, field.DDL().WriteBefore("  "))
	}
	for _, constraint := range s.Constraints {
		statements = append(statements, constraint.DDL().WriteBeforeF("  constraint %s", pq.QuoteIdentifier(constraint.Name)))
	}

	ddl = ddl.
		MergeAll(",\n", statements...).
		WriteString("\n);")

	return ddl
}
