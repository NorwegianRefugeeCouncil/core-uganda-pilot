package sqlmanager

import (
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/api/types"
)

// getSQLActionsForField returns the SQL Actions necessary to store the information for a single field
// getSQLActionsForForm calls this method for every field in a form
func getSQLActionsForField(formInterface types.FormInterface, fieldDefinition *types.FieldDefinition) (sqlActions, error) {
	fieldKind, err := fieldDefinition.FieldType.GetFieldKind()
	if err != nil {
		return nil, err
	}
	switch fieldKind {
	case types.FieldKindSingleSelect:
		return nil, errors.New("not implemented")
	case types.FieldKindSubForm:
		return subFormFieldActions(formInterface, fieldDefinition)
	case types.FieldKindReference:
		return referenceFieldActions(formInterface, fieldDefinition), nil
	case types.FieldKindQuantity:
		return quantityFieldActions(formInterface, fieldDefinition), nil
	case types.FieldKindDate:
		return dateFieldActions(formInterface, fieldDefinition), nil
	case types.FieldKindMonth:
		return monthFieldActions(formInterface, fieldDefinition), nil
	case types.FieldKindMultilineText:
		return multilineTextFieldActions(formInterface, fieldDefinition), nil
	case types.FieldKindText:
		return textFieldActions(formInterface, fieldDefinition), nil
	}
	return nil, fmt.Errorf("unable to convert field kind '%s'", fieldKind)
}
