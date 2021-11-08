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

func handleConsentRequestApproved(
	ctx context.Context,
	dispatch func(evt string),
	getConsentRequest func(consentChallenge string) (*models.ConsentRequest, error),
	hydraAdmin admin.ClientService,
) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {

	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {

		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StateConsentRequestApproved))

		l.Debug("getting consent request")
		consentRequest, err := getConsentRequest(authRequest.ConsentChallenge)
		if err != nil {
			l.Error("failed to get consent request", zap.Error(err))
			return err
		}

		l.Debug("accepting consent request")
		resp, err := hydraAdmin.AcceptConsentRequest(&admin.AcceptConsentRequestParams{
			Context:          ctx,
			ConsentChallenge: authRequest.ConsentChallenge,
			Body: &models.AcceptConsentRequest{

				GrantAccessTokenAudience: consentRequest.RequestedAccessTokenAudience,
				GrantScope:               consentRequest.RequestedScope,
				Remember:                 true,
				RememberFor:              0,
			},
		})
		if err != nil {
			l.Error("failed to accept consent request", zap.Error(err))
			return err
		}

		authRequest.PostConsentURL = *resp.Payload.RedirectTo
		dispatch(authrequest.EventAccept)
		return nil

	}
}
