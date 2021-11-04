package login

import (
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/nrc-no/core/pkg/store"
	"github.com/ory/hydra-client-go/models"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"net/http"
)

func handleUseIdentityProvider(
	w http.ResponseWriter,
	req *http.Request,
	getLoginRequest func(loginChallenge string) (*models.LoginRequest, error),
	pathParameters map[string]string,
	idpStore store.IdentityProviderStore,
	selfURL string,
) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {

	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {

		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("event", authrequest.EventUseIdentityProvider))

		l.Debug("getting login request")
		loginRequest, err := getLoginRequest(authRequest.LoginChallenge)
		if err != nil {
			l.Error("failed to get login request", zap.Error(err))
			return err
		}

		identityProviderID := pathParameters["identityProviderId"]
		l.Debug("getting identity provider", zap.String("identity_provider_id", identityProviderID))
		idp, err := idpStore.Get(ctx, identityProviderID, store.IdentityProviderGetOptions{ReturnClientSecret: true})
		if err != nil {
			l.Error("failed to get identity provider", zap.Error(err))
			return err
		}
		authRequest.IdentityProviderId = identityProviderID

		l.Debug("creating state variable")
		stateVar, err := createStateVariable()
		if err != nil {
			l.Error("failed to create state variable", zap.Error(err))
			return err
		}
		authRequest.StateVariable = stateVar

		l.Debug("getting identity provider oauth2 config")
		oauth2Config, _, _, err := getOauthProvider(ctx, idp, selfURL, loginRequest)
		if err != nil {
			l.Error("failed to get identity provider oauth2 config", zap.Error(err))
			return err
		}

		authCodeURL := oauth2Config.AuthCodeURL(stateVar, oauth2.SetAuthURLParam("login_hint", authRequest.Identifier))
		l.Debug("redirecting to oauth2 login")
		http.Redirect(w, req, authCodeURL, http.StatusSeeOther)

		return nil
	}
}
