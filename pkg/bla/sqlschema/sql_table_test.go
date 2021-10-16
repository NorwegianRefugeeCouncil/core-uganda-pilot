package sqlschema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSQLTable_DDL(t *testing.T) {
	type fields struct {
		Name        string
		Fields      []SQLField
		Constraints []SQLTableConstraint
	}
	tests := []struct {
		name   string
		fields fields
		want   DDL
	}{
		{
			name:   "empty",
			fields: fields{Name: "empty"},
			want:   NewDDL(`create table ?;`, "empty"),
		}, {
			name: "single field",
			fields: fields{Name: "singleField", Fields: []SQLField{
				{
					Name:     "field",
					DataType: SQLDataType{Int: &SQLDataTypeInt{}},
				},
			}},
			want: NewDDL(`create table ?
(
  ? int
);`, "singleField", "field"),
		}, {
			name: "multi field",
			fields: fields{Name: "multiField", Fields: []SQLField{
				{
					Name:     "field1",
					DataType: SQLDataType{Int: &SQLDataTypeInt{}},
				}, {
					Name:     "field2",
					DataType: SQLDataType{Int: &SQLDataTypeInt{}},
				},
			}},
			want: NewDDL(`create table ?
(
  ? int,
  ? int
);`, "multiField", "field1", "field2"),
		}, {
			name: "multi field and constraints",
			fields: fields{Name: "multiFieldConstraint", Fields: []SQLField{
				{
					Name:     "field1",
					DataType: SQLDataType{Int: &SQLDataTypeInt{}},
				}, {
					Name:     "field2",
					DataType: SQLDataType{Int: &SQLDataTypeInt{}},
				},
			},
				Constraints: []SQLTableConstraint{
					{
						Name: "pk_field1",
						PrimaryKey: &SQLTableConstraintPrimaryKey{
							ColumnNames: []string{
								"field1",
							},
						},
					}, {
						Name: "uq_field2",
						Unique: &SQLTableConstraintUnique{
							ColumnNames: []string{
								"field2",
							},
						},
					},
				}},
			want: NewDDL(`create table ?
(
  ? int,
  ? int,
  constraint ? primary key (?),
  constraint ? unique (?)
);`,
				"multiFieldConstraint",
				"field1",
				"field2",
				"pk_field1",
				"field1",
				"uq_field2",
				"field2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SQLTable{
				Name:        tt.fields.Name,
				Fields:      tt.fields.Fields,
				Constraints: tt.fields.Constraints,
			}
			assert.Equal(t, tt.want, s.DDL())
		})
	}
}
