package sessionmanager

import (
	"context"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"time"
)

type Options struct {
	ErrorFunc func(http.ResponseWriter, *http.Request, error)
}

type Store interface {
	Get(ctx context.Context, key string) interface{}
	GetBool(ctx context.Context, key string) bool
	GetString(ctx context.Context, key string) string
	GetBytes(ctx context.Context, key string) []byte
	GetFloat(ctx context.Context, key string) float64
	GetTime(ctx context.Context, key string) time.Time
	GetInt(ctx context.Context, key string) int
	Pop(ctx context.Context, key string) interface{}
	PopBool(ctx context.Context, key string) bool
	PopString(ctx context.Context, key string) string
	PopBytes(ctx context.Context, key string) []byte
	PopFloat(ctx context.Context, key string) float64
	PopTime(ctx context.Context, key string) time.Time
	PopInt(ctx context.Context, key string) int
	Put(ctx context.Context, key string, val interface{})
	LoadAndSave(next http.Handler) http.Handler
	Clear(ctx context.Context) error
	Destroy(ctx context.Context) error
	Commit(ctx context.Context) (string, time.Time, error)
}

func New(pool *redis.Pool, options Options) Store {
	sessionManager := scs.New()
	sessionManager.Store = redisstore.New(pool)
	if options.ErrorFunc != nil {
		sessionManager.ErrorFunc = options.ErrorFunc
	}
	return sessionManager
}
