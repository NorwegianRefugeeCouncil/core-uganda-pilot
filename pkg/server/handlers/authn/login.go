package authn

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/nrc-no/core/pkg/utils/sets"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"net/http"
)

func (h *Handler) Login(
	sessionKey string,
	silent bool,
	redirectUri string,
	defaultRedirectUri string,
	allowedRedirectUris sets.String,
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		l := logging.NewLogger(ctx)

		clearSession := func() {
			userSession, err := h.sessionStore.New(req, sessionKey)
			if err != nil {
				return
			}
			_ = userSession.Save(req, w)
		}

		if len(redirectUri) == 0 {
			redirectUri = defaultRedirectUri
		}
		if !allowedRedirectUris.Has(redirectUri) {
			l.Warn("illegal redirect uri passed to login endpoint", zap.String("redirect_uri", redirectUri))
			utils.ErrorResponse(w, meta.NewBadRequest("illegal redirect uri"))
			return
		}

		state, err := createStateVariable()
		if err != nil {
			l.Error("failed to create state variable", zap.Error(err))
			clearSession()
			utils.ErrorResponse(w, meta.NewInternalServerError(errors.New("failed to create state variable")))
			return
		}

		userSession, err := getSession(w, req, h.sessionStore, sessionKey)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		userSession.Values[constants.SessionState] = state
		userSession.Values[constants.SessionDesiredURL] = redirectUri
		if err := userSession.Save(req, w); err != nil {
			l.Error("failed to save user session", zap.Error(err))
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

func (h *Handler) RestfulLogin(
	sessionKey string,
	silent bool,
	defaultRedirectUri string,
	allowedRedirectUris sets.String,
) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		redirectUri := req.QueryParameter("redirect_uri")
		h.Login(sessionKey, silent, redirectUri, defaultRedirectUri, allowedRedirectUris)(res.ResponseWriter, req.Request)
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
