package api

import "context"

type GetRecordRequest struct {
	TableName string
	RecordID  string
	Revision  Revision
}

type GetRecordsRequest struct {
	TableName string   `json:"table"`
	RecordIDs []string `json:"recordIds"`
	Revisions bool     `json:"revisions"`
}

type GetChangesRequest struct {
	Since int64
}

type GetTablesRequest struct{}

type GetTableRequest struct {
	TableName string `json:"tableName"`
}

type GetTablesResponseItem struct {
	Name string `json:"name"`
}

type GetTablesResponse struct {
	Items []GetTablesResponseItem `json:"items"`
}

type ReadInterface interface {
	// GetRecord gets a single record from the database.
	GetRecord(ctx context.Context, request GetRecordRequest) (Record, error)
	// GetRecords gets multiple records from the database.
	GetRecords(ctx context.Context, request GetRecordsRequest) (RecordList, error)
	// GetChanges gets a change stream for a table
	GetChanges(ctx context.Context, request GetChangesRequest) (Changes, error)
	// GetTables gets a list of tables in the database
	GetTables(ctx context.Context, request GetTablesRequest) (GetTablesResponse, error)
	// GetTable gets a table definition
	GetTable(ctx context.Context, request GetTableRequest) (Table, error)
}

type PutRecordRequest struct {
	Record        Record
	IsReplication bool
}

type WriteInterface interface {
	// PutRecord puts a single record inside the database.
	PutRecord(ctx context.Context, request PutRecordRequest) (Record, error)
	// CreateTable creates a new table in the database.
	CreateTable(ctx context.Context, table Table) (Table, error)
}

type Engine interface {
	ReadInterface
	WriteInterface
}

type TxFactory func(ctx context.Context) (Transaction, error)

type Transaction interface {
	Query(ctx context.Context, query string, args []interface{}) (ResultReader, error)
	Exec(ctx context.Context, query string, args []interface{}) (interface{}, error)
	Commit() error
	Rollback() error
}

// Rand generates random bytes
type Rand interface {
	// Read puts random bytes into the given buffer.
	Read(b []byte) (n int, err error)
}

// UUIDGenerator generates UUIDs
type UUIDGenerator interface {
	// Generate generates a UUID
	Generate() (string, error)
}

// ResultReader reads results from a query
type ResultReader interface {
	// Next returns the next result
	Next() bool
	// Read returns the current result
	Read(columnKinds []ValueKind) (map[string]Value, error)
	// Close closes the reader
	Close() error
	// Err returns the last error
	Err() error
}

// RevisionGenerator generates revision hashes
type RevisionGenerator interface {
	// Generate generates a revision Hash
	Generate(num int, data map[string]interface{}) Revision
}

// Clock provides a clock
type Clock interface {
	// Now returns the current time
	Now() int64
}
