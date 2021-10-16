package store

import (
	"github.com/nrc-no/core/pkg/bla/types"
	"github.com/nrc-no/core/pkg/pointers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mapToFormDefinitions(t *testing.T) {
	tests := []struct {
		name    string
		forms   []*Form
		fields  []*Field
		want    []*types.FormDefinition
		wantErr bool
	}{
		{
			name: "simple",
			forms: []*Form{
				{
					ID:         "formId",
					DatabaseID: "db",
					FolderID:   "folder",
					Name:       "formName",
				},
			},
			fields: []*Field{
				{
					ID:         "field1",
					DatabaseID: "db",
					FormID:     "formId",
					RootFormID: "formId",
					Name:       "fieldName",
					Type:       FieldTypeText,
				},
			},
			want: []*types.FormDefinition{
				{
					ID:         "formId",
					DatabaseID: "db",
					FolderID:   "folder",
					Name:       "formName",
					Fields: []*types.FieldDefinition{
						{
							ID:   "field1",
							Name: "fieldName",
							FieldType: types.FieldType{
								Text: &types.FieldTypeText{},
							},
						},
					},
				},
			},
		}, {
			name: "with subform",
			forms: []*Form{
				{
					ID:         "formId1",
					DatabaseID: "db",
					FolderID:   "folder",
					Name:       "form1",
				}, {
					ID:         "formId2",
					DatabaseID: "db",
					FolderID:   "folder",
					RootID:     pointers.String("formId1"),
					ParentID:   pointers.String("formId1"),
					Name:       "form2",
				},
			},
			fields: []*Field{
				{
					ID:         "field1",
					DatabaseID: "db",
					FormID:     "formId1",
					RootFormID: "formId1",
					SubFormID:  pointers.String("formId2"),
					Name:       "sub",
					Type:       FieldTypeSubForm,
				},
			},
			want: []*types.FormDefinition{
				{
					ID:         "formId1",
					DatabaseID: "db",
					FolderID:   "folder",
					Name:       "form1",
					Fields: []*types.FieldDefinition{
						{
							ID:   "field1",
							Name: "sub",
							FieldType: types.FieldType{
								SubForm: &types.FieldTypeSubForm{
									ID:     "formId2",
									Name:   "form2",
									Fields: []*types.FieldDefinition{},
								},
							},
						},
					},
				},
			},
		}, {
			name: "with multiple subforms",
			forms: []*Form{
				{
					ID:         "formId1",
					DatabaseID: "db",
					FolderID:   "folder",
					Name:       "form1",
				}, {
					ID:         "formId2",
					DatabaseID: "db",
					FolderID:   "folder",
					RootID:     pointers.String("formId1"),
					ParentID:   pointers.String("formId1"),
					Name:       "form2",
				}, {
					ID:         "formId3",
					DatabaseID: "db",
					FolderID:   "folder",
					RootID:     pointers.String("formId1"),
					ParentID:   pointers.String("formId1"),
					Name:       "form3",
				},
			},
			fields: []*Field{
				{
					ID:         "fieldId1",
					DatabaseID: "db",
					FormID:     "formId1",
					RootFormID: "formId1",
					SubFormID:  pointers.String("formId2"),
					Name:       "field1",
					Type:       FieldTypeSubForm,
				}, {
					ID:         "fieldId2",
					DatabaseID: "db",
					FormID:     "formId1",
					RootFormID: "formId1",
					SubFormID:  pointers.String("formId3"),
					Name:       "field2",
					Type:       FieldTypeSubForm,
				},
			},
			want: []*types.FormDefinition{
				{
					ID:         "formId1",
					DatabaseID: "db",
					FolderID:   "folder",
					Name:       "form1",
					Fields: []*types.FieldDefinition{
						{
							ID:   "fieldId1",
							Name: "field1",
							FieldType: types.FieldType{
								SubForm: &types.FieldTypeSubForm{
									ID:     "formId2",
									Name:   "form2",
									Fields: []*types.FieldDefinition{},
								},
							},
						}, {
							ID:   "fieldId2",
							Name: "field2",
							FieldType: types.FieldType{
								SubForm: &types.FieldTypeSubForm{
									ID:     "formId3",
									Name:   "form3",
									Fields: []*types.FieldDefinition{},
								},
							},
						},
					},
				},
			},
		}, {
			name: "with nested subform",
			forms: []*Form{
				{
					ID:         "formId1",
					DatabaseID: "db",
					FolderID:   "folder",
					Name:       "form1",
				}, {
					ID:         "formId2",
					DatabaseID: "db",
					FolderID:   "folder",
					RootID:     pointers.String("formId1"),
					ParentID:   pointers.String("formId1"),
					Name:       "form2",
				}, {
					ID:         "formId3",
					DatabaseID: "db",
					FolderID:   "folder",
					RootID:     pointers.String("formId1"),
					ParentID:   pointers.String("formId2"),
					Name:       "form3",
				},
			},
			fields: []*Field{
				{
					ID:         "field1",
					DatabaseID: "db",
					FormID:     "formId1",
					RootFormID: "formId1",
					SubFormID:  pointers.String("formId2"),
					Name:       "field1Name",
					Type:       FieldTypeSubForm,
				}, {
					ID:         "field2",
					DatabaseID: "db",
					FormID:     "formId2",
					RootFormID: "formId1",
					SubFormID:  pointers.String("formId3"),
					Name:       "field2Name",
					Type:       FieldTypeSubForm,
				},
			},
			want: []*types.FormDefinition{
				{
					ID:         "formId1",
					DatabaseID: "db",
					FolderID:   "folder",
					Name:       "form1",
					Fields: []*types.FieldDefinition{
						{
							ID:   "field1",
							Name: "field1Name",
							FieldType: types.FieldType{
								SubForm: &types.FieldTypeSubForm{
									ID:   "formId2",
									Name: "form2",
									Fields: []*types.FieldDefinition{
										{
											ID:   "field2",
											Name: "field2Name",
											FieldType: types.FieldType{
												SubForm: &types.FieldTypeSubForm{
													ID:     "formId3",
													Name:   "form3",
													Fields: []*types.FieldDefinition{},
												},
											},
										},
									},
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
			gotFd, err := mapToFormDefinitions(tt.forms, tt.fields)
			if tt.wantErr && !assert.Error(t, err) {
				return
			}
			if !tt.wantErr && !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tt.want, gotFd)

			frms, flds, err := mapToFormFields(gotFd[0])
			for i, form := range tt.forms {
				assert.Equal(t, form, frms[i])
			}
			for i, fld := range tt.fields {
				assert.Equal(t, fld, flds[i])
			}

		})
	}
}
