package seeder

import "github.com/nrc-no/core/pkg/api/types"

var (
	ifOtherPleaseSpecify = &types.FieldDefinition{
    Name: "If other, please specify",
    FieldType: types.FieldType{Text: &types.FieldTypeText{}},
  }
)

func yesNoQuestion(name string) *types.FieldDefinition {
	return &types.FieldDefinition{
    Name: name,
    FieldType: types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: yesNoChoice}},
  }
}
