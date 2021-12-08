package testutils

import (
	"fmt"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/juju/fslock"
	"gopkg.in/matryer/try.v1"
	"time"
)

// lock is a file lock. This is needed because the go tests do not actually
// run on the same thread, and sometimes perhaps even not on the same process.
// We must use a file-based lock so that we don't try to instantiate/teardown
// postgres concurrently.
var lock = fslock.New("/tmp/embedded-postgres")

// LockPostgres will try to acquire the lock on the file
func LockPostgres() error {
	return lock.Lock()
}

// UnlockPostgres will try to release the lock on the file
func UnlockPostgres() error {
	return lock.Unlock()
}

// TryGetPostgres will try to start an instance of postgres
func TryGetPostgres() (func(), error) {
	if err := LockPostgres(); err != nil {
		return func() {}, err
	}

	done := func() {
		UnlockPostgres()
	}

	var pg *embeddedpostgres.EmbeddedPostgres
	err := try.Do(func(attempt int) (bool, error) {
		pg = embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().Port(15432))
		err := pg.Start()
		if err != nil {
			fmt.Println("failed to start postgres: " + err.Error())
			time.Sleep(time.Duration(attempt) * time.Second)
		}
		return attempt < 5, err
	})

	if err != nil {
		return done, err
	}

	done = func() {
		UnlockPostgres()
		pg.Stop()
	}

	return done, nil

}
