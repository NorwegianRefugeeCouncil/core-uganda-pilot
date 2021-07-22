package sessionmanager

import (
	"encoding/gob"
	"fmt"
	"github.com/boj/redistore"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Options struct {
	ErrorFunc func(http.ResponseWriter, *http.Request, error)
}

type Store interface {
	AddNotification(req *http.Request, w http.ResponseWriter, notification *Notification) error
	ConsumeNotifications(req *http.Request) ([]*Notification, error)
	Get(req *http.Request) (*sessions.Session, error)
	GetString(req *http.Request, key string) (string, error)
}

type RedisSessionManager struct {
	sessions.Store
}

func (r *RedisSessionManager) GetString(req *http.Request, key string) (string, error) {
	session, err := r.Get(req)
	if err != nil {
		return "", err
	}

	strIntf, ok := session.Values[key]
	if !ok {
		return "", fmt.Errorf("key %s not found in session", key)
	}

	str, ok := strIntf.(string)
	if !ok {
		return "", fmt.Errorf("key %s is not a string", key)
	}

	return str, nil
}

func (r *RedisSessionManager) Get(req *http.Request) (*sessions.Session, error) {
	session, err := r.Store.Get(req, varSession)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get session")
		return nil, err
	}
	return session, nil
}

func init() {
	gob.Register(&Notification{})
}

func New(redisStore *redistore.RediStore) (Store, error) {
	return &RedisSessionManager{
		redisStore,
	}, nil
}

// Notification contains "flash" messages shown to the user.
// The theme field should correspond to one of the bootstrap theme colors; eg. "success", "warning",
// "danger", etc. https://getbootstrap.com/docs/5.0/customize/color/
type Notification struct {
	Message string
	Theme   string
}

const (
	varSession       = "session"
	varNotifications = "notifications"
)

func (r RedisSessionManager) AddNotification(req *http.Request, w http.ResponseWriter, notification *Notification) error {
	session, err := r.Store.Get(req, varNotifications)
	if err != nil {
		return err
	}
	session.AddFlash(notification)
	if err := session.Save(req, w); err != nil {
		return err
	}
	return nil
}

func (r RedisSessionManager) ConsumeNotifications(req *http.Request) ([]*Notification, error) {
	session, err := r.Store.Get(req, varNotifications)
	if err != nil {
		return nil, err
	}
	flashes := session.Flashes()
	var notifications []*Notification
	for _, flash := range flashes {
		flashNotification, ok := flash.(*Notification)
		if ok {
			notifications = append(notifications, flashNotification)
		}
	}
	return notifications, nil
}
