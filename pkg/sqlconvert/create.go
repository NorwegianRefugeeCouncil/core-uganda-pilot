package sqlconvert

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/nrc-no/core/pkg/sqlschema"
	"github.com/nrc-no/core/pkg/types"
)

func CreateDatabase(db *sql.DB, database *types.Database) error {
	ddl := convertDatabaseToSqlSchema(*database).DDL()
	fmt.Println(ddl)
	_, err := db.Exec(ddl.Query, ddl.Args...)
	return err
}

func CreateForm(ctx context.Context, db *sql.DB, form *types.FormDefinition) error {
	allForms := expandSubForms(form)
	for _, expanded := range allForms {
		table := convertFormToSqlTable(expanded)
		err := createTable(ctx, db, table)
		if err != nil {
			return err
		}
	}
	return nil
}

func createTable(ctx context.Context, db *sql.DB, table sqlschema.SQLTable) error {

	ddl := table.DDL()

	for _, field := range table.Fields {
		if len(field.Comment) != 0 {
			ddl.WriteF("\ncomment on %s.%s.%s is $1;",
				pq.QuoteIdentifier(table.Schema),
				pq.QuoteIdentifier(table.Name),
				pq.QuoteIdentifier(field.Name))
		}
	}

	fmt.Println(ddl)

	query := ddl.Query
	_, err := db.ExecContext(ctx, query, ddl.Args...)
	if err != nil {
		return err
	}

	return nil

}

func deleteFormIfExists(db *sql.DB, formDef types.FormDefinition) error {
	qry := fmt.Sprintf("drop table if exists %s.%s cascade",
		pq.QuoteIdentifier(formDef.DatabaseID),
		pq.QuoteIdentifier(formDef.ID),
	)
	fmt.Println(qry)
	_, err := db.Exec(qry)
	return err
}

func deleteTableIfExists(db *sql.DB, schemaName, tableName string) error {
	_, err := db.Exec(fmt.Sprintf("drop table if exists %s.%s cascade",
		pq.QuoteIdentifier(schemaName),
		pq.QuoteIdentifier(tableName)))
	return err
}

func DeleteDatabaseIfExists(db *sql.DB, databaseID string) error {
	return deleteSchemaIfExists(db, databaseID)
}

func deleteSchemaIfExists(db *sql.DB, schemaName string) error {
	_, err := db.Exec(fmt.Sprintf("drop schema if exists %s cascade", pq.QuoteIdentifier(schemaName)))
	return err
}
