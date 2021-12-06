package schema

import (
	"github.com/lib/pq"
)

type SQLTableName string

type SQLTable struct {
	Schema      string              `json:"schema,omitempty" yaml:"schema,omitempty"`
	Name        string              `json:"name,omitempty" yaml:"name,omitempty"`
	Columns     SQLColumns          `json:"columns,omitempty" yaml:"columns,omitempty"`
	Constraints SQLTableConstraints `json:"constraints,omitempty" yaml:"constraints,omitempty"`
}

type SQLTables []SQLTable

func (s SQLTables) GetTable(name string) (SQLTable, error) {
	for _, table := range s {
		if table.Name == name {
			return table, nil
		}
	}
	return SQLTable{}, newTableNotFoundErr(name)
}

func (s SQLTables) ReplaceTable(table SQLTable) (SQLTables, error) {
	for i := 0; i < len(s); i++ {
		current := s[i]
		if current.Name == table.Name {
			s[i] = table
			return s, nil
		}
	}
	return nil, newTableNotFoundErr(table.Name)
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

func (s SQLTable) DDL() DDL {
	ddl := DDL{}

	ddl = ddl.WriteF(`create table %s.%s`,
		pq.QuoteIdentifier(s.Schema),
		pq.QuoteIdentifier(s.Name),
	)

	if len(s.Columns) == 0 && len(s.Constraints) == 0 {
		return ddl.Write("();")
	}

	ddl = ddl.Write("\n(\n")

	var statements []DDLGenerator

	for _, column := range s.Columns {
		statements = append(statements, column.DDL().WriteBefore("  "))
	}
	for _, constraint := range s.Constraints {
		statements = append(statements, NewDDL("  ").Merge(constraint.DDL()))
	}

	ddl = ddl.
		MergeAll(",\n", statements...).
		WriteString("\n);")

	return ddl
}
