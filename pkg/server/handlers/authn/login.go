package authn

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/emicklei/go-restful/v3"
	"github.com/gorilla/securecookie"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
)

func (h *Handler) Login(sessionKey string, silent bool, redirectUri string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		clearSession := func() {
			userSession, err := h.sessionStore.New(req, sessionKey)
			if err != nil {
				return
			}
			_ = userSession.Save(req, w)
		}

		state, err := createStateVariable()
		if err != nil {
			logrus.WithError(err).Errorf("failed ot create state variable: %s", err)
			clearSession()
			utils.ErrorResponse(w, meta.NewInternalServerError(errors.New("failed to create state variable")))
			return
		}

		userSession, err := h.sessionStore.Get(req, sessionKey)
		securecookie.MultiError{}.IsDecode()
		if err != nil {
			if cookieErr, ok := err.(securecookie.MultiError); ok {
				if !cookieErr.IsDecode() {
					logrus.WithError(err).Errorf("failed to retrieve user session: %s", err)
					clearSession()
					return
				}
			}
			if err := userSession.Save(req, w); err != nil {
				logrus.WithError(err).Errorf("failed to clear user session: %s", err)
				clearSession()
				return
			}
		}

		userSession.Values[constants.SessionState] = state
		userSession.Values[constants.SessionDesiredURL] = redirectUri
		if err := userSession.Save(req, w); err != nil {
			logrus.WithError(err).Errorf("failed to save user session: %s", err)
			clearSession()
			utils.ErrorResponse(w, meta.NewInternalServerError(errors.New("failed to save session")))
			return
		}

		var authCodeOptions []oauth2.AuthCodeOption
		if silent {
			authCodeOptions = append(authCodeOptions, oauth2.SetAuthURLParam("prompt", "none"))
		}
		authCodeURL := h.oauth2Config.AuthCodeURL(state, authCodeOptions...)

		http.Redirect(w, req, authCodeURL, http.StatusTemporaryRedirect)

	}
}

func (h *Handler) RestfulLogin(sessionKey string, silent bool) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		redirectUri := req.QueryParameter("redirect_uri")
		h.Login(sessionKey, silent, redirectUri)(res.ResponseWriter, req.Request)
	}
}

func createStateVariable() (string, error) {
	bts, err := generateBytes(32)
	if err != nil {
		return "", err
	}
	state := base64.StdEncoding.EncodeToString(bts)
	return state, nil
}

func generateBytes(count int) ([]byte, error) {
	b := make([]byte, count)
	_, err := rand.Read(b)
	if err != nil {
		return nil, meta.NewInternalServerError(err)
	}
	return b, nil
}
