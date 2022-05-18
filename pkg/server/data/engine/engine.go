package engine

import (
	"context"
	"fmt"
	"strings"

	"github.com/nrc-no/core/pkg/server/data/api"
)

type engine struct {
	// txFactory is the txFactory factory
	txFactory api.TxFactory
	// uuidGenerator generates uuids
	uuidGenerator api.UUIDGenerator
	// revisionGenerator generates revision hashes
	revisionGenerator api.RevisionGenerator
	// clock is used to get the current time
	clock api.Clock
	// dialect is the dialect used by the engine
	// available dialects are:
	// - "sqlite"
	dialect string
}

func NewEngine(
	ctx context.Context,
	txFactory api.TxFactory,
	uuidGenerator api.UUIDGenerator,
	revisionGenerator api.RevisionGenerator,
	clock api.Clock,
	dialect string,
) (api.Engine, error) {
	e := &engine{
		txFactory:         txFactory,
		uuidGenerator:     uuidGenerator,
		revisionGenerator: revisionGenerator,
		clock:             clock,
		dialect:           dialect,
	}
	if err := e.Init(ctx); err != nil {
		return nil, err
	}
	if dialect != "sqlite" {
		return nil, api.ErrUnsupportedDialect
	}
	return e, nil
}

// Init initializes the engine
// It creates supporting tables if they don't exist
func (e *engine) Init(ctx context.Context) error {
	_, err := e.doTransaction(ctx, func(t api.Transaction) (interface{}, error) {
		if err := e.initChangesTable(ctx, t); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

// PutRecord implements Engine.PutRecord
func (e *engine) PutRecord(ctx context.Context, request api.PutRecordRequest) (api.Record, error) {
	ret, err := e.doTransaction(ctx, func(tx api.Transaction) (interface{}, error) {
		var recPtr = &request.Record
		if err := e.putRecord(ctx, tx, recPtr, request.IsReplication); err != nil {
			return nil, err
		}
		return *recPtr, nil
	})
	if err != nil {
		return api.Record{}, err
	}
	return ret.(api.Record), nil
}

// GetRecord implements Engine.GetRecord
func (e *engine) GetRecord(ctx context.Context, request api.GetRecordRequest) (api.Record, error) {
	ret, err := e.doTransaction(ctx, func(tx api.Transaction) (interface{}, error) {
		if request.Revision.IsEmpty() {
			return e.getLatestRevision(ctx, tx, request)
		} else {
			return e.getRevision(ctx, tx, request)
		}
	})
	if err != nil {
		return api.Record{}, err
	}
	return ret.(api.Record), nil
}

// GetRecords implements Engine.GetRecords
func (e *engine) GetRecords(ctx context.Context, request api.GetRecordsRequest) (api.RecordList, error) {
	ret, err := e.doTransaction(ctx, func(t api.Transaction) (interface{}, error) {
		return e.getRecords(ctx, t, request)
	})
	if err != nil {
		return api.RecordList{}, err
	}
	return ret.(api.RecordList), nil
}

// CreateTable implements Engine.CreateTable
func (e *engine) CreateTable(ctx context.Context, table api.Table) (api.Table, error) {
	_, err := e.doTransaction(ctx, func(tx api.Transaction) (interface{}, error) {
		err := e.createTable(ctx, tx, table)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return table, err
}

// GetChanges implements Engine.GetChanges
func (e *engine) GetChanges(ctx context.Context, request api.GetChangesRequest) (api.Changes, error) {
	ret, err := e.doTransaction(ctx, func(tx api.Transaction) (interface{}, error) {
		return e.getChanges(tx, ctx, request.Since)
	})
	if err != nil {
		return api.Changes{}, err
	}
	return ret.(api.Changes), nil
}

// GetTables implements Engine.GetTables
func (e *engine) GetTables(ctx context.Context, request api.GetTablesRequest) (api.GetTablesResponse, error) {
	ret, err := e.doTransaction(ctx, func(tx api.Transaction) (interface{}, error) {
		return e.getTables(ctx, tx)
	})
	if err != nil {
		return api.GetTablesResponse{}, err
	}
	return ret.(api.GetTablesResponse), nil
}

// GetTable implements Engine.GetTable
func (e *engine) GetTable(ctx context.Context, request api.GetTableRequest) (api.Table, error) {
	ret, err := e.doTransaction(ctx, func(tx api.Transaction) (interface{}, error) {
		return e.getTable(ctx, tx, request.TableName)
	})
	if err != nil {
		return api.Table{}, err
	}
	return ret.(api.Table), nil
}

func (e *engine) getRevision(ctx context.Context, tx api.Transaction, request api.GetRecordRequest) (interface{}, error) {
	found, err := e.findRevision(ctx, tx, request.TableName, request.RecordID, request.Revision)
	if err != nil {
		return api.Record{}, err
	}
	if found == nil {
		return api.Record{}, api.ErrRecordNotFound
	}
	return *found, nil
}

func (e *engine) getLatestRevision(ctx context.Context, tx api.Transaction, request api.GetRecordRequest) (interface{}, error) {
	found, err := e.findLatestRevision(ctx, tx, request.TableName, request.RecordID)
	if err != nil {
		return api.Record{}, err
	}
	if found == nil {
		return api.Record{}, api.ErrRecordNotFound
	}
	return *found, nil
}

func (e *engine) getRecords(ctx context.Context, tx api.Transaction, request api.GetRecordsRequest) (api.RecordList, error) {

	// retrieve the information about the table
	columnNames, columnKinds, err := e.getValueTypeInfoForTable(ctx, tx, request.TableName)
	if err != nil {
		return api.RecordList{}, err
	}

	// build the query
	sqlQuery := &StringBuilder{}
	sqlQuery.WriteString("SELECT " + strings.Join(columnNames, ",") + " FROM " + request.TableName)

	var whereClauses []string
	var params []interface{}
	if len(request.RecordIDs) > 0 {
		// if the request contains explicit record IDs, we need to filter by those
		whereClauses = append(whereClauses, api.KeyRecordID+" IN ("+joinStrings(repeatStrings("?", len(request.RecordIDs)), ",")+")")
		for _, recordId := range request.RecordIDs {
			params = append(params, recordId)
		}
	}

	var groupByClauses []string
	var orderByClauses []string
	if request.Revisions == false {
		// if the request does not request revisions, we need
		// to find the latest revision for each record
		subQuery := &StringBuilder{}
		subQuery.WriteString(`
SELECT MAX(` + api.KeyRevision + `) AS ` + api.KeyRevision + `
FROM ` + request.TableName + `
GROUP BY ` + api.KeyRecordID + `
`)
		whereClauses = append(whereClauses, api.KeyRevision+" IN ("+subQuery.String()+")")
	}

	if len(whereClauses) > 0 {
		sqlQuery.WriteString(" WHERE " + strings.Join(whereClauses, " AND "))
	}
	if len(groupByClauses) > 0 {
		sqlQuery.WriteString(" GROUP BY " + strings.Join(groupByClauses, ","))
	}
	if len(orderByClauses) > 0 {
		sqlQuery.WriteString(" ORDER BY " + strings.Join(orderByClauses, ","))
	}

	sqlStatement := sqlQuery.String()
	result, err := tx.Query(ctx, sqlStatement, params)
	if err != nil {
		return api.RecordList{}, err
	}
	defer closeRows(result)

	var records = make([]api.Record, 0)
	for result.Next() {
		values, err := result.Read(columnKinds)
		if err != nil {
			return api.RecordList{}, err
		}
		record, err := readInRecord(request.TableName, values)
		if err != nil {
			return api.RecordList{}, err
		}
		records = append(records, record)
	}
	if err := result.Err(); err != nil {
		return api.RecordList{}, err
	}

	return api.RecordList{Items: records}, nil

}

func (e *engine) getTable(ctx context.Context, tx api.Transaction, tableName string) (api.Table, error) {
	columnNames, columnKinds, err := e.getColumnTypeInfoForTable(ctx, tx, tableName)
	if err != nil {
		return api.Table{}, err
	}
	var ret = api.Table{}
	ret.Name = tableName
	for i := range columnNames {
		columnName := columnNames[i]
		columnKind := columnKinds[i]
		var column = api.Column{
			Name: columnName,
			Type: columnKind,
		}
		ret.Columns = append(ret.Columns, column)
	}
	return ret, nil
}

func (e *engine) getChanges(tx api.Transaction, ctx context.Context, checkpoint int64) (api.Changes, error) {

	// retrieve the information about the table
	columnNames, columnKinds, err := e.getValueTypeInfoForTable(ctx, tx, api.ChangeStreamTableName)
	if err != nil {
		return api.Changes{}, err
	}

	// build the query
	sqlQuery := `
SELECT ` + joinStrings(columnNames, ",") + ` FROM "` + api.ChangeStreamTableName + `"
WHERE "` + api.KeyCSSequence + `" > ?
ORDER BY "` + api.KeyCSSequence + `" ASC;
`
	// execute the query
	rows, err := tx.Query(ctx, sqlQuery, []interface{}{checkpoint})
	if err != nil {
		return api.Changes{}, err
	}
	defer closeRows(rows)

	// prepare the result
	var records = make([]api.ChangeItem, 0)

	// iterate over the rows
	for rows.Next() {

		var values map[string]api.Value
		var rec api.Record

		// read the values
		if values, err = rows.Read(columnKinds); err != nil {
			return api.Changes{}, err
		}
		// create the record
		if rec, err = readInRecord(api.ChangeStreamTableName, values); err != nil {
			return api.Changes{}, err
		}

		// parse the record
		changeItem, err := parseChangeStreamItem(rec)
		if err != nil {
			return api.Changes{}, err
		}

		// add the record to the result
		records = append(records, changeItem)
	}

	// check if there was an error while iterating
	if err := rows.Err(); err != nil {
		return api.Changes{}, err
	}

	return api.Changes{
		Items: records,
	}, nil
}

// getColumnKinds returns the column types for the given table
func (e *engine) getColumnKinds(ctx context.Context, tx api.Transaction, table string) (map[string]api.ValueKind, error) {
	columnTypes, err := e.getColumnTypes(ctx, tx, table)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]api.ValueKind)
	for columnName, columnType := range columnTypes {
		// map to the correct kind
		switch columnType {
		case "integer":
			ret[columnName] = api.ValueKindInt
		case "real":
			ret[columnName] = api.ValueKindFloat
		case "text", "varchar":
			ret[columnName] = api.ValueKindString
		case "bool":
			ret[columnName] = api.ValueKindBool
		default:
			return nil, fmt.Errorf("unknown column type: %s", columnType)
		}
	}
	return ret, nil
}

// getColumnTypes returns the column types for the given table
func (e *engine) getColumnTypes(ctx context.Context, tx api.Transaction, table string) (map[string]string, error) {

	// these will be the returned column types
	columnKinds := []api.ValueKind{
		api.ValueKindString,
		api.ValueKindString,
	}

	// build the query
	sql := `SELECT name, type FROM PRAGMA_TABLE_INFO(?);`

	// execute the query
	rows, err := tx.Query(ctx, sql, []interface{}{table})
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)

	// prepare the result
	columnTypes := make(map[string]string)
	for rows.Next() {

		// read the values
		row, err := rows.Read(columnKinds)
		if err != nil {
			return nil, err
		}

		// get the name of the column
		columnName := row["name"].Str.ValueOrZero()
		// get the type of the column
		columnType := strings.ToLower(row["type"].Str.ValueOrZero())

		columnTypes[columnName] = columnType

	}
	// check if there was an error while iterating
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// if there are no columns, that means that the table does not exist
	if len(columnTypes) == 0 {
		return nil, fmt.Errorf("table %s does not exist", table)
	}

	return columnTypes, nil
}

// findLatestRevision finds the latest revision of the record
func (e *engine) findLatestRevision(ctx context.Context, tx api.Transaction, table string, id string) (*api.Record, error) {

	// retrieve the column information
	columnNames, columnKinds, err := e.getValueTypeInfoForTable(ctx, tx, table)
	if err != nil {
		return nil, err
	}

	// build the query
	sqlQuery := `SELECT ` + joinStrings(columnNames, ",") + ` FROM "` + table + `" WHERE "` + api.KeyRecordID + `" = ? ORDER BY "` + api.KeyRevision + `" DESC LIMIT 1;`

	// execute the query
	rows, err := tx.Query(ctx, sqlQuery, []interface{}{id})
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)

	// if there are no rows, that means that the record does not exist
	// but this method will return nil instead of an error
	if !rows.Next() {
		return nil, nil
	}

	// read the values
	data, err := rows.Read(columnKinds)
	if err != nil {
		return nil, err
	}

	// create the record
	var record api.Record
	if record, err = readInRecord(table, data); err != nil {
		return nil, err
	}

	return &record, nil
}

// getColumnTypeInfoForTable return the column names and types for the given table
func (e *engine) getColumnTypeInfoForTable(ctx context.Context, tx api.Transaction, table string) ([]string, []string, error) {
	columnTypes, err := e.getColumnTypes(ctx, tx, table)
	if err != nil {
		return nil, nil, err
	}
	var columnNames []string
	var columnKinds []string
	for columnName := range columnTypes {
		columnNames = append(columnNames, columnName)
	}
	sortStrings(columnNames)
	for _, columnName := range columnNames {
		columnKinds = append(columnKinds, columnTypes[columnName])
	}
	return columnNames, columnKinds, nil
}

// getValueTypeInfoForTable returns the column names and value kinds for the given table
func (e *engine) getValueTypeInfoForTable(ctx context.Context, tx api.Transaction, table string) ([]string, []api.ValueKind, error) {
	columnTypes, err := e.getColumnKinds(ctx, tx, table)
	if err != nil {
		return nil, nil, err
	}
	var columnNames []string
	var columnKinds []api.ValueKind
	for columnName := range columnTypes {
		columnNames = append(columnNames, columnName)
	}
	sortStrings(columnNames)
	for _, columnName := range columnNames {
		columnKinds = append(columnKinds, columnTypes[columnName])
	}
	return columnNames, columnKinds, nil
}

// putRecord inserts a record into the database with the given timestamp
func (e *engine) putRecord(ctx context.Context, tx api.Transaction, record *api.Record, isReplication bool) error {

	// insert the record into the database
	if err := e.appendToHistory(ctx, tx, record, isReplication); err != nil {
		return err
	}

	// save to change stream
	if !isLocalTable(record.Table) {
		if err := e.appendToChangeStream(ctx, tx, record.Table, record.ID, record.Revision.String()); err != nil {
			return err
		}
	}

	// update the view
	if err := e.updateView(ctx, tx, record); err != nil {
		return err
	}

	return nil
}

// appendToHistory adds the given record to the history table
func (e *engine) appendToHistory(ctx context.Context, tx api.Transaction, record *api.Record, isReplication bool) error {

	// record must have an id
	if len(record.ID) == 0 {
		return fmt.Errorf("record id is empty")
	}

	// is the record is marked as not new, then it must have a revision
	if isReplication && record.Revision.IsEmpty() {
		return fmt.Errorf("record revision is empty")
	}

	// if the record is new, then the record revision must found in the database
	// it basically means that this record is the next revision of revision.Revision
	if !isReplication && !record.Revision.IsEmpty() {
		// find previous revision
		_, err := e.findRevision(ctx, tx, record.Table, record.ID, record.Revision)
		if err != nil {
			return err
		}
	}

	// generate a new revision is this is a new record
	if !isReplication {
		record.PreviousRevision = record.Revision
		newRevision := generateRevision(e.revisionGenerator, record.ID, record.Revision, record.Attributes)
		record.Revision = newRevision
	}

	// build the query
	sqlBuilder := &StringBuilder{}
	fields, placeHolders, params := getUpdateRecordSQLArgs(record.ID, record.PreviousRevision, record.Revision, record.Attributes)
	sqlBuilder.WriteString("INSERT INTO \"" + getHistoryTableName(record.Table) + "\" (" + joinStrings(fields, ", ") + ") VALUES (" + joinStrings(placeHolders, ", ") + ");")

	// execute the query
	if _, err := tx.Exec(ctx, sqlBuilder.String(), params); err != nil {
		return err
	}

	return nil
}

// updateView updates the reconciled view for the given record
func (e *engine) updateView(ctx context.Context, tx api.Transaction, rec *api.Record) error {

	// retrieve information about the table
	columnNames, columnTypes, err := e.getValueTypeInfoForTable(ctx, tx, rec.Table)
	if err != nil {
		return err
	}

	// build the query
	sqlQuery := `SELECT ` + joinStrings(columnNames, ",") + ` FROM "` + getHistoryTableName(rec.Table) + `" WHERE "` + api.KeyDeleted + `" = false AND "` + api.KeyRecordID + `" = ? ORDER BY "` + api.KeyRevision + `" DESC LIMIT 1;`

	// execute the query
	rows, err := tx.Query(ctx, sqlQuery, []interface{}{rec.ID})
	if err != nil {
		return err
	}
	defer closeRows(rows)

	// it there are no rows, that means that the last version of the record is deleted,
	// or that there was no history. We need to delete the record from the view to
	// reflect this
	if !rows.Next() {
		_, err := tx.Exec(ctx, `DELETE FROM "`+rec.Table+`" WHERE "`+api.KeyRecordID+`" = ?`, []interface{}{rec.ID})
		return err
	}

	// read the record from the query result
	data, err := rows.Read(columnTypes)
	if err != nil {
		return err
	}

	// build the update query
	var fields []string
	var placeholders []string
	var values []interface{}
	for k := range data {
		if k == api.KeyPrevision || k == api.KeyDeleted {
			continue
		}
		fields = append(fields, k)
		placeholders = append(placeholders, "?")
		value := data[k]
		values = append(values, &value)
	}
	sqlBuilder := &StringBuilder{}
	sqlBuilder.WriteString(`INSERT INTO "` + rec.Table + `" ("` + strings.Join(fields, `", "`) + `") VALUES (` + strings.Join(placeholders, ", ") + `)`)
	sqlBuilder.WriteString(` ON CONFLICT ("` + api.KeyRecordID + `") DO UPDATE SET `)
	var i int
	for _, k := range fields {
		if k == api.KeyRecordID {
			continue
		}
		if i != 0 {
			sqlBuilder.WriteString(", ")
		}
		sqlBuilder.WriteString(`"` + k + `" = excluded."` + k + `"`)
		i++
	}
	sqlBuilder.WriteString(`;`)

	// execute the query
	_, err = tx.Exec(ctx, sqlBuilder.String(), values)
	return err
}

// findRevision finds the revision of a record
func (e *engine) findRevision(ctx context.Context, tx api.Transaction, table string, id string, revision api.Revision) (*api.Record, error) {

	// get information about the table
	columnNames, columnTypes, err := e.getValueTypeInfoForTable(ctx, tx, table)
	if err != nil {
		return nil, err
	}

	// build the query
	sqlQuery := &StringBuilder{}
	sqlQuery.WriteString("SELECT " + joinStrings(columnNames, ",") + " FROM \"" + getHistoryTableName(table) + "\" WHERE \"" + api.KeyRecordID + "\" = ? AND \"" + api.KeyRevision + "\" = ?")

	// execute the query
	result, err := tx.Query(ctx, sqlQuery.String(), []interface{}{id, revision.String()})
	if err != nil {
		return nil, err
	}
	defer closeResult(result)

	// if there are no rows, return nil
	if !result.Next() {
		return nil, nil
	}

	// if there was an error while iterating, return it
	if err := result.Err(); err != nil {
		return nil, err
	}

	// get the values
	values, err := result.Read(columnTypes)
	if err != nil {
		return nil, err
	}

	// build the record
	ret, err := readInRecord(table, values)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

// createTable creates a table
func (e *engine) createTable(ctx context.Context, tx api.Transaction, table api.Table) error {

	// validate the table structure
	if err := validateTable(table); err != nil {
		return err
	}

	// check that a table with the same name does not already exist
	if err := e.checkTableDoesNotExist(ctx, tx, table.Name); err != nil {
		return err
	}

	// create the table
	if err := e.createTableInternal(ctx, tx, table, false); err != nil {
		return err
	}

	// create the history table
	if err := e.createTableInternal(ctx, tx, table, true); err != nil {
		return err
	}

	return nil
}

// createTableInternal creates a table
func (e *engine) createTableInternal(ctx context.Context, tx api.Transaction, table api.Table, isHistoryTable bool) error {

	// build the query
	sqlBuilder := &StringBuilder{}
	sqlBuilder.WriteString("CREATE TABLE IF NOT EXISTS \"")
	if isHistoryTable {
		sqlBuilder.WriteString(getHistoryTableName(table.Name))
	} else {
		sqlBuilder.WriteString(table.Name)
	}
	sqlBuilder.WriteString("\" (")
	columns := make([]api.Column, 0)

	// append the primary key column
	columns = append(columns, api.Column{
		Name: api.KeyRecordID,
		Type: "varchar",
	})

	// append the previous revision column if this is the history table
	if isHistoryTable {
		columns = append(columns, api.Column{
			Name: api.KeyPrevision,
			Type: "varchar",
			Constraints: []api.ColumnConstraint{
				{NotNull: &api.NotNullColumnConstraint{}},
			},
		})
	}

	// append the revision column
	columns = append(columns, api.Column{
		Name: api.KeyRevision,
		Type: "varchar",
		Constraints: []api.ColumnConstraint{
			{NotNull: &api.NotNullColumnConstraint{}},
		},
	})

	// build the constraints
	var tableConstraints []api.TableConstraint
	for _, c := range table.Constraints {
		tableConstraints = append(tableConstraints, c)
	}

	if isHistoryTable {
		// add the deleted column
		columns = append(columns, api.Column{
			Name:    api.KeyDeleted,
			Type:    "boolean",
			Default: "false",
			Constraints: []api.ColumnConstraint{
				{NotNull: &api.NotNullColumnConstraint{}},
			},
		})
		// add the primary key constraint for the history table
		tableConstraints = append(tableConstraints, api.TableConstraint{
			PrimaryKey: &api.PrimaryKeyTableConstraint{
				Columns: []string{api.KeyRecordID, api.KeyRevision},
			},
		})
	} else {
		// add the primary key constraint for the view table
		tableConstraints = append(tableConstraints, api.TableConstraint{
			PrimaryKey: &api.PrimaryKeyTableConstraint{
				Columns: []string{api.KeyRecordID},
			},
		})
	}

	for _, c := range table.Columns {
		columns = append(columns, c)
	}

	// write the SQL for the columns
	for i, column := range columns {
		writeColumnDefinition(sqlBuilder, column)
		if i < len(columns)-1 {
			sqlBuilder.WriteString(", ")
		}
	}

	// write the SQL for the constraints
	for _, c := range tableConstraints {
		sqlBuilder.WriteString(", ")
		writeConstraintDefinition(sqlBuilder, c)
	}

	sqlBuilder.WriteString(")")

	// execute the query
	if _, err := tx.Exec(ctx, sqlBuilder.String(), []interface{}{}); err != nil {
		return err
	}

	return nil
}

// checkTableDoesNotExist checks that a table does not exist
func (e *engine) checkTableDoesNotExist(ctx context.Context, tx api.Transaction, name string) error {
	if e.dialect == "sqlite" {
		sql := `SELECT "name" FROM "sqlite_master" WHERE "type" = 'table' AND "name" = ?`
		res, err := tx.Query(ctx, sql, []interface{}{name})
		if err != nil {
			return err
		}
		defer closeResult(res)
		if !res.Next() {
			return nil
		}
		if err := res.Err(); err != nil {
			return err
		}
		return api.NewTableAlreadyExistsErr(name)
	} else {
		return api.NewError(api.ErrCodeInternalError, "not implemented")
	}
}

// doTransaction is a helper method that wraps SQL operations within a transaction
func (e *engine) doTransaction(ctx context.Context, fn func(t api.Transaction) (interface{}, error)) (interface{}, error) {
	tr, err := e.txFactory(ctx)
	if err != nil {
		return nil, err
	}
	errored := false
	defer func() {
		if !errored {
			handleCommit(tr)
		} else {
			handleRollback(tr)
		}
	}()
	ret, err := fn(tr)
	if err != nil {
		errored = true
	}
	return ret, err
}

// initChangesTable creates the changes table if it does not exist
func (e *engine) initChangesTable(ctx context.Context, t api.Transaction) error {
	sqlStatement := `
CREATE TABLE IF NOT EXISTS "` + api.ChangeStreamTableName + `" (
	"` + api.KeyCSSequence + `" integer PRIMARY KEY AUTOINCREMENT,
    "` + api.KeyCSTableName + `" varchar NOT NULL,
	"` + api.KeyCSRecordID + `" varchar NOT NULL,
	"` + api.KeyCSRecordRevision + `" varchar NOT NULL
);`
	_, err := t.Exec(ctx, sqlStatement, nil)
	if err != nil {
		return err
	}
	return nil
}

// appendToChangeStream appends a change to the change stream
func (e *engine) appendToChangeStream(ctx context.Context, t api.Transaction, tableName string, recordId, revision string) error {
	sqlStatement := `
INSERT INTO "` + api.ChangeStreamTableName + `" ("` + api.KeyCSTableName + `", "` + api.KeyCSRecordID + `", "` + api.KeyCSRecordRevision + `")
VALUES (?, ?, ?)
`
	_, err := t.Exec(ctx, sqlStatement, []interface{}{tableName, recordId, revision})
	if err != nil {
		return err
	}
	return nil
}

func (e *engine) getTables(ctx context.Context, tx api.Transaction) (api.GetTablesResponse, error) {
	sqlStatement := `
select name from sqlite_master
where type = 'table'
and "name" not in ('sqlite_sequence', '` + api.ChangeStreamTableName + `')
and "name" not like '%_history'
`
	res, err := tx.Query(ctx, sqlStatement, nil)
	if err != nil {
		return api.GetTablesResponse{}, err
	}
	defer closeResult(res)
	var tableNames = make([]api.GetTablesResponseItem, 0)
	for res.Next() {
		var tableName string
		row, err := res.Read([]api.ValueKind{api.ValueKindString})
		if err != nil {
			return api.GetTablesResponse{}, err
		}
		tableName = row["name"].Str.ValueOrZero()
		tableNames = append(tableNames, api.GetTablesResponseItem{
			Name: tableName,
		})
	}
	if err := res.Err(); err != nil {
		return api.GetTablesResponse{}, err
	}
	return api.GetTablesResponse{
		Items: tableNames,
	}, nil
}

// generateRevision generates a new revision for a record
func generateRevision(revisionGenerator api.RevisionGenerator, recordId string, previousRevision api.Revision, data map[string]api.Value) api.Revision {
	revisionData := map[string]interface{}{
		api.KeyRecordID: recordId,
	}
	if previousRevision != api.EmptyRevision {
		revisionData[api.KeyPrevision] = previousRevision
	}
	for k, v := range data {
		revisionData[k] = v
	}
	return revisionGenerator.Generate(previousRevision.Num+1, revisionData)
}

func parseChangeStreamItem(rec api.Record) (api.ChangeItem, error) {
	var changeItem api.ChangeItem

	recordIDValue, err := rec.GetFieldValue(api.KeyCSRecordID)
	if err != nil {
		return api.ChangeItem{}, err
	}
	if recordIDValue.Kind != api.ValueKindString {
		return api.ChangeItem{}, fmt.Errorf("recordID is not a string")
	}
	changeItem.RecordID = recordIDValue.Str.ValueOrZero()

	tableNameValue, err := rec.GetFieldValue(api.KeyCSTableName)
	if err != nil {
		return api.ChangeItem{}, err
	}
	if tableNameValue.Kind != api.ValueKindString {
		return api.ChangeItem{}, fmt.Errorf("table name is not a string")
	}
	changeItem.TableName = tableNameValue.Str.ValueOrZero()

	revisionValue, err := rec.GetFieldValue(api.KeyCSRecordRevision)
	if err != nil {
		return api.ChangeItem{}, err
	}
	if revisionValue.Kind != api.ValueKindString {
		return api.ChangeItem{}, fmt.Errorf("revision is not a string")
	}
	revision, err := api.ParseRevision(revisionValue.Str.ValueOrZero())
	if err != nil {
		return api.ChangeItem{}, err
	}
	changeItem.RecordRevision = revision

	sequenceValue, err := rec.GetFieldValue(api.KeyCSSequence)
	if err != nil {
		return api.ChangeItem{}, err
	}
	if sequenceValue.Kind != api.ValueKindInt {
		return api.ChangeItem{}, fmt.Errorf("sequence is not an int")
	}
	changeItem.Sequence = sequenceValue.Int.ValueOrZero()

	return changeItem, nil
}

// readInRecord builds a record from a database result
func readInRecord(table string, data map[string]api.Value) (api.Record, error) {
	var err error
	var record = api.Record{
		Attributes: make(map[string]api.Value),
	}
	for columnName, columnValue := range data {
		switch columnName {
		case api.KeyRecordID:
			record.ID = columnValue.Str.ValueOrZero()
			continue
		case api.KeyRevision:
			record.Revision, err = api.ParseRevision(columnValue.Str.ValueOrZero())
			if err != nil {
				return api.Record{}, err
			}
		case api.KeyPrevision:
			value := columnValue.Str.ValueOrZero()
			if len(value) == 0 {
				record.PreviousRevision = api.EmptyRevision
			} else {
				record.PreviousRevision, err = api.ParseRevision(value)
				if err != nil {
					return api.Record{}, err
				}
			}
		default:
			record.Attributes[columnName] = columnValue
		}
	}
	record.Table = table
	return record, nil
}

func writeColumnDefinition(sqlBuilder *StringBuilder, column api.Column) {
	sqlBuilder.WriteString("\"" + column.Name + "\" " + column.Type)
	if column.Default != "" {
		sqlBuilder.WriteString(" DEFAULT " + column.Default)
	}
	for _, constraint := range column.Constraints {
		if constraint.NotNull != nil {
			sqlBuilder.WriteString(" NOT NULL")
		}
		if constraint.PrimaryKey != nil {
			sqlBuilder.WriteString(" PRIMARY KEY")
		}
	}
}

func writeConstraintDefinition(builder *StringBuilder, constraint api.TableConstraint) {
	if constraint.PrimaryKey != nil {
		builder.WriteString("PRIMARY KEY (")
		for i, column := range constraint.PrimaryKey.Columns {
			builder.WriteString("\"" + column + "\"")
			if i < len(constraint.PrimaryKey.Columns)-1 {
				builder.WriteString(", ")
			}
		}
		builder.WriteString(")")
	}
}

func handleCommit(tr api.Transaction) {
	if err := tr.Commit(); err != nil {
		fmt.Printf("Error on commit: %v\n", err)
		handleRollback(tr)
	}
}

func handleRollback(tr api.Transaction) {
	if err := tr.Rollback(); err != nil {
		fmt.Printf("error while rolling back transaction: %v", err)
	}
}

func closeRows(rows api.ResultReader) {
	if err := rows.Close(); err != nil {
		fmt.Printf("error closing rows: %v\n", err)
	}
}

func closeResult(result api.ResultReader) {
	if err := result.Close(); err != nil {
		fmt.Printf("error closing result: %v\n", err)
	}
}

func isLocalTable(tableName string) bool {
	if len(tableName) <= 6 {
		return false
	}
	return tableName[len(tableName)-6:] == "local_"
}

func getUpdateRecordSQLArgs(id string, previousRevision, currentRevision api.Revision, data map[string]api.Value) (fields, placeholders []string, values []interface{}) {
	fields = append(fields, `"`+api.KeyRecordID+`"`, `"`+api.KeyRevision+`"`, `"`+api.KeyPrevision+`"`)
	placeholders = append(placeholders, "?", "?", "?")
	values = append(values, id, currentRevision.String(), previousRevision.String())

	for field, value := range data {
		fields = append(fields, `"`+field+`"`)
		placeholders = append(placeholders, "?")
		values = append(values, &value)
	}
	return
}

func getHistoryTableName(tableName string) string {
	return tableName + "_history"
}
