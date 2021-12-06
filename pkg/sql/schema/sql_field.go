package schema

import (
	"github.com/lib/pq"
)

type SQLColumn struct {
	Name        string                `json:"name" yaml:"name"`
	DataType    SQLDataType           `json:"dataType" yaml:"dataType"`
	Collate     string                `json:"collate,omitempty" yaml:"collate,omitempty"`
	Default     string                `json:"default,omitempty" yaml:"default,omitempty"`
	Options     []string              `json:"options,omitempty" yaml:"options,omitempty"`
	Constraints []SQLColumnConstraint `json:"constraints,omitempty" yaml:"constraints,omitempty"`
	Comment     string                `json:"comment,omitempty" yaml:"comment,omitempty"`
}

type SQLColumns []SQLColumn

func (s SQLColumns) GetColumn(name string) (SQLColumn, error) {
	for _, field := range s {
		if field.Name == name {
			return field, nil
		}
	}
	return SQLColumn{}, newColumnNotFoundErr(name)
}

func (s SQLColumn) DDL() DDL {
	ddl := DDL{}
	ddl = ddl.
		WriteF("%s ", pq.QuoteIdentifier(s.Name)).
		Merge(s.DataType)

	for _, constraint := range s.Constraints {
		ddl = ddl.MergeAll("", NewDDL(" "), constraint)
	}

	if len(s.Collate) > 0 {
		ddl = ddl.
			WriteString(" collate ?").
			WriteArgs(s.Collate)
	}

	if len(s.Default) > 0 {
		ddl = ddl.
			WriteString(" default " + s.Default)
	}

	return ddl
}
