package login

import (
	"context"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"go.uber.org/zap"
	"net/http"
)

func handleValidatingIdentifier(ctx context.Context, req *http.Request, dispatch func(evt string)) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {

		l := logging.NewLogger(ctx).With(zap.String("event", authrequest.EventUseIdentityProvider))

		l.Debug("parsing form data")
		if err := req.ParseForm(); err != nil {
			l.Error("failed to parse form data", zap.Error(err))
			dispatch(authrequest.EventProvideInvalidIdentifier)
			return nil
		}

		q := req.Form

		l.Debug("getting email domain")
		email := q.Get("email")
		emailDomain, err := getEmailDomain(email)
		if err != nil {
			l.Error("failed to get email domain", zap.Error(err))
			dispatch(authrequest.EventProvideInvalidIdentifier)
			return nil
		}

		authRequest.Identifier = email
		authRequest.EmailDomain = emailDomain

		dispatch(authrequest.EventProvideValidIdentifier)
		return nil
	}
}
