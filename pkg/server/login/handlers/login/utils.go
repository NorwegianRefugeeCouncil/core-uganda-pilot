package login

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/ory/hydra-client-go/models"
	"golang.org/x/oauth2"
	"strings"
)

func getEmailDomain(val string) (string, error) {
	components := strings.Split(val, "@")
	if len(components) != 2 {
		return "", errors.New("invalid email address")
	}
	if len(components[0]) == 0 {
		return "", errors.New("invalid email address")
	}
	return components[1], nil
}

func createStateVariable() (string, error) {
	bts, err := generateBytes(32)
	if err != nil {
		return "", err
	}
	state := base64.StdEncoding.EncodeToString(bts)
	return state, nil
}

func generateBytes(count int) ([]byte, error) {
	b := make([]byte, count)
	_, err := rand.Read(b)
	if err != nil {
		return nil, meta.NewInternalServerError(err)
	}
	return b, nil
}

func getOauthProvider(
	ctx context.Context,
	client *types.IdentityProvider,
	selfURL string,
	loginRequest *models.LoginRequest,
) (*oauth2.Config, *oidc.Provider, *oidc.IDTokenVerifier, error) {

	// getting oidc provider
	oidcProvider, err := oidc.NewProvider(ctx, client.Domain)
	if err != nil {
		return nil, nil, nil, err
	}

	// getting oauth2 config
	redirectUri := fmt.Sprintf("%s/login/oidc/callback", selfURL)
	oauthConfig := &oauth2.Config{
		ClientID:     client.ClientID,
		ClientSecret: client.ClientSecret,
		Endpoint:     oidcProvider.Endpoint(),
		RedirectURL:  redirectUri,
	}
	if loginRequest != nil {
		oauthConfig.Scopes = loginRequest.RequestedScope
	}

	verifier := oidcProvider.Verifier(&oidc.Config{
		ClientID: client.ClientID,
	})

	return oauthConfig, oidcProvider, verifier, nil
}
