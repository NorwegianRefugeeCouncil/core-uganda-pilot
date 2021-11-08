package login

import (
	"context"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"go.uber.org/zap"
)

func handleFindingAuthMethod(ctx context.Context, dispatch func(evt string)) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StateFindingAuthMethod))
		l.Debug("using oidc auth")
		dispatch(authrequest.EventUseOidcLogin)
		return nil
	}
}
