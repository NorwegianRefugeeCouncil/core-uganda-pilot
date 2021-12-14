package sqlmanager

import (
	"fmt"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/mocks"
	"github.com/nrc-no/core/pkg/sql/schema"
	"github.com/nrc-no/core/pkg/testutils"
	"github.com/nrc-no/core/pkg/utils/pointers"
	uuid "github.com/satori/go.uuid"
	"github.com/snabb/isoweek"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strings"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	db   *gorm.DB
	done func()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func TestMain(m *testing.M) {
	// embedded-postgres runs as a separate process altogether
	// We must setup/teardown here, otherwise the process does
	// not get properly cleaned up
	done, err := testutils.TryGetPostgres()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer func() {
		recover()
		done()
	}()
	exitVal := m.Run()
	done()
	os.Exit(exitVal)
}

func (s *Suite) SetupSuite() {

	sqlDb, err := gorm.Open(postgres.Open("host=localhost port=15432 user=postgres password=postgres dbname=postgres sslmode=disable"))
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	dbFactory := mocks.NewMockFactory(sqlDb)

	db, err := dbFactory.Get()
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}

	s.db = db

}

func (s *Suite) TearDownSuite() {
	if s.done != nil {
		s.done()
	}
}

func TestSchema(t *testing.T) {
	s := writer{}
	var err error

	s, err = s.handleCreateTable(sqlActionCreateTable{
		sqlTable: schema.SQLTable{
			Schema:      "schema",
			Name:        "tableName",
			Columns:     schema.SQLColumns{},
			Constraints: []schema.SQLTableConstraint{},
		},
	})

	s, err = s.handleCreateColumn(sqlActionCreateColumn{
		tableName:  "tableName",
		schemaName: "schemaName",
		sqlColumn: schema.SQLColumn{
			Name: "field",
			DataType: schema.SQLDataType{
				VarChar: &schema.SQLDataTypeVarChar{
					Length: 10,
				},
			},
		},
	})

	s, err = s.handleCreateConstraint(sqlActionCreateConstraint{
		tableName:  "tableName",
		schemaName: "schemaName",
		sqlConstraint: schema.SQLTableConstraint{
			Name: "uk_constraint",
			Unique: &schema.SQLTableConstraintUnique{
				ColumnNames: []string{"column1"},
			},
		},
	})

	assert.Equal(t, []schema.DDL{
		schema.NewDDL(`create table "schema"."tableName"();`),
		schema.NewDDL(`alter table "schemaName"."tableName" add "field" varchar(10);`),
		schema.NewDDL(`alter table "schemaName"."tableName" add constraint "uk_constraint" unique ("column1");`),
	}, s.Statements)

	assert.NoError(t, err)
}

