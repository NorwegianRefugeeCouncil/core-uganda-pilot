package authn

import (
	"errors"
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ExposedSession struct {
	Active           bool      `json:"active"`
	Expiry           time.Time `json:"expiry,omitempty"`
	ExpiresInSeconds int       `json:"expiresInSeconds"`
	Subject          string    `json:"subject,omitempty"`
	Username         string    `json:"username,omitempty"`
}

func (h *Handler) Session(sessionKey string, hydraAdmin admin.ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		l := logging.NewLogger(ctx)

		l.Debug("getting user session")
		userSession, err := getSession(w, req, h.sessionStore, sessionKey)
		if err != nil {
			l.Error("failed to get user session", zap.Error(err))
			clearSession(w, req, h.sessionStore, sessionKey)
			utils.ErrorResponse(w, meta.NewInternalServerError(errors.New("failed to get user session")))
		}

		l.Debug("getting access token from session")
		accessToken, ok := strFromSession(userSession, constants.SessionAccessToken)
		if !ok {
			l.Debug("access token not in session")
			utils.JSONResponse(w, http.StatusOK, &ExposedSession{Active: false})
		}

		l.Debug("introspecting access token")
		introspection, err := hydraAdmin.IntrospectOAuth2Token(&admin.IntrospectOAuth2TokenParams{
			Token:   accessToken,
			Context: ctx,
		})
		if err != nil {
			l.Error("failed to introspect access token", zap.Error(err))
			utils.ErrorResponse(w, meta.NewInternalServerError(err))
			return
		}

		exp := time.Unix(introspection.Payload.Exp, 0).UTC()
		expIn := int(exp.Sub(time.Now()).Seconds())
		if exp.Unix() == 0 {
			expIn = 0
		}

		s := ExposedSession{
			Active:           *introspection.Payload.Active,
			Expiry:           exp,
			ExpiresInSeconds: expIn,
			Subject:          introspection.Payload.Sub,
			Username:         introspection.Payload.Username,
		}

		utils.JSONResponse(w, http.StatusOK, &s)
	}
}

func (h *Handler) RestfulSession(sessionKey string, hydraAdmin admin.ClientService) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		h.Session(sessionKey, hydraAdmin)(res.ResponseWriter, req.Request)
	}
}
