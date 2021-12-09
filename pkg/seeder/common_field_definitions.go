package seeder

import "github.com/nrc-no/core/pkg/api/types"

var (
	ifOtherPleaseSpecify = newFieldDefinition("If other, please specify", "", false, false, types.FieldType{Text: &types.FieldTypeText{}})
)

func yesNoQuestion(name string) *types.FieldDefinition {
	return newFieldDefinition(name, "", false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: yesNoChoice}})
}
