package sessionmanager

import (
	"context"
	"encoding/gob"
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
	AddNotification(ctx context.Context, notification *Notification)
	ConsumeNotifications(ctx context.Context) []*Notification
}

type RedisSessionManager struct {
	*scs.SessionManager
}

func init() {
	gob.Register([]*Notification{})
}

func New(pool *redis.Pool, options Options) Store {
	sessionManager := scs.New()
	sessionManager.Store = redisstore.New(pool)
	if options.ErrorFunc != nil {
		sessionManager.ErrorFunc = options.ErrorFunc
	}
	return RedisSessionManager{
		SessionManager: sessionManager,
	}
}

// Notification contains "flash" messages shown to the user.
// The theme field should correspond to one of the bootstrap theme colors; eg. "success", "warning",
// "danger", etc. https://getbootstrap.com/docs/5.0/customize/color/
type Notification struct {
	Message string
	Theme   string
}

func (r RedisSessionManager) AddNotification(ctx context.Context, notification *Notification) {
	notifs, ok := r.Get(ctx, "notifications").([]*Notification)
	if !ok {
		notifs = []*Notification{}
	}
	notifs = append(notifs, notification)
	r.Put(ctx, "notifications", notifs)
}
func (r RedisSessionManager) ConsumeNotifications(ctx context.Context) []*Notification {
	notifs, ok := r.Pop(ctx, "notifications").([]*Notification)
	if !ok {
		notifs = []*Notification{}
	}
	return notifs
}
