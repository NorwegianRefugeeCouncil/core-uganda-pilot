package convert

import (
	"github.com/nrc-no/core/pkg/api/types"
	sqlschema2 "github.com/nrc-no/core/pkg/sql/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_convertFormToSqlTable(t *testing.T) {
	tests := []struct {
		name string
		args types.FormDefinition
		want sqlschema2.SQLTable
	}{
		{
			name: "simple",
			args: types.FormDefinition{
				DatabaseName: "database",
				Name:         "form",
				Fields: []types.FieldDefinition{
					{
						Name: "field",
						FieldType: types.FieldType{
							Text: &types.FieldTypeText{},
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
			args: types.FormDefinition{
				DatabaseName: "database",
				Name:         "form",
				Fields: []types.FieldDefinition{
					{
						Name: "field",
						FieldType: types.FieldType{
							Reference: &types.FieldTypeReference{
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
		args types.FormDefinition
		want []types.FormDefinition
	}{
		{
			name: "simple",
			args: types.FormDefinition{
				DatabaseName: "db",
				Name:         "name",
				Fields:       []types.FieldDefinition{},
			},
			want: []types.FormDefinition{
				{
					DatabaseName: "db",
					Name:         "name",
					Fields:       []types.FieldDefinition{},
				},
			},
		}, {
			name: "nested",
			args: types.FormDefinition{
				DatabaseName: "db",
				Name:         "root",
				Fields: []types.FieldDefinition{
					{
						Name: "parent_field",
						FieldType: types.FieldType{
							SubForm: &types.FieldTypeSubForm{
								Fields: []types.FieldDefinition{
									{
										Name: "child_field",
										FieldType: types.FieldType{
											Text: &types.FieldTypeText{},
										},
									},
								},
							},
						},
					},
				},
			},
			want: []types.FormDefinition{
				{
					DatabaseName: "db",
					Name:         "root",
					Fields: []types.FieldDefinition{
						{
							Name: "parent_field",
							FieldType: types.FieldType{
								SubForm: &types.FieldTypeSubForm{
									Fields: []types.FieldDefinition{
										{
											Name: "child_field",
											FieldType: types.FieldType{
												Text: &types.FieldTypeText{},
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
					Fields: []types.FieldDefinition{
						{
							Name: "child_field",
							FieldType: types.FieldType{
								Text: &types.FieldTypeText{},
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
