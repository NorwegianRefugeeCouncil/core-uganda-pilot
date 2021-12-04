package store

import (
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/utils/pointers"
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
					RootOwnerID: "formId",
					OwnerID:     "formId",
					ID:          "formId",
					DatabaseID:  "db",
					FolderID:    pointers.String("folder"),
					Name:        "formName",
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
					RootOwnerID: "formId1",
					OwnerID:     "formId1",
					ID:          "formId1",
					DatabaseID:  "db",
					FolderID:    pointers.String("folder"),
					Name:        "form1",
				}, {
					ID:          "formId2",
					DatabaseID:  "db",
					FolderID:    pointers.String("folder"),
					RootOwnerID: "formId1",
					OwnerID:     "formId1",
					Name:        "form2",
				},
			},
			fields: []*Field{
				{
					ID:         "formId2",
					DatabaseID: "db",
					FormID:     "formId1",
					RootFormID: "formId1",
					SubFormID:  pointers.String("formId2"),
					Name:       "form2",
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
							ID:   "formId2",
							Name: "form2",
							FieldType: types.FieldType{
								SubForm: &types.FieldTypeSubForm{
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
					RootOwnerID: "formId1",
					OwnerID:     "formId1",
					ID:          "formId1",
					DatabaseID:  "db",
					FolderID:    pointers.String("folder"),
					Name:        "form1",
				}, {
					ID:          "formId2",
					DatabaseID:  "db",
					FolderID:    pointers.String("folder"),
					RootOwnerID: "formId1",
					OwnerID:     "formId1",
					Name:        "form2",
				}, {
					ID:          "formId3",
					DatabaseID:  "db",
					FolderID:    pointers.String("folder"),
					RootOwnerID: "formId1",
					OwnerID:     "formId1",
					Name:        "form3",
				},
			},
			fields: []*Field{
				{
					ID:         "formId2",
					DatabaseID: "db",
					FormID:     "formId1",
					RootFormID: "formId1",
					SubFormID:  pointers.String("formId2"),
					Name:       "form2",
					Type:       types.FieldKindSubForm,
				}, {
					ID:         "formId3",
					DatabaseID: "db",
					FormID:     "formId1",
					RootFormID: "formId1",
					SubFormID:  pointers.String("formId3"),
					Name:       "form3",
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
							ID:   "formId2",
							Name: "form2",
							FieldType: types.FieldType{
								SubForm: &types.FieldTypeSubForm{
									Fields: []*types.FieldDefinition{},
								},
							},
						}, {
							ID:   "formId3",
							Name: "form3",
							FieldType: types.FieldType{
								SubForm: &types.FieldTypeSubForm{
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
					ID:          "formId1",
					RootOwnerID: "formId1",
					OwnerID:     "formId1",
					DatabaseID:  "db",
					FolderID:    pointers.String("folder"),
					Name:        "form1",
				}, {
					ID:          "formId2",
					DatabaseID:  "db",
					FolderID:    pointers.String("folder"),
					RootOwnerID: "formId1",
					OwnerID:     "formId1",
					Name:        "form2",
				}, {
					ID:          "formId3",
					DatabaseID:  "db",
					FolderID:    pointers.String("folder"),
					RootOwnerID: "formId1",
					OwnerID:     "formId2",
					Name:        "form3",
				},
			},
			fields: []*Field{
				{
					ID:         "formId2",
					DatabaseID: "db",
					FormID:     "formId1",
					RootFormID: "formId1",
					SubFormID:  pointers.String("formId2"),
					Name:       "form2",
					Type:       types.FieldKindSubForm,
				}, {
					ID:         "formId3",
					DatabaseID: "db",
					FormID:     "formId2",
					RootFormID: "formId1",
					SubFormID:  pointers.String("formId3"),
					Name:       "form3",
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
							ID:   "formId2",
							Name: "form2",
							FieldType: types.FieldType{
								SubForm: &types.FieldTypeSubForm{
									Fields: []*types.FieldDefinition{
										{
											ID:   "formId3",
											Name: "form3",
											FieldType: types.FieldType{
												SubForm: &types.FieldTypeSubForm{
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
