package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSQLTable_DDL(t *testing.T) {
	type args struct {
		Name        string
		Fields      []SQLField
		Constraints []SQLTableConstraint
	}
	tests := []struct {
		name   string
		fields args
		want   DDL
	}{
		{
			name:   "empty",
			fields: args{Name: "empty"},
			want:   NewDDL(`create table "public"."empty";`),
		}, {
			name: "single field",
			fields: args{Name: "singleField", Fields: []SQLField{
				{
					Name:     "field",
					DataType: SQLDataType{Int: &SQLDataTypeInt{}},
				},
			}},
			want: NewDDL(`create table "public"."singleField"
(
  "field" int
);`),
		}, {
			name: "multi field",
			fields: args{Name: "multiField", Fields: []SQLField{
				{
					Name:     "field1",
					DataType: SQLDataType{Int: &SQLDataTypeInt{}},
				}, {
					Name:     "field2",
					DataType: SQLDataType{Int: &SQLDataTypeInt{}},
				},
			}},
			want: NewDDL(`create table "public"."multiField"
(
  "field1" int,
  "field2" int
);`),
		}, {
			name: "multi field and constraints",
			fields: args{Name: "multiFieldConstraint", Fields: []SQLField{
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
			want: NewDDL(`create table "public"."multiFieldConstraint"
(
  "field1" int,
  "field2" int,
  constraint "pk_field1" primary key ("field1"),
  constraint "uq_field2" unique ("field2")
);`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SQLTable{
				Schema:      "public",
				Name:        tt.fields.Name,
				Fields:      tt.fields.Fields,
				Constraints: tt.fields.Constraints,
			}
			assert.Equal(t, tt.want, s.DDL())
		})
	}
}
