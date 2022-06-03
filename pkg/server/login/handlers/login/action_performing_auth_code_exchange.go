package login

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/nrc-no/core/pkg/store"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"net/http"
	"strconv"
	"text/template"
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
		processedToken, err := processOidcToken(req.Context(), tokenFromExchange, verifier, idp)
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
	idp *types.IdentityProvider,
) (*ProcessedToken, error) {

	l := logging.NewLogger(ctx)

	// getting id token from exchange
	rawIDTokenIntf := token.Extra("id_token")
	if rawIDTokenIntf == nil {
		l.Warn("id token not present in token")
		var err = errors.New("id token not present in token exchange response")
		return nil, err
	}

	// converting id token to string
	rawIDToken, ok := rawIDTokenIntf.(string)
	if !ok {
		l.Warn(fmt.Sprintf("id token in response was not a string but was: %T", rawIDTokenIntf))
		var err = errors.New("id token in exchange response was not a string")
		return nil, err
	}

	// verifying id token
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		l.Warn("failed to verify ID token", zap.Error(err))
		return nil, err
	}


	var claimInft interface{}
	if err := idToken.Claims(&claimInft); err != nil {
		l.Error("failed to unmarshal claims from ID token", zap.Error(err))
		return nil, err
	}

	userProfile, err := extractIdentityProfileTemplateVersion(ctx, claimInft, idp)

	result := &ProcessedToken{
		OriginalToken: token,
		RawIDToken:    rawIDToken,
		AccessToken:   token.AccessToken,
		RefreshToken:  token.RefreshToken,
		IDToken:       idToken,
		Claims:        userProfile,
	}

	return result, nil
}


func extractIdentityProfileTemplateVersion(
	ctx context.Context,
	claimInft interface{},
	idp *types.IdentityProvider,
) (*authrequest.Claims, error){

	l := logging.NewLogger(ctx)

	var userProfile authrequest.Claims

	mappedClaimInft, ok := claimInft.(map[string]interface{})
	if !ok {
		l.Error("failed to cast claims into type map[string]interface{}")
		return nil, errors.New("failed to cast claims into type map[string]interface{}")
	}

	subjectTpl := template.New("")
	subjectTplParsed, err := subjectTpl.Parse(idp.ClaimMappings.Subject)
	if err != nil { return nil, err}
	var subjectBuffer bytes.Buffer
	if err := subjectTplParsed.Execute(&subjectBuffer, mappedClaimInft); err != nil {
		l.Error("failed to execute template for subject claim", zap.Error(err))
		return nil, err
	}
	subjectString := subjectBuffer.String()
	if err != nil { return nil, err}
	userProfile.Subject = subjectString

	displayNameTpl := template.New("")
	displayNameTplParsed, err := displayNameTpl.Parse(idp.ClaimMappings.DisplayName)
	if err != nil { return nil, err}
	var displayNameBuffer bytes.Buffer
	if err := displayNameTplParsed.Execute(&displayNameBuffer, mappedClaimInft); err != nil {
		l.Error("failed to execute template for displayName claim", zap.Error(err))
		return nil, err
	}
	displayNameString := displayNameBuffer.String()
	if err != nil { return nil, err}
	userProfile.DisplayName = displayNameString

	fullNameTpl := template.New("")
	fullNameTplParsed, err := fullNameTpl.Parse(idp.ClaimMappings.FullName)
	if err != nil { return nil, err}
	var fullNameBuffer bytes.Buffer
	if err := fullNameTplParsed.Execute(&fullNameBuffer, mappedClaimInft); err != nil {
		l.Error("failed to execute template for FullName claim", zap.Error(err))
		return nil, err
	}
	fullNameString := fullNameBuffer.String()
	if err != nil { return nil, err}
	userProfile.FullName = fullNameString

	emailTpl := template.New("")
	emailTplParsed, err := emailTpl.Parse(idp.ClaimMappings.Email)
	if err != nil { return nil, err}
	var emailBuffer bytes.Buffer
	if err := emailTplParsed.Execute(&emailBuffer, mappedClaimInft); err != nil {
		l.Error("failed to execute template for FullName claim", zap.Error(err))
		return nil, err
	}
	emailString := emailBuffer.String()
	if err != nil { return nil, err}
	userProfile.Email = emailString

	emailVerifiedTpl := template.New("")
	emailVerifiedTplParsed, err := emailVerifiedTpl.Parse(idp.ClaimMappings.EmailVerified)
	if err != nil { return nil, err}
	var emailVerifiedBuffer bytes.Buffer
	if err := emailVerifiedTplParsed.Execute(&emailVerifiedBuffer, mappedClaimInft); err != nil {
		l.Error("failed to execute template for FullName claim", zap.Error(err))
		return nil, err
	}
	emailVerifiedString := emailVerifiedBuffer.String()
	emailVerifiedBool, _ := strconv.ParseBool(emailVerifiedString)
	if err != nil { return nil, err}
	userProfile.EmailVerified = emailVerifiedBool


	return &userProfile, nil
}
