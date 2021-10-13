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
	AddFlash(req *http.Request, w http.ResponseWriter, notification *FlashMessage) error
	ConsumeFlashes(req *http.Request, w http.ResponseWriter) ([]*FlashMessage, error)
	Get(req *http.Request) (*sessions.Session, error)
	GetString(req *http.Request, key string) (string, error)
	FindString(req *http.Request, key string) (string, bool)
}

type RedisSessionManager struct {
	sessions.Store
}

func (r *RedisSessionManager) FindString(req *http.Request, key string) (string, bool) {
	str, err := r.GetString(req, key)
	if err != nil {
		return "", false
	}

	return str, true
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

		return session, err
	}

	return session, nil
}

func New(redisStore *redistore.RediStore) (Store, error) {
	return &RedisSessionManager{
		redisStore,
	}, nil
}

// FlashMessage contains "flash" messages shown to the user.
// The theme field should correspond to one of the bootstrap theme colors; eg. "success", "warning",
// "danger", etc. https://getbootstrap.com/docs/5.0/customize/color/
type FlashMessage struct {
	Message string
	Theme   string
}

// We must register the data type so that it can be encoded/decoded
func init() {
	gob.Register(&FlashMessage{})
}

const (
	varSession = "session"
	varFlashes = "flashes"
)

func (r RedisSessionManager) AddFlash(req *http.Request, w http.ResponseWriter, flash *FlashMessage) error {
	session, err := r.Store.Get(req, varFlashes)
	if err != nil {
		return err
	}

	session.AddFlash(flash)

	return session.Save(req, w)
}

func (r RedisSessionManager) ConsumeFlashes(req *http.Request, w http.ResponseWriter) ([]*FlashMessage, error) {
	session, err := r.Store.Get(req, varFlashes)
	if err != nil {
		return nil, err
	}

	var flashes []*FlashMessage

	if values := session.Flashes(); len(flashes) > 0 {
		for _, val := range values {
			var flash = &FlashMessage{}
			flash, ok := val.(*FlashMessage)
			if !ok {
				// TODO handle unexpected type
				continue
			}

			flashes = append(flashes, flash)
		}
	}

	return flashes, session.Save(req, w)
}
