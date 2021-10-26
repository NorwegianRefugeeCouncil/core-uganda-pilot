package store

import (
	"github.com/nrc-no/core/pkg/pointers"
	"github.com/nrc-no/core/pkg/types"
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
					RootID:     "formId",
					ParentID:   "formId",
					ID:         "formId",
					DatabaseID: "db",
					FolderID:   pointers.String("folder"),
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
					Type:       types.FieldKindText,
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
					RootID:     "formId1",
					ParentID:   "formId1",
					ID:         "formId1",
					DatabaseID: "db",
					FolderID:   pointers.String("folder"),
					Name:       "form1",
				}, {
					ID:         "formId2",
					DatabaseID: "db",
					FolderID:   pointers.String("folder"),
					RootID:     "formId1",
					ParentID:   "formId1",
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
					Type:       types.FieldKindSubForm,
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
					RootID:     "formId1",
					ParentID:   "formId1",
					ID:         "formId1",
					DatabaseID: "db",
					FolderID:   pointers.String("folder"),
					Name:       "form1",
				}, {
					ID:         "formId2",
					DatabaseID: "db",
					FolderID:   pointers.String("folder"),
					RootID:     "formId1",
					ParentID:   "formId1",
					Name:       "form2",
				}, {
					ID:         "formId3",
					DatabaseID: "db",
					FolderID:   pointers.String("folder"),
					RootID:     "formId1",
					ParentID:   "formId1",
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
					Type:       types.FieldKindSubForm,
				}, {
					ID:         "fieldId2",
					DatabaseID: "db",
					FormID:     "formId1",
					RootFormID: "formId1",
					SubFormID:  pointers.String("formId3"),
					Name:       "field2",
					Type:       types.FieldKindSubForm,
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
					RootID:     "formId1",
					ParentID:   "formId1",
					DatabaseID: "db",
					FolderID:   pointers.String("folder"),
					Name:       "form1",
				}, {
					ID:         "formId2",
					DatabaseID: "db",
					FolderID:   pointers.String("folder"),
					RootID:     "formId1",
					ParentID:   "formId1",
					Name:       "form2",
				}, {
					ID:         "formId3",
					DatabaseID: "db",
					FolderID:   pointers.String("folder"),
					RootID:     "formId1",
					ParentID:   "formId2",
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
					Type:       types.FieldKindSubForm,
				}, {
					ID:         "field2",
					DatabaseID: "db",
					FormID:     "formId2",
					RootFormID: "formId1",
					SubFormID:  pointers.String("formId3"),
					Name:       "field2Name",
					Type:       types.FieldKindSubForm,
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