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

func handleLoginRequested(
	ctx context.Context,
	queryParameters url.Values,
	dispatch func(evt string),
	getLoginRequest func(loginChallenge string) (*models.LoginRequest, error),
) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {

		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StateLoginRequested))

		l.Debug("getting login challenge")
		loginChallenge := queryParameters.Get("login_challenge")
		loginRequest, err := getLoginRequest(loginChallenge)
		if err != nil {
			l.Error("failed to get login challenge", zap.Error(err))
			return err
		}

		authRequest.LoginChallenge = loginChallenge

		l.Debug("checking if should skip login request or not")
		if loginRequest.Skip != nil && *loginRequest.Skip {
			l.Debug("skipping login request")
			authRequest.Subject = *loginRequest.Subject
			dispatch(authrequest.EventSkipLoginRequest)
		} else {
			l.Debug("performing login request")
			dispatch(authrequest.EventPerformLogin)
		}

		return nil
	}
}
