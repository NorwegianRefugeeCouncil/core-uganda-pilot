package authn

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		clearSession := func() {
			userSession, err := h.sessionStore.New(req, constants.SessionKey)
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

		userSession, err := h.sessionStore.Get(req, constants.SessionKey)
		if err != nil {
			logrus.WithError(err).Errorf("failed to retrieve user session: %s", err)
			clearSession()
			utils.ErrorResponse(w, meta.NewInternalServerError(errors.New("failed to retrieve session")))
			return
		}

		userSession.Values[constants.SessionState] = state
		if err := userSession.Save(req, w); err != nil {
			logrus.WithError(err).Errorf("failed to save user session: %s", err)
			clearSession()
			utils.ErrorResponse(w, meta.NewInternalServerError(errors.New("failed to save session")))
			return
		}

		authCodeURL := h.oauth2Config.AuthCodeURL(state)
		http.Redirect(w, req, authCodeURL, http.StatusTemporaryRedirect)
	}
}

func (h *Handler) RestfulLogin(request *restful.Request, response *restful.Response) {
	handler := h.Login()
	handler(response.ResponseWriter, request.Request)
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
