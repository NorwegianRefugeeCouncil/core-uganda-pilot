package authrequest

import (
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/utils/sets"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransitions(t *testing.T) {
	a := NewAuthRequest(map[string]fsm.Callback{})
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
