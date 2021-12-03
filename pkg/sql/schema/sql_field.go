package schema

import (
	"github.com/lib/pq"
)

type SQLField struct {
	Name        string                `json:"name" yaml:"name"`
	DataType    SQLDataType           `json:"dataType" yaml:"dataType"`
	Collate     string                `json:"collate,omitempty" yaml:"collate,omitempty"`
	Default     string                `json:"default,omitempty" yaml:"default,omitempty"`
	Options     []string              `json:"options,omitempty" yaml:"options,omitempty"`
	Constraints []SQLColumnConstraint `json:"constraints,omitempty" yaml:"constraints,omitempty"`
	Comment     string                `json:"comment,omitempty" yaml:"comment,omitempty"`
}

func NewSQLField(name string) SQLField {
	return SQLField{Name: name}
}

func (s SQLField) WithCollate(collate string) SQLField {
	s.Collate = collate
	return s
}

func (s SQLField) WithSerialDataType() SQLField {
	s.DataType.Serial = &SQLDataTypeSerial{}
	return s
}

func (s SQLField) WithVarCharDataType(length int) SQLField {
	s.DataType.VarChar = &SQLDataTypeVarChar{Length: length}
	return s
}

func (s SQLField) WithTimeStampDataType(tzMode *SQLDataTypeTimestampTZMode) SQLField {
	s.DataType.Timestamp = &SQLDataTypeTimestamp{Timezone: tzMode}
	return s
}

func (s SQLField) WithIntDataType() SQLField {
	s.DataType.Int = &SQLDataTypeInt{}
	return s
}

func (s SQLField) WithConstraints(constraints ...SQLColumnConstraint) SQLField {
	s.Constraints = append(s.Constraints, constraints...)
	return s
}

func (s SQLField) WithPrimaryKeyConstraint(name string) SQLField {
	s.Constraints = append(s.Constraints, SQLColumnConstraint{
		Name:       name,
		PrimaryKey: &PrimaryKeySQLColumnConstraint{},
	})
	return s
}

func (s SQLField) WithNotNullConstraint() SQLField {
	s.Constraints = append(s.Constraints, SQLColumnConstraint{
		NotNull: &NotNullSQLColumnConstraint{},
	})
	return s
}

func (s SQLField) WithReferenceConstraint(
	name string,
	schema string,
	table string,
	column string,
	onDelete SQLForeignKeyAction,
	onUpdate SQLForeignKeyAction,
) SQLField {
	s.Constraints = append(s.Constraints, SQLColumnConstraint{
		Name: name,
		Reference: &ReferenceSQLColumnConstraint{
			Schema:   schema,
			Table:    table,
			Column:   column,
			OnDelete: onDelete,
			OnUpdate: onUpdate,
		},
	})
	return s
}

func (s SQLField) DDL() DDL {
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
