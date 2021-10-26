package sqlconvert

import (
	"github.com/nrc-no/core/pkg/sqlschema"
	"github.com/nrc-no/core/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_convertFormToSqlTable(t *testing.T) {
	tests := []struct {
		name string
		args types.FormDefinition
		want sqlschema.SQLTable
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
			want: sqlschema.SQLTable{
				Schema: "database",
				Name:   `form`,
				Fields: []sqlschema.SQLField{
					sqlschema.NewSQLField("id").
						WithSerialDataType().
						WithPrimaryKeyConstraint("pk_database_form"),
					sqlschema.NewSQLField("field").
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
			want: sqlschema.SQLTable{
				Name:   "form",
				Schema: "database",
				Fields: []sqlschema.SQLField{
					sqlschema.NewSQLField("id").
						WithSerialDataType().
						WithPrimaryKeyConstraint("pk_database_form"),
					sqlschema.NewSQLField("field").
						WithIntDataType().
						WithNotNullConstraint().
						WithReferenceConstraint("fk__database_form_field__remote_other_id",
							"remote",
							"other",
							"id",
							sqlschema.ActionCascade,
							sqlschema.ActionCascade,
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
