package login

import (
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/nrc-no/core/pkg/server/login/templates"
	"go.uber.org/zap"
	"net/http"
)

func handlePromptingForIdentifier(w http.ResponseWriter, req *http.Request) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StatePromptingForIdentifier))

		l.Debug("prompting user for identifier")
		err := templates.Template.ExecuteTemplate(w, "login_subject", map[string]interface{}{
			"Error": authRequest.IdentifierError,
		})
		if err != nil {
			l.Error("failed to prompt user for identifier", zap.Error(err))
			return err
		}
		return nil
	}
}
