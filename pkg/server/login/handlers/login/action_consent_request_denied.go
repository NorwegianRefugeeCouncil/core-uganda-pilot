package login

import (
	"context"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"go.uber.org/zap"
)

func handleConsentRequestDenied(ctx context.Context, dispatch func(evt string), hydraAdmin admin.ClientService) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StateConsentRequestDeclined))
		l.Debug("rejecting consent request")
		resp, err := hydraAdmin.RejectConsentRequest(&admin.RejectConsentRequestParams{
			Context:          ctx,
			ConsentChallenge: authRequest.ConsentChallenge,
			Body:             &models.RejectRequest{},
		})
		if err != nil {
			l.Error("failed to reject consent request", zap.Error(err))
			return err
		}
		authRequest.PostConsentURL = *resp.Payload.RedirectTo
		l.Debug("declining auth request")
		dispatch(authrequest.EventDecline)
		return nil
	}
}
