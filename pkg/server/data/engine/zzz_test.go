package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nrc-no/core/pkg/server/data/api"
	"github.com/nrc-no/core/pkg/server/data/test"
	"github.com/nrc-no/core/pkg/server/data/utils"
)

// Bench is a testing utility that makes it easy to setup a test database
// with all the dependencies
type Bench struct {
	DBName            string
	Engine            api.Engine
	DB                *sqlx.DB
	Recorder          *test.DBRecorder
	Clock             *test.MockClock
	UUIDGenerator     *test.MockUUIDGenerator
	RevisionGenerator *utils.Md5RevGenerator
	Cancel            context.CancelFunc
	Ctx               context.Context
}

// TearDown closes the database connection and cleans up the test database
func (b *Bench) TearDown() error {
	if b.Cancel != nil {
		b.Cancel()
	}
	if b.DB != nil {
		if err := b.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Reset resets the database to a clean state
func (b *Bench) Reset() error {

	var err error
	if err = b.TearDown(); err != nil {
		return err
	}

	b.Ctx, b.Cancel = context.WithCancel(context.Background())

	// create the database connection
	b.DB, err = sqlx.ConnectContext(b.Ctx, "sqlite3", b.DBName)
	if err != nil {
		return err
	}

	// drop all data
	if err := test.DropAll(b.DB); err != nil {
		panic(err)
	}

	b.Engine, err = NewEngine(
		context.Background(),
		TxFactory(b.DB, b.Recorder),
		b.UUIDGenerator,
		b.RevisionGenerator,
		b.Clock,
		"sqlite",
	)
	return err
}

// NewTestBench creates a new test bench
func NewTestBench(dbName string) *Bench {
	b := &Bench{
		DBName:   dbName,
		Recorder: &test.DBRecorder{},
		Clock: &test.MockClock{
			TheTime: time.Now().Unix(),
		},
		UUIDGenerator: &test.MockUUIDGenerator{
			ReturnUUID: "12345678-1234-1234-1234-123456789012",
		},
		RevisionGenerator: &utils.Md5RevGenerator{},
	}
	return b
}

// TxFactory creates a new transaction factory
func TxFactory(conn *sqlx.DB, recorder *test.DBRecorder) func(ctx context.Context) (api.Transaction, error) {
	return func(ctx context.Context) (api.Transaction, error) {
		return utils.NewTransaction(ctx, conn, func(qry string, args []interface{}) {
			fmt.Println(qry, args)
			recorder.Record(qry, args)
		})
	}
}
