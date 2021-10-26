package sqlconvert

import (
	sqlschema2 "github.com/nrc-no/core/pkg/sqlschema"
	types2 "github.com/nrc-no/core/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_convertFormToSqlTable(t *testing.T) {
	tests := []struct {
		name string
		args types2.FormDefinition
		want sqlschema2.SQLTable
	}{
		{
			name: "simple",
			args: types2.FormDefinition{
				DatabaseName: "database",
				Name:         "form",
				Fields: []types2.FieldDefinition{
					{
						Name: "field",
						FieldType: types2.FieldType{
							Text: &types2.FieldTypeText{},
						},
						Required: true,
					},
				},
			},
			want: sqlschema2.SQLTable{
				Schema: "database",
				Name:   `form`,
				Fields: []sqlschema2.SQLField{
					sqlschema2.NewSQLField("id").
						WithSerialDataType().
						WithPrimaryKeyConstraint("pk_database_form"),
					sqlschema2.NewSQLField("field").
						WithVarCharDataType(1024).
						WithNotNullConstraint(),
				},
			},
		}, {
			name: "reference",
			args: types2.FormDefinition{
				DatabaseName: "database",
				Name:         "form",
				Fields: []types2.FieldDefinition{
					{
						Name: "field",
						FieldType: types2.FieldType{
							Reference: &types2.FieldTypeReference{
								DatabaseName: "remote",
								FormName:     "other",
							},
						},
						Required: true,
					},
				},
			},
			want: sqlschema2.SQLTable{
				Name:   "form",
				Schema: "database",
				Fields: []sqlschema2.SQLField{
					sqlschema2.NewSQLField("id").
						WithSerialDataType().
						WithPrimaryKeyConstraint("pk_database_form"),
					sqlschema2.NewSQLField("field").
						WithIntDataType().
						WithNotNullConstraint().
						WithReferenceConstraint("fk__database_form_field__remote_other_id",
							"remote",
							"other",
							"id",
							sqlschema2.ActionCascade,
							sqlschema2.ActionCascade,
						),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, convertFormToSqlTable(tt.args))
		})
	}
}

func Test_expandSubForms(t *testing.T) {
	tests := []struct {
		name string
		args types2.FormDefinition
		want []types2.FormDefinition
	}{
		{
			name: "simple",
			args: types2.FormDefinition{
				DatabaseName: "db",
				Name:         "name",
				Fields:       []types2.FieldDefinition{},
			},
			want: []types2.FormDefinition{
				{
					DatabaseName: "db",
					Name:         "name",
					Fields:       []types2.FieldDefinition{},
				},
			},
		}, {
			name: "nested",
			args: types2.FormDefinition{
				DatabaseName: "db",
				Name:         "root",
				Fields: []types2.FieldDefinition{
					{
						Name: "parent_field",
						FieldType: types2.FieldType{
							SubForm: &types2.FieldTypeSubForm{
								Fields: []types2.FieldDefinition{
									{
										Name: "child_field",
										FieldType: types2.FieldType{
											Text: &types2.FieldTypeText{},
										},
									},
								},
							},
						},
					},
				},
			},
			want: []types2.FormDefinition{
				{
					DatabaseName: "db",
					Name:         "root",
					Fields: []types2.FieldDefinition{
						{
							Name: "parent_field",
							FieldType: types2.FieldType{
								SubForm: &types2.FieldTypeSubForm{
									Fields: []types2.FieldDefinition{
										{
											Name: "child_field",
											FieldType: types2.FieldType{
												Text: &types2.FieldTypeText{},
											},
										},
									},
								},
							},
						},
					},
				}, {
					DatabaseName: "db",
					Name:         "root_parent_fields",
					Fields: []types2.FieldDefinition{
						{
							Name: "child_field",
							FieldType: types2.FieldType{
								Text: &types2.FieldTypeText{},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, expandSubForms(tt.args))
		})
	}
}
