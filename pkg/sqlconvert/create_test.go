package sqlconvert

import (
	"database/sql"
	"fmt"
	"github.com/nrc-no/core/pkg/sqlschema"
	"github.com/nrc-no/core/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_createTable(t *testing.T) {
	tests := []struct {
		name    string
		args    sqlschema.SQLTable
		wantErr bool
	}{
		{
			name: "simple",
			args: sqlschema.SQLTable{
				Schema: "db",
				Name:   "example",
				Fields: []sqlschema.SQLField{
					{
						Name: "field",
						DataType: sqlschema.SQLDataType{
							Int: &sqlschema.SQLDataTypeInt{},
						},
					},
				},
			},
			wantErr: false,
		}, {
			name: "primaryKey",
			args: sqlschema.NewSQLTable("db", "simple").
				WithField(sqlschema.NewSQLField("id").
					WithSerialDataType().
					WithPrimaryKeyConstraint("pk_simple_id")),
			wantErr: false,
		}, {
			name: "uniqueConstraint",
			args: sqlschema.NewSQLTable("db", "uniqueConstraint").
				WithField(sqlschema.NewSQLField("id").
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
		args    types.Database
		wantErr bool
	}{
		{
			name:    "simple",
			args:    types.Database{Name: "bla"},
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
		database types.Database
		form     types.FormDefinition
		wantErr  bool
	}{
		{
			name: "simple",
			form: types.FormDefinition{
				Name:         "form",
				DatabaseName: "db",
				Fields: []types.FieldDefinition{
					{
						Name: "field",
						FieldType: types.FieldType{
							Text: &types.FieldTypeText{},
						},
					},
				},
			},
			database: types.NewDatabase("db"),
			wantErr:  false,
		}, {
			name: "nested",
			form: types.FormDefinition{
				Name:         "form",
				DatabaseName: "db",
				Fields: []types.FieldDefinition{
					{
						Name: "nested",
						FieldType: types.FieldType{
							SubForm: &types.FieldTypeSubForm{
								Fields: []types.FieldDefinition{
									{
										Name: "field",
										FieldType: types.FieldType{
											Text: &types.FieldTypeText{},
										},
									},
								},
							},
						},
					},
				},
			},
			database: types.NewDatabase("db"),
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
