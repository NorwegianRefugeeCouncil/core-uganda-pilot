package store

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/logging"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

// actionContext is the boilerplate code for all database operations
// it includes creating a logger, getting a database connection, setting the database context.Context, and
// also logging the execution time of a database operation.
// ctx is the request context
//
// factory is the database Factory
// storeName is the name of the store, for logging purposes
// actionName is the name of the action being done, for logging purposes
// fields is an arbitrary list of fields for logging
//
// Returns
// context.Context the new populated request context with the logger
// gorm.DB the database connection
// zap.Logger the logger
// doneFunc the function that logs the execution time
func actionContext(ctx context.Context, factory Factory, storeName, actionName string, fields ...zap.Field) (context.Context, *gorm.DB, *zap.Logger, doneFunc func(), error) {
	ctx, l := logging.NewStoreLogger(ctx, storeName, actionName, fields...)
	l.Debug("getting database connection")
	db, err := factory.Get()
	if err != nil {
		l.Error("failed to get database connection", zap.Error(err))
		return nil, nil, nil, nil, err
	}
	db = db.WithContext(ctx)
	start := time.Now()
	return ctx, db, l, func() {
		l.Debug(fmt.Sprintf("%s.%s completed", storeName, actionName), zap.Duration("duration", time.Now().Sub(start)))
	}, nil
}
