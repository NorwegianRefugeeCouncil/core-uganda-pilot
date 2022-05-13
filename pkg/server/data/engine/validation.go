package engine

import (
	"github.com/nrc-no/core/pkg/server/data/api"
)

// validateColumnType checks that the data type is a valid column type.
func validateColumnType(typeName string) bool {
	var allowedTypes = map[string]bool{
		"varchar":   true,
		"integer":   true,
		"timestamp": true,
		"boolean":   true,
	}
	return allowedTypes[typeName]
}

// validateIdentifier validates that a given name is a valid identifier
// for use in the database, such as a table or column name.
func validateIdentifier(name string) bool {
	if len(name) == 0 {
		return false
	}
	var prohibitedNames = map[string]bool{
		api.KeyRecordID:  true,
		api.KeyRevision:  true,
		api.KeyPrevision: true,
	}
	if prohibitedNames[name] {
		return false
	}
	for _, c := range name {
		if c != '_' && (c < '0' || c > 'z') {
			return false
		}
	}
	return true
}

// validateTable validates that a given table definition is valid
func validateTable(table api.Table) error {
	if !validateIdentifier(table.Name) {
		return api.ErrInvalidTableName
	}

	if len(table.Columns) == 0 {
		return api.ErrEmptyColumns
	}

	var columnNames = map[string]bool{}
	for _, column := range table.Columns {
		if columnNames[column.Name] {
			return api.NewDuplicateColumnNameErr(column.Name)
		}
		if !validateIdentifier(column.Name) {
			return api.NewInvalidColumnNameErr(column.Name)
		}
		if !validateColumnType(column.Type) {
			return api.NewInvalidColumnTypeErr(column.Type)
		}
		columnNames[column.Name] = true
	}
	return nil
}
