package seeder

import "github.com/nrc-no/core/pkg/api/types"

// NOTE: the variables and functions defined here are very simple, most of the fields are
// default (ie: nil or false). Don't use them if you need to specify other non-default field values.

var (
	ifOtherPleaseSpecify = &types.FieldDefinition{
		Name:      "If other, please specify",
		FieldType: types.FieldType{Text: &types.FieldTypeText{}},
	}
)

func yesNo(name string) *types.FieldDefinition {
	return &types.FieldDefinition{
		Name:      name,
		Required:  true,
		FieldType: types.FieldType{Boolean: &types.FieldTypeBoolean{}},
	}
}

func text(name string, required bool) *types.FieldDefinition {
	return &types.FieldDefinition{
		Name:      name,
		Required:  required,
		FieldType: types.FieldType{Text: &types.FieldTypeText{}},
	}
}

func textarea(name string, required bool) *types.FieldDefinition {
	return &types.FieldDefinition{
		Name:      name,
		Required:  required,
		FieldType: types.FieldType{MultilineText: &types.FieldTypeMultilineText{}},
	}
}

func date(name string, required bool) *types.FieldDefinition {
	return &types.FieldDefinition{
		Name:      name,
		Required:  required,
		FieldType: types.FieldType{Date: &types.FieldTypeDate{}},
	}
}

func dropdown(name string, options []*types.SelectOption, required bool) *types.FieldDefinition {
	return &types.FieldDefinition{
		Name:      name,
		Required:  required,
		FieldType: types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: options}},
	}
}

func multiSelect(name string, options []*types.SelectOption, required bool) *types.FieldDefinition {
	return &types.FieldDefinition{
		Name:      name,
		Required:  required,
		FieldType: types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: options}},
	}
}

func quantity(name string, required bool) *types.FieldDefinition {
	return &types.FieldDefinition{
		Name:      name,
		Required:  required,
		FieldType: types.FieldType{Quantity: &types.FieldTypeQuantity{}},
	}
}
