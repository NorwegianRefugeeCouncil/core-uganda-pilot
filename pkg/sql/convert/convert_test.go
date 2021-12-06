package convert

import (
	"github.com/nrc-no/core/pkg/api/types"
	sqlschema2 "github.com/nrc-no/core/pkg/sql/schema"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func Test_convertFormToSqlTable(t *testing.T) {
	tests := []struct {
		name    string
		form    string
		want    string
		wantErr bool
	}{
		{
			name: "simple",
			form: `
id: form
databaseId: database
fields:
- id: field
  required: true
  fieldType:
    text: { }
`,
			want: `
schema: database
name: form
columns:
- name: id
  dataType:
    varChar:
      length: 36
  constraints:
    - name: pk_database_form
      primaryKey: { }
- name: seq
  dataType:
    serial: { }
- name: created_at
  dataType:
    timestamp:
      timezone: WithoutTimezone
  default: NOW()
  constraints:
    - notNull: { }
- name: field
  dataType:
    varChar:
      length: 1024
  constraints:
    - notNull: { }
`,
		}, {
			name: "reference",
			form: `
id: form
databaseId: database
fields:
- id: field
  fieldType:
    reference:
      formId: otherform
      databaseId: otherdatabase
`,
			want: `
schema: database
name: form
columns:
- name: id
  dataType:
    varChar:
      length: 36
  constraints:
    - name: pk_database_form
      primaryKey: { }
- name: seq
  dataType:
    serial: { }
- name: created_at
  dataType:
    timestamp:
      timezone: WithoutTimezone
  default: NOW()
  constraints:
    - notNull: { }
- name: field
  dataType:
    varChar:
      length: 36
  constraints:
    - name: fkref__form_field__otherform_id
      reference:
        column: id
        onDelete: ActionCascade
        onUpdate: ActionCascade
        schema: otherdatabase
        table: otherform
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var form types.FormDefinition
			if err := yaml.Unmarshal([]byte(tt.form), &form); !assert.NoError(t, err) {
				return
			}

			var want sqlschema2.SQLTable
			if err := yaml.Unmarshal([]byte(tt.want), &want); !assert.NoError(t, err) {
				return
			}

			got, err := convertFormToSqlTable(&form, &types.FormDefinitionList{})
			if tt.wantErr {
				if !assert.Error(t, err) {
					return
				}
				return
			}
			if !tt.wantErr {
				if !assert.NoError(t, err) {
					return
				}
			}

			gotYaml, err := yaml.Marshal(got)
			if !assert.NoError(t, err) {
				return
			}

			assert.YAMLEq(t, tt.want, string(gotYaml))
		})
	}
}

func Test_expandSubForms(t *testing.T) {
	tests := []struct {
		name string
		args types.FormDefinition
		want []*types.FormDefinition
	}{
		{
			name: "simple",
			args: types.FormDefinition{
				DatabaseID: "db",
				ID:         "name",
				Fields:     []*types.FieldDefinition{},
			},
			want: []*types.FormDefinition{
				{
					DatabaseID: "db",
					ID:         "name",
					Fields:     []*types.FieldDefinition{},
				},
			},
		}, {
			name: "nested",
			args: types.FormDefinition{
				DatabaseID: "db",
				ID:         "root",
				Fields: []*types.FieldDefinition{
					{
						ID: "parent_field",
						FieldType: types.FieldType{
							SubForm: &types.FieldTypeSubForm{
								Fields: []*types.FieldDefinition{
									{
										ID: "child_field",
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
			want: []*types.FormDefinition{
				{
					DatabaseID: "db",
					ID:         "root",
				}, {
					DatabaseID: "db",
					ID:         "parent_field",
					Fields: []*types.FieldDefinition{
						{
							ID: "child_field",
							FieldType: types.FieldType{
								Text: &types.FieldTypeText{},
							},
						}, {
							ID:       "owner_id",
							Name:     "owner_id",
							Code:     "owner_id",
							Required: true,
							FieldType: types.FieldType{
								Reference: &types.FieldTypeReference{
									DatabaseID: "db",
									FormID:     "root",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expandSubForms(&tt.args)
			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
