package sqlmanager

import "github.com/nrc-no/core/pkg/sql/schema"

// sqlAction represents an SQL Operation that we want to execute.
// Only one of createColumn or createTable can be specified at any time
type sqlAction struct {
	// createColumn represents an action that creates an SQL sqlColumn
	createColumn *sqlActionCreateColumn
	// createTable represents an action that creates an SQL table
	createTable *sqlActionCreateTable
	// createUniqueConstraint represents an action that creates an SQL Index
	createUniqueConstraint *sqlActionCreateConstraint
	// insertRow represents an action that adds a row in an SQL table
	insertRow *sqlActionInsertRow
}

// sqlActions is a list of sqlAction
type sqlActions []sqlAction

// sqlActionCreateColumn contain the parameters for creating a schema.SQLColumn
type sqlActionCreateColumn struct {
	// tableName represents the name of the SQL table in which the column will be added
	tableName string
	// schemaName represents the name of the SQL schemaName in which the column will be added
	schemaName string
	// Field represents the SQL Definition of the column to add
	sqlColumn schema.SQLColumn
}

// sqlActionCreateTable contains the parameters for creating a schema.SQLTable
type sqlActionCreateTable struct {
	// sqlTable represents the configuration of the schema.SQLTable
	sqlTable schema.SQLTable
}

// sqlActionCreateConstraint contains the parameters for creating a schema.SQLTable
type sqlActionCreateConstraint struct {
	// tableName represents the sql constraint table name
	tableName string
	// schemaName represents the sql constraint schema name
	schemaName string
	// sqlConstraint represents the configuration of the schema.SQLTableConstraint
	sqlConstraint schema.SQLTableConstraint
}

// sqlActionInsertRow contains the parameters for inserting a new row
type sqlActionInsertRow struct {
	schemaName string
	tableName  string
	columns    []string
	values     []interface{}
}