func TestWriterFormConversion(t *testing.T) {

	const createTableDDL = `create table "databaseId"."formId"( "id" varchar(36) primary key, "created_at" timestamp with time zone not null default NOW());`
	const formId = "formId"
	const databaseId = "databaseId"

	tests := []struct {
		name    string
		args    []*types.FormDefinition
		want    []schema.DDL
		wantErr bool
	}{
		{
			name: "empty form",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
			},
		},
		{
			name: "form with text field",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "textField",
						FieldType: types.FieldType{
							Text: &types.FieldTypeText{},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `alter table "databaseId"."formId" add "textField" varchar(1024);`},
			},
		},
		{
			name: "form with multiLine text field",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "multiLineTextField",
						FieldType: types.FieldType{
							MultilineText: &types.FieldTypeMultilineText{},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `alter table "databaseId"."formId" add "multiLineTextField" text;`},
			},
		},
		{
			name: "form with date field",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "dateField",
						FieldType: types.FieldType{
							Date: &types.FieldTypeDate{},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `alter table "databaseId"."formId" add "dateField" date;`},
			},
		},
		{
			name: "form with week field",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "weekField",
						FieldType: types.FieldType{
							Week: &types.FieldTypeWeek{},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `alter table "databaseId"."formId" add "weekField" date;`},
			},
		},
		{
			name: "form with month field",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "monthField",
						FieldType: types.FieldType{
							Month: &types.FieldTypeMonth{},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `alter table "databaseId"."formId" add "monthField" date;`},
			},
		},
		{
			name: "form with quantity field",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "quantityField",
						FieldType: types.FieldType{
							Quantity: &types.FieldTypeQuantity{},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `alter table "databaseId"."formId" add "quantityField" int;`},
			},
		},
		{
			name: "form with subForm and key fields",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "subFormField",
						FieldType: types.FieldType{
							SubForm: &types.FieldTypeSubForm{
								Fields: []*types.FieldDefinition{
									{
										ID:  "subTextField",
										Key: true,
										FieldType: types.FieldType{
											Text: &types.FieldTypeText{},
										},
									},
								},
							},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `
create table "databaseId"."subFormField"(
  "id" varchar(36) primary key,
  "created_at" timestamp with time zone not null default NOW(),
  "owner_id" varchar(36) not null references "databaseId"."formId" ("id")
);`},
				{Query: `alter table "databaseId"."subFormField" add "subTextField" varchar(1024) not null;`},
				{Query: `alter table "databaseId"."subFormField" add constraint "uk_key_subFormField" unique ("subTextField");`},
			},
		},
		{
			name: "form with reference field",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "referenceField",
						FieldType: types.FieldType{
							Reference: &types.FieldTypeReference{
								DatabaseID: "otherDatabaseId",
								FormID:     "otherFormId",
							},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `alter table "databaseId"."formId" add "referenceField" varchar(36) references "otherDatabaseId"."otherFormId" ("id");`},
			},
		},
		{
			name: "form with single select field",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "singleSelectField",
						FieldType: types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{
								Options: []*types.SelectOption{
									{
										ID:   "option1",
										Name: "Option 1",
									}, {
										ID:   "option2",
										Name: "Option 2",
									},
								},
							},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `create table "databaseId"."singleSelectField_options"( "id" varchar(36) primary key, "name" varchar(128) not null unique);`},
				{Query: `insert into "databaseId"."singleSelectField_options" ("id","name") values ($1,$2);`, Args: []interface{}{"option1", "Option 1"}},
				{Query: `insert into "databaseId"."singleSelectField_options" ("id","name") values ($1,$2);`, Args: []interface{}{"option2", "Option 2"}},
				{Query: `alter table "databaseId"."formId" add "singleSelectField" varchar(36) references "databaseId"."singleSelectField_options" ("id");`},
			},
		},
		{
			name: "form with multi select field",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "multiSelectField",
						FieldType: types.FieldType{
							MultiSelect: &types.FieldTypeMultiSelect{
								Options: []*types.SelectOption{
									{
										ID:   "option1",
										Name: "Option 1",
									}, {
										ID:   "option2",
										Name: "Option 2",
									},
								},
							},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `create table "databaseId"."multiSelectField_options"( "id" varchar(36) primary key, "name" varchar(128) not null unique);`},
				{Query: `insert into "databaseId"."multiSelectField_options" ("id","name") values ($1,$2);`, Args: []interface{}{"option1", "Option 1"}},
				{Query: `insert into "databaseId"."multiSelectField_options" ("id","name") values ($1,$2);`, Args: []interface{}{"option2", "Option 2"}},
				{Query: `
create table "databaseId"."multiSelectField_associations"(
  "id" varchar(36) not null references "databaseId"."formId" ("id"),
  "option_id" varchar(36) not null references "databaseId"."multiSelectField_options" ("id"),
  constraint "uk_key_multiSelectField" unique ("id", "option_id")
);`},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			manager := New()
			manager, err := manager.PutForms(&types.FormDefinitionList{Items: test.args})
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			if !assert.NoError(t, err) {
				return
			}
			statements := manager.GetStatements()
			for i := 0; i < len(statements); i++ {
				ddl := statements[i]
				ddl.Query = strings.ReplaceAll(ddl.Query, "\n", "")
				ddl.Query = strings.ReplaceAll(ddl.Query, "  ", " ")
				statements[i] = ddl
			}
			for i := 0; i < len(test.want); i++ {
				ddl := test.want[i]
				ddl.Query = strings.ReplaceAll(ddl.Query, "\n", "")
				ddl.Query = strings.ReplaceAll(ddl.Query, "  ", " ")
				test.want[i] = ddl
			}
			assert.Equal(t, test.want, statements)
		})
	}
}

