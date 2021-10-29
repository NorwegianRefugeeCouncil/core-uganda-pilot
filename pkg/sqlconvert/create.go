package sqlconvert

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/nrc-no/core/pkg/sqlschema"
	"github.com/nrc-no/core/pkg/types"
	"gorm.io/gorm"
)

func CreateDatabase(db *sql.DB, database *types.Database) error {
	ddl := convertDatabaseToSqlSchema(*database).DDL()
	fmt.Println(ddl)
	_, err := db.Exec(ddl.Query, ddl.Args...)
	return err
}

// CreateForm will create the SQL schemas and tables necessary to store the form records.
// The 4th parameter, referencedForms, is necessary when the form has Reference fields.
// In order to create the appropriate foreign keys, we need to obtain the information on the
// keys of the Referenced forms.
//
// ```
// 		# For example
// 		# Given Form A
// 		#   with text field (id1)
// 		#   with text field (id2)
// 		#   with key (id1, id2)
//
// 		# Given Form B
// 		#    with text field (name)
// 		#    with text field (id)
// 		#    with reference field (snip) that references Form A
// 		#    with key on (snip)
//
// 		# Given Form C
// 		#    with text field (comment)
// 		#    with reference field (snap) that references Form B
// 		#    with reference field (snop) that references Form B
// 		#    with key on (snap)
//
// 		# ================
// 		# Table for Form A
// 		# ================
//
// 		# id1 field
// 		# =========
// 		id1  varchar(36)
//
// 		# id2 field
// 		# =========
// 		id2  varchar(36)
//
// 		# constraints
// 		# ===========
// 		constraint pk_form_a primary key (id1, id2)  # Primary Key on id1, id2
//
// 		# ================
// 		# Table for Form B
// 		# ================
//
// 		# name field
// 		# ==========
// 		name varchar(1024)
//
// 		# id field
// 		# ========
// 		id   varchar(1024)
//
// 		# snip field
// 		# ==========
// 		# Here, 2 columns are required for "snip"
// 		# because it reference FormA and FormA
// 		# has 2 fields in its key (id1, id2)
//
// 		snip_id1 varchar(1024)
// 		snip_id2 varchar(1024)
//
// 		# constraints
// 		# ===========
//
// 		# Primary Key of FormA
// 		constraint pk_form_b
// 		   primary_key (snip_id1, snip_id2)
//
// 		# Foreign Key constraint on FormA (id1, id2)
// 		constraint fk_form_b_form_a
// 		   foreign key (snip_id1, snip_id2)
// 		   references table_a (id1, id2)
//
//
// 		# ================
// 		# Table for Form C
// 		# ================
//
// 		# snap field
// 		# ==========
// 		# Here, 2 columns are required for the "snap" field
// 		# because it references a form that has 2 fields in its key (FormB)
//
// 		snap_id1 varchar(36)
// 		snap_id2 varchar(36)
//
// 		# snop field
// 		# ==========
// 		# Here, 2 columns are also required for the "snop" field,
// 		# because it also references a form that has 2 fields in its key (FormB)
//
// 		snop_id1 varchar(36)
// 		snop_id2 varchar(36)
//
// 		# constraints
// 		# ===========
//
// 		# Primary key on field snap
// 		constraint pk_form_b
// 		   primary key (snap_id1, snap_id2)
//
// 		# Foreign key (snap) -> (FormB.snip)
// 		constraint fk_form_b_snap_form_a
// 		   foreign key (snap_id1, snap_id2)
// 		   references table_b (snip_id1, snip_id2)
//
// 		# Foreign key (snop) -> (FormB.snip)
// 		constraint fk_form_b_snop_form_a
// 		   foreign key (snop_id1, snop_id2)
// 		   references table_b (snip_id1, snip_id2)
// ```
//
func CreateForm(ctx context.Context, db *gorm.DB, form *types.FormDefinition, referencedForms *types.FormDefinitionList) error {
	allForms := expandSubForms(form)
	for _, expanded := range allForms {
		table, err := convertFormToSqlTable(expanded, referencedForms)
		if err != nil {
			return err
		}
		err = createTable(ctx, db, table)
		if err != nil {
			return err
		}
	}
	return nil
}

func createTable(ctx context.Context, db *gorm.DB, table sqlschema.SQLTable) error {

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

	if err := db.WithContext(ctx).Exec(query, ddl.Args...).Error; err != nil {
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

func DeleteDatabaseIfExists(db *gorm.DB, databaseID string) error {
	return deleteSchemaIfExists(db, databaseID)
}

func deleteSchemaIfExists(db *gorm.DB, schemaName string) error {
	err := db.Exec(fmt.Sprintf("drop schema if exists %s cascade", pq.QuoteIdentifier(schemaName))).Error
	return err
}
