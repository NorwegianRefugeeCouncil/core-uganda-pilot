package sqlmanager

import (
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/schema"
)

// getSQLActionsForForm returns the sqlActions that are necessary to create the SQL tables, columns, etc.
// necessary to store the form records.
func getSQLActionsForForm(formInterface types.FormInterface) (sqlActions, error) {
	result := sqlActions{}

	// all SQL Tables have an "id" column
	columns := []schema.SQLColumn{
		buildTableIDColumn(),
		buildTableCreatedAtColumn(),
	}

	// If the form is a subform, then we also add the "owner_id" column
	if formInterface.HasOwner() {
		columns = append(columns, buildSubFormOwnerColumn(formInterface))
	}

	// Append the sqlActionCreateTable to the list of actions
	result = append(result, sqlAction{
		createTable: &sqlActionCreateTable{
			sqlTable: schema.SQLTable{
				Schema:      formInterface.GetDatabaseID(),
				Name:        formInterface.GetFormID(),
				Columns:     columns,
				Constraints: []schema.SQLTableConstraint{},
			},
		},
	})

	// Append the SQL Actions for each field of the FormInterface
	for _, fieldDefinition := range formInterface.GetFields() {
		actionsForField, err := getSQLActionsForField(formInterface, fieldDefinition)
		if err != nil {
			return nil, err
		}
		result = append(result, actionsForField...)
	}

	// Append the SQL Actions for configuring the form's key columns
	result = append(result, buildKeyFieldActions(formInterface)...)

	return result, nil
}