func (s *Suite) TestWriterActions() {

	const publicSchema = "public"
	const formId = "formId"
	const otherFormId = "other-form"
	const field1Id = "field-1"

	otherForm := &types.FormDefinition{
		ID:         otherFormId,
		DatabaseID: publicSchema,
		Fields: []*types.FieldDefinition{
			{
				ID: field1Id,
				FieldType: types.FieldType{
					Text: &types.FieldTypeText{},
				},
			},
		},
	}

	formDef := &types.FormDefinition{
		ID:         formId,
		DatabaseID: publicSchema,
		Fields: []*types.FieldDefinition{
			{
				ID:  "textField",
				Key: true,
				FieldType: types.FieldType{
					Text: &types.FieldTypeText{},
				},
			},
			{
				ID:  "dateField",
				Key: true,
				FieldType: types.FieldType{
					Date: &types.FieldTypeDate{},
				},
			},
			{
				ID: "referenceField",
				FieldType: types.FieldType{
					Reference: &types.FieldTypeReference{
						DatabaseID: publicSchema,
						FormID:     otherFormId,
					},
				},
			},
			{
				ID: "subFormField",
				FieldType: types.FieldType{
					SubForm: &types.FieldTypeSubForm{
						Fields: []*types.FieldDefinition{
							{
								ID: "sub-field",
								FieldType: types.FieldType{
									Text: &types.FieldTypeText{},
								},
							},
						},
					},
				},
			},
		},
	}

	manager := New()
	var err error
	manager, err = manager.PutForms(&types.FormDefinitionList{
		Items: []*types.FormDefinition{
			otherForm,
			formDef,
		},
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	statements := manager.GetStatements()

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		db, err := tx.DB()
		if err != nil {
			return err
		}
		var ddl string
		var args []interface{}
		for _, statement := range statements {
			ddl = ddl + statement.Query + "\n"
			args = append(args, statement.Args...)
		}

		s.T().Logf("executing statement: \n%s", ddl)

		if _, err := db.Exec(ddl, args...); err != nil {
			s.T().Logf("error with statement %s", ddl)
			return err
		}

		return nil

	}); !assert.NoError(s.T(), err) {
		return
	}

}

