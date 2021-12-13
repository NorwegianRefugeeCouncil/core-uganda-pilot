package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDynamicFactory(t *testing.T) {
	getDsn := func() (string, error) {
		return "sqlite://file::memory:", nil
	}
	factory, err := NewDynamicFactory(getDsn)
	if !assert.NoError(t, err) {
		return
	}
	db, err := factory.Get()
	if !assert.NoError(t, err) {
		return
	}
	sqlDB, err := db.DB()
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NoError(t, sqlDB.Ping()) {
		return
	}
}

func TestDynamicFactoryReload(t *testing.T) {
	var dsn1 = "sqlite://test-1.db"
	var dsn2 = "sqlite://test-2.db"
	var dsn = &dsn1
	getDsn := func() (string, error) {
		return *dsn, nil
	}
	factory, err := NewDynamicFactory(getDsn)
	if !assert.NoError(t, err) {
		return
	}
	first, err := factory.Get()
	if !assert.NoError(t, err) {
		return
	}
	firstAgain, err := factory.Get()
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, first, firstAgain)
	dsn = &dsn2
	second, err := factory.Get()
	if !assert.NoError(t, err) {
		return
	}
	assert.NotEqual(t, first, second)
}

func TestExtractDBOptions(t *testing.T) {

	tests := []struct {
		name          string
		dsn           string
		expectDsn     string
		expectOptions sqlDbOptions
		expectErr     bool
	}{
		{
			name:      "max_open_conns",
			dsn:       "sqlite://bla?max_open_conns=10",
			expectDsn: "sqlite://bla",
			expectOptions: sqlDbOptions{
				maxOpenConns: 10,
			},
		}, {
			name:      "max_idle_conns",
			dsn:       "sqlite://bla?max_idle_conns=10",
			expectDsn: "sqlite://bla",
			expectOptions: sqlDbOptions{
				maxIdleConns: 10,
			},
		}, {
			name:      "conn_max_idle_time",
			dsn:       "sqlite://bla?conn_max_idle_time=10",
			expectDsn: "sqlite://bla",
			expectOptions: sqlDbOptions{
				connMaxIdleTime: 10,
			},
		}, {
			name:      "conn_max_lifetime",
			dsn:       "sqlite://bla?conn_max_lifetime=10",
			expectDsn: "sqlite://bla",
			expectOptions: sqlDbOptions{
				connMaxLifetime: 10,
			},
		}, {
			name:      "allAtOnce",
			dsn:       "sqlite://bla?sslMode=disable&max_open_conns=10&max_idle_conns=10&conn_max_idle_time=10&conn_max_lifetime=10",
			expectDsn: "sqlite://bla?sslMode=disable",
			expectOptions: sqlDbOptions{
				connMaxIdleTime: 10,
				connMaxLifetime: 10,
				maxIdleConns:    10,
				maxOpenConns:    10,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dsn, opts, err := extractDbOptions(test.dsn)
			if test.expectErr {
				assert.Error(t, err)
				return
			}
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, test.expectOptions, opts)
			assert.Equal(t, test.expectDsn, dsn)
		})
	}

}
