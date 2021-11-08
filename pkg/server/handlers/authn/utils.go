package authn

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/logging"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

func strFromSession(session *sessions.Session, key string) (string, bool) {
	valueIntf, ok := session.Values[key]
	if !ok {
		return "", false
	}
	value, ok := valueIntf.(string)
	if !ok {
		return "", false
	}
	return value, true
}

func timeFromSession(session *sessions.Session, key string) (*time.Time, bool) {
	valueIntf, ok := session.Values[key]
	if !ok {
		return nil, false
	}
	value, ok := valueIntf.(*time.Time)
	if !ok {
		return nil, false
	}
	return value, true
}

func tokenFromSession(userSession *sessions.Session) (*StoredToken, bool) {
	idToken, ok := strFromSession(userSession, constants.SessionIDToken)
	if !ok {
		return nil, false
	}
	accessToken, ok := strFromSession(userSession, constants.SessionAccessToken)
	if !ok {
		return nil, false
	}
	refreshToken, ok := strFromSession(userSession, constants.SessionRefreshToken)
	if !ok {
		return nil, false
	}
	tokenType, ok := strFromSession(userSession, constants.SessionTokenType)
	if !ok {
		return nil, false
	}
	tokenExpiry, ok := timeFromSession(userSession, constants.SessionTokenExpiry)
	if !ok {
		return nil, false
	}
	return &StoredToken{
		Token: &oauth2.Token{
			AccessToken:  accessToken,
			TokenType:    tokenType,
			RefreshToken: refreshToken,
			Expiry:       *tokenExpiry,
		},
		IDToken: idToken,
	}, true
}

func clearSession(w http.ResponseWriter, req *http.Request, sessionStore sessions.Store, sessionKey string) func() {
	return func() {
		userSession, err := sessionStore.New(req, sessionKey)
		if err != nil {
			return
		}
		_ = userSession.Save(req, w)
	}
}

func getSession(
	w http.ResponseWriter,
	req *http.Request,
	sessionStore sessions.Store,
	sessionKey string,
) (*sessions.Session, error) {

	ctx := req.Context()
	l := logging.NewLogger(ctx)

	l.Debug("getting user session")
	userSession, err := sessionStore.Get(req, sessionKey)
	securecookie.MultiError{}.IsDecode()
	if err != nil {
		if cookieErr, ok := err.(securecookie.MultiError); ok {
			if !cookieErr.IsDecode() {
				l.Error("failed to retrieve user session", zap.Error(err))
				return nil, meta.NewInternalServerError(err)
			}
		}
		if err := userSession.Save(req, w); err != nil {
			l.Error("failed to clear user session", zap.Error(err))
			return nil, meta.NewInternalServerError(err)
		}
	}

	return userSession, nil

}
