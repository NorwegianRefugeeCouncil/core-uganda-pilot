package store

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/logging"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

func actionContext(ctx context.Context, factory Factory, storeName, actionName string, fields ...zap.Field) (context.Context, *gorm.DB, *zap.Logger, func(), error) {
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
