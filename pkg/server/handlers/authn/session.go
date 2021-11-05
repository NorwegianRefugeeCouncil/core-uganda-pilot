package authn

import (
	"errors"
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
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

func (h *Handler) Session(sessionKey string) http.HandlerFunc {
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

		l.Debug("getting token from session")
		tok, ok := tokenFromSession(userSession)
		if !ok {
			l.Debug("access token not in session")
			utils.JSONResponse(w, http.StatusOK, &ExposedSession{Active: false})
			return
		}

		idTok, err := h.tokenVerifier.Verify(ctx, tok.IDToken)
		if err != nil {
			l.Debug("token is invalid")
			utils.JSONResponse(w, http.StatusOK, &ExposedSession{Active: false})
			return
		}

		expIn := int(idTok.Expiry.Sub(time.Now()).Seconds())
		if tok.Expiry.Unix() == 0 {
			expIn = 0
		}

		s := ExposedSession{
			Active:           true,
			Expiry:           tok.Expiry,
			ExpiresInSeconds: expIn,
			Subject:          idTok.Subject,
		}

		utils.JSONResponse(w, http.StatusOK, &s)
	}
}

func (h *Handler) RestfulSession(sessionKey string) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		h.Session(sessionKey)(res.ResponseWriter, req.Request)
	}
}
