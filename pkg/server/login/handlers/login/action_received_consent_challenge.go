package login

import (
	"context"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/ory/hydra-client-go/models"
	"go.uber.org/zap"
	"net/url"
)

func handleReceivedConsentChallenge(
	ctx context.Context,
	dispatch func(evt string),
	requestParameters url.Values,
	getConsentRequest func(consentChallenge string) (*models.ConsentRequest, error),
) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {

	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StateReceivedConsentChallenge))

		l.Debug("getting consent request")
		consentChallenge := requestParameters.Get("consent_challenge")
		consentRequest, err := getConsentRequest(consentChallenge)
		if err != nil {
			l.Error("failed to get consent request", zap.Error(err))
			return err
		}

		authRequest.ConsentChallenge = consentChallenge

		if consentRequest.Skip {
			l.Debug("skipping consent request")
			dispatch(authrequest.EventSkipConsentChallenge)
		} else {
			l.Debug("presenting consent challenge")
			dispatch(authrequest.EventPresentConsentChallenge)
		}
		return nil
	}
}
