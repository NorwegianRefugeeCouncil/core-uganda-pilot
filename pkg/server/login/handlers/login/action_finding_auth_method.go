package login

import (
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"go.uber.org/zap"
	"net/http"
)

func handleFindingAuthMethod(w http.ResponseWriter, req *http.Request, userSession *sessions.Session, enqueue func(fn func())) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StateFindingAuthMethod))
		l.Debug("entered state")

		l.Debug("saving auth request")
		if err := authRequest.Save(w, req, userSession); err != nil {
			l.Error("failed to save auth request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		l.Debug("using oidc auth")
		enqueue(func() {
			if err := authRequest.UseOidcAuth(); err != nil {
				l.Error("failed to use oidc auth", zap.Error(err))
				enqueue(func() {
					_ = authRequest.Fail(err)
				})
				return
			}
		})

	}
}
