package login

import (
	"context"
	"errors"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/nrc-no/core/pkg/store"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"net/http"
)

func handlePerformingAuthCodeExchange(
	req *http.Request,
	dispatch func(evt string),
	idpStore store.IdentityProviderStore,
	selfURL string,
) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {

	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {

		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StatePerformingAuthCodeExchange))

		l.Debug("getting identity provider", zap.String("identity_provider_id", authRequest.IdentityProviderId))
		idp, err := idpStore.Get(ctx, authRequest.IdentityProviderId, store.IdentityProviderGetOptions{ReturnClientSecret: true})
		if err != nil {
			l.Error("failed to get identity provider", zap.Error(err))
			return err
		}

		l.Debug("getting identity provider oauth configuration")
		oauth2Config, _, verifier, err := getOauthProvider(ctx, idp, selfURL, nil)
		if err != nil {
			l.Error("failed to get identity provider oauth configuration", zap.Error(err))
			return err
		}

		l.Debug("getting state from query")
		stateFromQuery := req.URL.Query().Get("state")
		if len(stateFromQuery) == 0 {
			l.Error("state not found in query parameters")
			err := errors.New("state not found in response")
			return err
		}

		l.Debug("getting authorization code from query")
		authorizationCodeFromQuery := req.URL.Query().Get("code")
		if len(authorizationCodeFromQuery) == 0 {
			l.Error("authorization code not found in query")
			err := errors.New("code not found in response")
			return err
		}

		l.Debug("exchanging authorization code")
		tokenFromExchange, err := oauth2Config.Exchange(req.Context(), authorizationCodeFromQuery)
		if err != nil {
			l.Error("failed to exchange authorization code", zap.Error(err))
			return err
		}

		l.Debug("verifying token")
		processedToken, err := processOidcToken(req.Context(), tokenFromExchange, verifier)
		if err != nil {
			l.Error("failed to verify token", zap.Error(err))
			return err
		}

		l.Debug("saving token in auth request")
		saveTokenInAuthRequest(authRequest, processedToken)

		l.Debug("succeeding auth code exchange")
		dispatch(authrequest.EventSucceedCodeExchange)
		return nil

	}
}

func saveTokenInAuthRequest(authRequest *authrequest.AuthRequest, processedToken *ProcessedToken) {
	authRequest.Claims = processedToken.Claims
	if len(authRequest.IDToken) != 0 && len(processedToken.RawIDToken) != 0 {
		authRequest.IDToken = processedToken.RawIDToken
	}
	authRequest.AccessToken = processedToken.AccessToken
	if len(processedToken.RefreshToken) > 0 {
		authRequest.RefreshToken = processedToken.RefreshToken
	}
	authRequest.TokenType = processedToken.OriginalToken.TokenType
	authRequest.TokenExpiry = processedToken.OriginalToken.Expiry
}

type ProcessedToken struct {
	OriginalToken *oauth2.Token
	RawIDToken    string
	AccessToken   string
	RefreshToken  string
	IDToken       *oidc.IDToken
	Claims        *authrequest.Claims
}

func processOidcToken(
	ctx context.Context,
	token *oauth2.Token,
	verifier *oidc.IDTokenVerifier,
) (*ProcessedToken, error) {

	// getting id token from exchange
	rawIDTokenIntf := token.Extra("id_token")
	if rawIDTokenIntf == nil {
		logrus.Warnf("id token not present in token")
		var err = errors.New("id token not present in token exchange response")
		return nil, err
	}

	// converting id token to string
	rawIDToken, ok := rawIDTokenIntf.(string)
	if !ok {
		logrus.Warnf("id token in response was not a string but was: %T", rawIDTokenIntf)
		var err = errors.New("id token in exchange response was not a string")
		return nil, err
	}

	// verifying id token
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		logrus.Warnf("failed to verify ID token: %v", err)
		return nil, err
	}

	// unmarshal claims
	var userProfile authrequest.Claims
	if err := idToken.Claims(&userProfile); err != nil {
		logrus.WithError(err).Warnf("failed to unmarshal claims from ID token")
		return nil, err
	}

	result := &ProcessedToken{
		OriginalToken: token,
		RawIDToken:    rawIDToken,
		AccessToken:   token.AccessToken,
		RefreshToken:  token.RefreshToken,
		IDToken:       idToken,
		Claims:        &userProfile,
	}

	return result, nil
}