func TestWriterPutRecords(t *testing.T) {

	const (
		formId      = "formId"
		databaseId  = "databaseId"
		ownerFormId = "ownerFormId"
	)

	formOpts := []testutils.FormOption{
		testutils.FormID(formId),
		testutils.FormDatabaseID(databaseId),
	}

	formWithFields := func(options ...testutils.FormOption) types.FormInterface {
		return testutils.AForm(append(formOpts, options...)...)
	}

	subFormWithFields := func(options ...testutils.FormOption) types.FormInterface {
		return testutils.ASubForm(ownerFormId, append(formOpts, options...)...)
	}

	aRecord := func(options ...testutils.RecordOption) *types.Record {
		opts := []testutils.RecordOption{
			testutils.RecordID("recordId"),
			testutils.RecordFormID(formId),
			testutils.RecordDatabaseID(databaseId),
		}
		opts = append(opts, options...)
		return testutils.RecordOptions(opts...)(&types.Record{})
	}

	aSingleRecord := func(options ...testutils.RecordOption) *types.RecordList {
		return &types.RecordList{
			Items: []*types.Record{
				aRecord(options...),
			},
		}
	}

	mustParseTime := func(layout string, str string) time.Time {
		tm, err := time.Parse(layout, str)
		if err != nil {
			panic(err)
		}
		return tm
	}

	aFormReference := types.FormRef{
		DatabaseID: uuid.NewV4().String(),
		FormID:     uuid.NewV4().String(),
	}

	tests := []struct {
		name    string
		form    types.FormInterface
		records *types.RecordList
		want    []schema.DDL
		wantErr bool
	}{
		{
			name:    "record with text value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeText()))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewStringValue("myValue"))),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" ("id","fieldId") values ($1,$2);`,
					Args:  []interface{}{"recordId", "myValue"},
				},
			},
		}, {
			name:    "record with null text value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeText()))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewNullValue())),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" ("id","fieldId") values ($1,$2);`,
					Args:  []interface{}{"recordId", nil},
				},
			},
		}, {
			name:    "record with multiline text value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeMultilineText()))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewStringValue("myValue"))),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" ("id","fieldId") values ($1,$2);`,
					Args:  []interface{}{"recordId", "myValue"},
				},
			},
		}, {
			name:    "record with nil multiline text value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeMultilineText()))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewNullValue())),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" ("id","fieldId") values ($1,$2);`,
					Args:  []interface{}{"recordId", nil},
				},
			},
		}, {
			name:    "record with month value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeMonth()))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewStringValue("2020-01"))),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" ("id","fieldId") values ($1,$2);`,
					Args:  []interface{}{"recordId", mustParseTime(monthFieldFormat, "2020-01")},
				},
			},
		}, {
			name:    "record with nil month value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeMonth()))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewNullValue())),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" ("id","fieldId") values ($1,$2);`,
					Args:  []interface{}{"recordId", nil},
				},
			},
		}, {
			name:    "record with bad month value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeMonth()))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewStringValue("abc"))),
			wantErr: true,
		}, {
			name:    "record with week value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeWeek()))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewStringValue("2020-W01"))),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" ("id","fieldId") values ($1,$2);`,
					Args:  []interface{}{"recordId", isoweek.StartTime(2020, 1, time.UTC)},
				},
			},
		}, {
			name:    "record with nil week value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeWeek()))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewNullValue())),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" ("id","fieldId") values ($1,$2);`,
					Args:  []interface{}{"recordId", nil},
				},
			},
		}, {
			name:    "record with bad week value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeWeek()))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewStringValue("2020-W75"))),
			wantErr: true,
		}, {
			name:    "record with date value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeDate()))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewStringValue("2020-01-01"))),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" ("id","fieldId") values ($1,$2);`,
					Args:  []interface{}{"recordId", mustParseTime(dateFieldFormat, "2020-01-01")},
				},
			},
		}, {
			name:    "record with nil date value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeDate()))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewNullValue())),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" ("id","fieldId") values ($1,$2);`,
					Args:  []interface{}{"recordId", nil},
				},
			},
		}, {
			name:    "record with bad date value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeDate()))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewStringValue("abc"))),
			wantErr: true,
		}, {
			name:    "record with reference value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeReference(aFormReference)))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewStringValue("refId"))),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" ("id","fieldId") values ($1,$2);`,
					Args:  []interface{}{"recordId", "refId"},
				},
			},
		}, {
			name:    "record with nil reference value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeReference(aFormReference)))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewNullValue())),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" ("id","fieldId") values ($1,$2);`,
					Args:  []interface{}{"recordId", nil},
				},
			},
		}, {
			name:    "record with quantity value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeQuantity()))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewStringValue("3"))),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" ("id","fieldId") values ($1,$2);`,
					Args:  []interface{}{"recordId", 3},
				},
			},
		}, {
			name:    "record with nil quantity value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeQuantity()))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewNullValue())),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" ("id","fieldId") values ($1,$2);`,
					Args:  []interface{}{"recordId", nil},
				},
			},
		}, {
			name:    "record with bad quantity value",
			form:    formWithFields(testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeQuantity()))),
			records: aSingleRecord(testutils.RecordValue("fieldId", types.NewStringValue("abc"))),
			wantErr: true,
		}, {
			name: "record with owner",
			form: subFormWithFields(
				testutils.FormField(testutils.AField(testutils.FieldID("fieldId"), testutils.FieldTypeQuantity())),
			),
			records: aSingleRecord(
				testutils.RecordOwnerID(pointers.String("ownerId")),
				testutils.RecordValue("fieldId", types.NewStringValue("123")),
			),
			want: []schema.DDL{
				{
					Query: `insert into "databaseId"."formId" ("id","owner_id","fieldId") values ($1,$2,$3);`,
					Args:  []interface{}{"recordId", "ownerId", 123},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f Writer = &writer{}
			var err error
			f, err = f.PutRecords(tt.form, tt.records)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if !assert.NoError(t, err) {
				return
			}
			got := f.GetStatements()
			t.Log(got)
			assert.Equal(t, tt.want, got)
		})
	}
}
