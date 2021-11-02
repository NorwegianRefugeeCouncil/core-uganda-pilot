package authrequest

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthRequestPasswordFlow(t *testing.T) {
	a := NewAuthRequest()

	// start
	assert.Equal(t, StateStart, a.currentState())

	// request login
	assert.NoError(t, a.RequestLogin(RequestLoginProps{}))
	assert.Equal(t, StateLoginRequested, a.currentState())

	// perform login
	assert.NoError(t, a.PerformLogin())
	assert.Equal(t, StatePromptingForIdentifier, a.currentState())

	// provide identifier
	assert.NoError(t, a.ProvideIdentifier("identifier"))
	assert.Equal(t, StateValidatingIdentifier, a.currentState())

	// provide invalid identifier
	assert.NoError(t, a.FailIdentifierValidation(errors.New("invalid")))
	assert.Equal(t, StatePromptingForIdentifier, a.currentState())

	// provide identifier
	assert.NoError(t, a.ProvideIdentifier("identifier"))
	assert.Equal(t, StateValidatingIdentifier, a.currentState())

	// provide valid identifier
	assert.NoError(t, a.SucceedIdentifierValidation())
	assert.Equal(t, StateFindingAuthMethod, a.currentState())

	// use password auth
	assert.NoError(t, a.UsePasswordAuth())
	assert.Equal(t, StatePromptingForPassword, a.currentState())

	// provide password
	assert.NoError(t, a.ProvidePassword("password"))
	assert.Equal(t, StateValidatingPassword, a.currentState())

	// provide invalid password
	assert.NoError(t, a.FailPasswordValidation())
	assert.Equal(t, StatePromptingForPassword, a.currentState())

	// provide password
	assert.NoError(t, a.ProvidePassword("password"))
	assert.Equal(t, StateValidatingPassword, a.currentState())

	// provide valid password
	assert.NoError(t, a.SucceedPasswordValidation())
	assert.Equal(t, StatePasswordLoginSucceeded, a.currentState())

	// accept login request
	assert.NoError(t, a.AcceptLoginRequest())
	assert.Equal(t, StateAwaitingConsentChallenge, a.currentState())

	// present consent
	assert.NoError(t, a.PresentConsentChallenge())
	assert.Equal(t, StatePresentingConsent, a.currentState())

	// approve consent
	assert.NoError(t, a.ApproveConsentRequest())
	assert.Equal(t, StateAccepted, a.currentState())

}

func TestAuthRequestOidcFlow(t *testing.T) {
	a := NewAuthRequest()

	// start
	assert.Equal(t, StateStart, a.currentState())

	// request login
	assert.NoError(t, a.RequestLogin(RequestLoginProps{}))
	assert.Equal(t, StateLoginRequested, a.currentState())

	// perform login
	assert.NoError(t, a.PerformLogin())
	assert.Equal(t, StatePromptingForIdentifier, a.currentState())

	// provide identifier
	assert.NoError(t, a.ProvideIdentifier("identifier"))
	assert.Equal(t, StateValidatingIdentifier, a.currentState())

	// identifier valid
	assert.NoError(t, a.SucceedIdentifierValidation())
	assert.Equal(t, StateFindingAuthMethod, a.currentState())

	// use oidc auth
	assert.NoError(t, a.UseOidcAuth())
	assert.Equal(t, StatePromptingForIdentityProvider, a.currentState())

	// select identity provider
	assert.NoError(t, a.UseIdentityProvider("identityprovider1"))
	assert.Equal(t, StateAwaitingIdpCallback, a.currentState())

	// callback called
	assert.NoError(t, a.CallIdentityProviderCallback())
	assert.Equal(t, StatePerformingAuthCodeExchange, a.currentState())

	// succeed callback
	assert.NoError(t, a.SucceedAuthCodeExchange())
	assert.Equal(t, StateAuthCodeExchangeSucceeded, a.currentState())

	// accept login request
	assert.NoError(t, a.AcceptLoginRequest())
	assert.Equal(t, StateAwaitingConsentChallenge, a.currentState())

	// present consent
	assert.NoError(t, a.PresentConsentChallenge())
	assert.Equal(t, StatePresentingConsent, a.currentState())

	// approve consent
	assert.NoError(t, a.ApproveConsentRequest())
	assert.Equal(t, StateAccepted, a.currentState())

}
