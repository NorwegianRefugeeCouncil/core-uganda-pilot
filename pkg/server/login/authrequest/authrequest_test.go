package authrequest

import (
	"errors"
	"github.com/nrc-no/core/pkg/utils/sets"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthRequestPasswordFlow(t *testing.T) {
	a := NewAuthRequest(Handlers{})

	// start
	assert.Equal(t, StateStart, a.currentState())

	// request login
	assert.NoError(t, a.RequestLogin())
	assert.Equal(t, StateLoginRequested, a.currentState())

	// perform login
	assert.NoError(t, a.PerformLogin())
	assert.Equal(t, StatePromptingForIdentifier, a.currentState())

	// provide identifier
	assert.NoError(t, a.ProvideIdentifier())
	assert.Equal(t, StateValidatingIdentifier, a.currentState())

	// provide invalid identifier
	assert.NoError(t, a.FailIdentifierValidation(errors.New("invalid")))
	assert.Equal(t, StatePromptingForIdentifier, a.currentState())

	// provide identifier
	assert.NoError(t, a.ProvideIdentifier())
	assert.Equal(t, StateValidatingIdentifier, a.currentState())

	// provide valid identifier
	assert.NoError(t, a.SucceedIdentifierValidation())
	assert.Equal(t, StateFindingAuthMethod, a.currentState())

	// use password auth
	assert.NoError(t, a.UsePasswordAuth())
	assert.Equal(t, StatePromptingForPassword, a.currentState())

	// provide password
	assert.NoError(t, a.ProvidePassword())
	assert.Equal(t, StateValidatingPassword, a.currentState())

	// provide invalid password
	assert.NoError(t, a.FailPasswordValidation())
	assert.Equal(t, StatePromptingForPassword, a.currentState())

	// provide password
	assert.NoError(t, a.ProvidePassword())
	assert.Equal(t, StateValidatingPassword, a.currentState())

	// provide valid password
	assert.NoError(t, a.SucceedPasswordValidation())
	assert.Equal(t, StatePasswordLoginSucceeded, a.currentState())

	// accept login request
	assert.NoError(t, a.AcceptLoginRequest())
	assert.Equal(t, StateAwaitingConsentChallenge, a.currentState())

	// receive consent
	assert.NoError(t, a.ReceiveConsentChallenge())
	assert.Equal(t, StateReceivedConsentChallenge, a.currentState())

	// present consent
	assert.NoError(t, a.PresentConsentChallenge())
	assert.Equal(t, StatePresentingConsent, a.currentState())

	// approve consent
	assert.NoError(t, a.ApproveConsentRequest())
	assert.Equal(t, StateConsentRequestApproved, a.currentState())

	// approve
	assert.NoError(t, a.Accept())
	assert.Equal(t, StateAccepted, a.currentState())

}

func TestAuthRequestOidcFlow(t *testing.T) {
	a := NewAuthRequest(Handlers{})

	// start
	assert.Equal(t, StateStart, a.currentState())

	// request login
	assert.NoError(t, a.RequestLogin())
	assert.Equal(t, StateLoginRequested, a.currentState())

	// perform login
	assert.NoError(t, a.PerformLogin())
	assert.Equal(t, StatePromptingForIdentifier, a.currentState())

	// provide identifier
	assert.NoError(t, a.ProvideIdentifier())
	assert.Equal(t, StateValidatingIdentifier, a.currentState())

	// identifier valid
	assert.NoError(t, a.SucceedIdentifierValidation())
	assert.Equal(t, StateFindingAuthMethod, a.currentState())

	// use oidc auth
	assert.NoError(t, a.UseOidcAuth())
	assert.Equal(t, StatePromptingForIdentityProvider, a.currentState())

	// select identity provider
	assert.NoError(t, a.UseIdentityProvider())
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

	// receive consent
	assert.NoError(t, a.ReceiveConsentChallenge())
	assert.Equal(t, StateReceivedConsentChallenge, a.currentState())

	// present consent
	assert.NoError(t, a.PresentConsentChallenge())
	assert.Equal(t, StatePresentingConsent, a.currentState())

	// approve consent
	assert.NoError(t, a.ApproveConsentRequest())
	assert.Equal(t, StateConsentRequestApproved, a.currentState())

	// approve consent
	assert.NoError(t, a.Accept())
	assert.Equal(t, StateAccepted, a.currentState())

}

func TestTransitions(t *testing.T) {
	a := NewAuthRequest(Handlers{})
	assertOnlyEventAccepted(t, a, StateStart, sets.NewString(EventRequestLogin, EventFail))
	assertOnlyEventAccepted(t, a, StateLoginRequested, sets.NewString(EventSkipLoginRequest, EventPerformLogin, EventFail))
	assertOnlyEventAccepted(t, a, StatePromptingForIdentifier, sets.NewString(EventProvideIdentifier, EventFail))
	assertOnlyEventAccepted(t, a, StateValidatingIdentifier, sets.NewString(EventProvideValidIdentifier, EventProvideInvalidIdentifier, EventFail))
	assertOnlyEventAccepted(t, a, StateFindingAuthMethod, sets.NewString(EventUseOidcLogin, EventUsePasswordLogin, EventFail))
	assertOnlyEventAccepted(t, a, StatePromptingForPassword, sets.NewString(EventProvidePassword, EventFail))
	assertOnlyEventAccepted(t, a, StateValidatingPassword, sets.NewString(EventProvideInvalidPassword, EventProvideValidPassword, EventFail))
	assertOnlyEventAccepted(t, a, StatePasswordLoginSucceeded, sets.NewString(EventAcceptLoginRequest, EventFail))
	assertOnlyEventAccepted(t, a, StatePromptingForIdentityProvider, sets.NewString(EventUseIdentityProvider, EventFail))
	assertOnlyEventAccepted(t, a, StateAwaitingIdpCallback, sets.NewString(EventCallOidcCallback, EventFail))
	assertOnlyEventAccepted(t, a, StatePerformingAuthCodeExchange, sets.NewString(EventSucceedCodeExchange, EventFailCodeExchange, EventFail))
	assertOnlyEventAccepted(t, a, StateAuthCodeExchangeSucceeded, sets.NewString(EventAcceptLoginRequest, EventFail))
	assertOnlyEventAccepted(t, a, StateAwaitingConsentChallenge, sets.NewString(EventReceiveConsentChallenge, EventFail))
	assertOnlyEventAccepted(t, a, StateReceivedConsentChallenge, sets.NewString(EventPresentConsentChallenge, EventSkipConsentChallenge, EventFail))
	assertOnlyEventAccepted(t, a, StatePresentingConsent, sets.NewString(EventApproveConsentChallenge, EventDeclineConsentChallenge, EventFail))
	assertOnlyEventAccepted(t, a, StateConsentRequestApproved, sets.NewString(EventAccept, EventFail))
	assertOnlyEventAccepted(t, a, StateConsentRequestDeclined, sets.NewString(EventDecline, EventFail))
	assertOnlyEventAccepted(t, a, StateAccepted, sets.NewString(EventFail, EventRequestLogin))
	assertOnlyEventAccepted(t, a, StateDeclined, sets.NewString(EventFail))
	assertOnlyEventAccepted(t, a, StateFailed, sets.NewString())
	assertOnlyEventAccepted(t, a, StateRefreshingIdentity, sets.NewString(EventSetRefreshedIdentity, EventFail))
}

func assertOnlyEventAccepted(t *testing.T, a *AuthRequest, fromState string, accepted sets.String) {
	allEvts := sets.NewString(allEvents...)
	notAccepted := allEvts.Difference(accepted)
	for evt := range notAccepted {
		a.fsm.SetState(fromState)
		assert.Error(t, a.fsm.Event(evt), "applying event %s to state %s should produce error but did not", evt, fromState)
	}
	for evt := range accepted {
		a.fsm.SetState(fromState)
		assert.NoError(t, a.fsm.Event(evt), "applying event %s to state %s should not produce error but did", evt, fromState)
	}
}
