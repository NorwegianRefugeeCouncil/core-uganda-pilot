package sqlconvert

import (
	"database/sql"
	"fmt"
	sqlschema2 "github.com/nrc-no/core/pkg/sqlschema"
	types2 "github.com/nrc-no/core/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_createTable(t *testing.T) {
	tests := []struct {
		name    string
		args    sqlschema2.SQLTable
		wantErr bool
	}{
		{
			name: "simple",
			args: sqlschema2.SQLTable{
				Schema: "db",
				Name:   "example",
				Fields: []sqlschema2.SQLField{
					{
						Name: "field",
						DataType: sqlschema2.SQLDataType{
							Int: &sqlschema2.SQLDataTypeInt{},
						},
					},
				},
			},
			wantErr: false,
		}, {
			name: "primaryKey",
			args: sqlschema2.NewSQLTable("db", "simple").
				WithField(sqlschema2.NewSQLField("id").
					WithSerialDataType().
					WithPrimaryKeyConstraint("pk_simple_id")),
			wantErr: false,
		}, {
			name: "uniqueConstraint",
			args: sqlschema2.NewSQLTable("db", "uniqueConstraint").
				WithField(sqlschema2.NewSQLField("id").
					WithSerialDataType().
					WithPrimaryKeyConstraint("pk_simple_id")).
				WithUniqueConstraint("uq_simple_id", "id"),
			wantErr: false,
		},
	}

	db := dbOpen(t)
	defer db.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := deleteTableIfExists(db, tt.args.Schema, tt.args.Name); err != nil {
				t.Fatal(err)
			}
			if err := createTable(ctx, db, tt.args); (err != nil) != tt.wantErr {
				t.Errorf("createTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createDatabase(t *testing.T) {
	tests := []struct {
		name    string
		args    types2.Database
		wantErr bool
	}{
		{
			name:    "simple",
			args:    types2.Database{Name: "bla"},
			wantErr: false,
		},
	}

	db := dbOpen(t)
	defer db.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !assert.NoError(t, DeleteDatabaseIfExists(db, tt.args.Name)) {
				return
			}
			if err := CreateDatabase(db, tt.args); (err != nil) != tt.wantErr {
				t.Errorf("CreateDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createForm(t *testing.T) {
	tests := []struct {
		name     string
		database types2.Database
		form     types2.FormDefinition
		wantErr  bool
	}{
		{
			name: "simple",
			form: types2.FormDefinition{
				Name:         "form",
				DatabaseName: "db",
				Fields: []types2.FieldDefinition{
					{
						Name: "field",
						FieldType: types2.FieldType{
							Text: &types2.FieldTypeText{},
						},
					},
				},
			},
			database: types2.NewDatabase("db"),
			wantErr:  false,
		}, {
			name: "nested",
			form: types2.FormDefinition{
				Name:         "form",
				DatabaseName: "db",
				Fields: []types2.FieldDefinition{
					{
						Name: "nested",
						FieldType: types2.FieldType{
							SubForm: &types2.FieldTypeSubForm{
								Fields: []types2.FieldDefinition{
									{
										Name: "field",
										FieldType: types2.FieldType{
											Text: &types2.FieldTypeText{},
										},
									},
								},
							},
						},
					},
				},
			},
			database: types2.NewDatabase("db"),
			wantErr:  false,
		},
	}

	db := dbOpen(t)
	defer db.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !assert.NoError(t, deleteFormIfExists(db, tt.form)) {
				return
			}
			if !assert.NoError(t, DeleteDatabaseIfExists(db, tt.database.Name)) {
				return
			}
			if err := CreateDatabase(db, tt.database); !assert.NoError(t, err) {
				return
			}
			if err := CreateForm(db, tt.form); !assert.NoError(t, err) {
				return
			}
		})
	}
}

func dbOpen(t *testing.T) *sql.DB {
	host := "localhost"
	port := 5435
	user := "postgres"
	password := "postgres"
	dbname := "core"

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
	return db
}
