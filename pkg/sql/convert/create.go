package convert

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/nrc-no/core/pkg/api/types"
	sqlschema2 "github.com/nrc-no/core/pkg/sql/schema"
	"gorm.io/gorm"
)

func CreateDatabase(db *gorm.DB, database *types.Database) error {
	ddl := convertDatabaseToSqlSchema(*database).DDL()
	fmt.Println(ddl)
	err := db.Exec(ddl.Query, ddl.Args...).Error
	return err
}

func DeleteTableIfExists(db *gorm.DB, schemaName, tableName string) error {
	err := db.Exec(fmt.Sprintf("drop table if exists %s.%s cascade",
		pq.QuoteIdentifier(schemaName),
		pq.QuoteIdentifier(tableName))).Error
	return err
}

func DeleteDatabaseSchemaIfExist(db *gorm.DB, databaseID string) error {
	return deleteSchemaIfExists(db, databaseID)
}

func deleteSchemaIfExists(db *gorm.DB, schemaName string) error {
	err := db.Exec(fmt.Sprintf("drop schema if exists %s cascade", pq.QuoteIdentifier(schemaName))).Error
	return err
}

func convertDatabaseToSqlSchema(database types.Database) sqlschema2.SQLSchema {
	return sqlschema2.SQLSchema{
		Name: database.ID,
	}
}
