package store

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/nrc-no/core/pkg/logging"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	keyMaxOpenConns = "max_open_conns"
	keyMaxIdleConns = "max_idle_conns"
	keyMaxIdleTime  = "conn_max_idle_time"
	keyMaxLifetime  = "conn_max_lifetime"
)

// Factory is the store.Factory that returns an instance of a gorm.DB
// This is useful because we can implement logic that allows us to
// renew a database connection, for example, when a connection string changes
// (credential rotation)
type Factory interface {
	Get() (*gorm.DB, error)
}

// NewFactory returns a new instance of Factory
func NewFactory(dsn string) (Factory, error) {
	sqlDialect, opts, err := getSqlDialect(dsn)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(sqlDialect)
	if err != nil {
		return nil, err
	}
	if err := applyDbOptions(db, opts); err != nil {
		return nil, err
	}
	return &factory{
		db: db,
	}, nil
}

// NewDynamicFactory returns a new dynamicFactory
func NewDynamicFactory(getDsn dsnFn) (Factory, error) {
	return &dynamicFactory{
		lock:   &sync.Mutex{},
		getDsn: getDsn,
	}, nil
}

// factory is the implementation of Factory that will use
// a static db connection
type factory struct {
	// db is the gorm.DB database connection
	db *gorm.DB
}

// dynamicFactory is an implementation of Factory that can
// reload connections when the dsn change
type dynamicFactory struct {
	// lock for controlled access to the database connection
	lock *sync.Mutex
	// prevDsnHash is the hash of the last used dsn
	prevDsnHash []byte
	// getDsn is a function that returns the current dsn
	getDsn dsnFn
	// db is the current gorm.DB database connection
	db *gorm.DB
}

// Get implements Factory.Get
func (f factory) Get() (*gorm.DB, error) {
	return f.db, nil
}

type dsnFn func() (string, error)

// Get implements Factory.Get
func (f *dynamicFactory) Get() (*gorm.DB, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	l := logging.NewLogger(context.Background())

	dsn, err := f.getDsn()
	if err != nil {
		return nil, err
	}

	// calculating if the dsn hash has changed. If yes, that means
	// that we must create a new connection. Otherwise, return
	// the previously used connection
	dsnHash := sha256.New().Sum([]byte(dsn))
	if f.db != nil && bytes.Equal(dsnHash, f.prevDsnHash) {
		l.Debug("reusing same database factory configuration")
		return f.db, nil
	}

	l.Debug("factory configuration changed")

	// getting sql dialect
	sqlDialect, sqlOptions, err := getSqlDialect(dsn)
	if err != nil {
		return nil, err
	}

	// opening a connection
	db, err := gorm.Open(sqlDialect)
	if err != nil {
		return nil, err
	}

	// applying sqlDbOptions
	if err := applyDbOptions(db, sqlOptions); err != nil {
		return nil, err
	}

	f.db = db
	f.prevDsnHash = dsnHash
	return db, nil
}

// applyDbOptions applies the sqlDbOptions on the given gorm.DB
func applyDbOptions(db *gorm.DB, opts sqlDbOptions) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	if opts.maxIdleConns != 0 {
		sqlDB.SetMaxIdleConns(opts.maxIdleConns)
	}
	if opts.maxOpenConns != 0 {
		sqlDB.SetMaxOpenConns(opts.maxOpenConns)
	}
	if opts.connMaxIdleTime != 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(opts.connMaxIdleTime) * time.Second)
	}
	if opts.connMaxLifetime != 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(opts.connMaxLifetime) * time.Second)
	}
	return nil
}

// sqlDbOptions represent the options extracted from the
// dsn that must be configured on the sql.DB connection.
type sqlDbOptions struct {
	// connMaxIdleTime is the sql connection max idle time
	connMaxIdleTime int
	// connMaxLifetime is the sql connection max lifetime
	connMaxLifetime int
	// maxIdleConns is the maximum number of idle connections
	maxIdleConns int
	// maxOpenConns is the maximum number of open connections
	maxOpenConns int
}

// getSqlDialect returns the gorm.Dialector and sqlDbOptions for the given dsn
func getSqlDialect(dsn string) (gorm.Dialector, sqlDbOptions, error) {
	if strings.HasPrefix(dsn, "sqlite://") {
		return getSqliteDialect(dsn)
	} else if strings.HasPrefix(dsn, "postgres") {
		return getPostgresDialect(dsn)
	} else {
		return nil, sqlDbOptions{}, fmt.Errorf("unknown database dialect")
	}
}

// getPostgresDialect returns a gorm.Dialector for a postgres connection,
// as well as a sqlDbOptions
func getPostgresDialect(dsn string) (gorm.Dialector, sqlDbOptions, error) {
	connDsn, opts, err := extractDbOptions(dsn)
	if err != nil {
		return nil, sqlDbOptions{}, err
	}
	sqlDialect := postgres.Open(connDsn)
	return sqlDialect, opts, nil
}

// getSqliteDialect returns the gorm.Dialector for a sqlite connection
// as well as a sqlDbOptions
func getSqliteDialect(dsn string) (gorm.Dialector, sqlDbOptions, error) {
	connDsn, opts, err := extractDbOptions(dsn)
	if err != nil {
		return nil, sqlDbOptions{}, err
	}
	connDsn = strings.TrimPrefix(connDsn, "sqlite://")
	sqlDialect := sqlite.Open(connDsn)
	return sqlDialect, opts, nil
}

// extractDbOptions will retrieve the sql.DB configuration parameters from a dsn.
// It will also remove these query parameters from the dsn, since the driver
// does not support these parameters directly, but they are used
// to subsequently configure the sql.DB connection
func extractDbOptions(dsn string) (string, sqlDbOptions, error) {
	dsnUrl, err := url.Parse(dsn)
	if err != nil {
		return "", sqlDbOptions{}, err
	}
	qry := dsnUrl.Query()
	opts := &sqlDbOptions{}

	if err := parseInt(qry, keyMaxOpenConns, &opts.maxOpenConns); err != nil {
		return "", sqlDbOptions{}, err
	}

	if err := parseInt(qry, keyMaxIdleConns, &opts.maxIdleConns); err != nil {
		return "", sqlDbOptions{}, err
	}

	if err := parseInt(qry, keyMaxIdleTime, &opts.connMaxIdleTime); err != nil {
		return "", sqlDbOptions{}, err
	}
	if err := parseInt(qry, keyMaxLifetime, &opts.connMaxLifetime); err != nil {
		return "", sqlDbOptions{}, err
	}
	qry.Del(keyMaxOpenConns)
	qry.Del(keyMaxIdleConns)
	qry.Del(keyMaxIdleTime)
	qry.Del(keyMaxLifetime)
	dsnUrl.RawQuery = qry.Encode()
	return dsnUrl.String(), *opts, nil
}

// parseInt is a utility method to populate an int field from an
// url.Values query parameter. Used by extractDbOptions to populate
func parseInt(qry url.Values, key string, value *int) error {
	qryParam := qry.Get(key)
	if len(qryParam) != 0 {
		qryParamValue, err := strconv.Atoi(qryParam)
		if err != nil {
			return err
		}
		*value = qryParamValue
	}
	return nil
}
