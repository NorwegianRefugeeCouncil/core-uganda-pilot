package login

import (
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/nrc-no/core/pkg/server/login/templates"
	"go.uber.org/zap"
	"net/http"
)

func handlePromptingForIdentifier(w http.ResponseWriter, req *http.Request, userSession *sessions.Session, enqueue func(fn func())) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StatePromptingForIdentifier))
		l.Debug("entered state")

		l.Debug("saving auth request")
		if err := authRequest.Save(w, req, userSession); err != nil {
			l.Error("failed to save auth request", zap.Error(err))
			enqueue(func() {
				_ = authRequest.Fail(err)
			})
			return
		}

		l.Debug("prompting user for identifier")
		err := templates.Template.ExecuteTemplate(w, "login_subject", map[string]interface{}{
			"Error": authRequest.IdentifierError,
		})
		if err != nil {
			l.Error("failed to prompt user for identifier", zap.Error(err))
		}
		return
	}
}
