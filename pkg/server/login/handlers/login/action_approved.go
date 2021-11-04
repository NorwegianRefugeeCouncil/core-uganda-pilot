package login

import (
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"go.uber.org/zap"
	"net/http"
)

func handleApproved(
	w http.ResponseWriter,
	req *http.Request,
) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StateAccepted))
		l.Debug("redirecting to post-consent url")
		http.Redirect(w, req, authRequest.PostConsentURL, http.StatusTemporaryRedirect)
		return nil
	}
}
