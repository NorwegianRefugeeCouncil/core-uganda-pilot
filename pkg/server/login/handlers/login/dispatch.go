package login

import (
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/sirupsen/logrus"
	"net/http"
)

func dispatchAction(w http.ResponseWriter, action string, enqueue func(fn func()), authRequest *authrequest.AuthRequest) {

	logrus.Infof("dispatching login action: %s", action)

	switch action {
	case authrequest.EventRequestLogin:
		enqueue(func() {
			if err := authRequest.RequestLogin(); err != nil {
				handleError(w, http.StatusBadRequest, err)
			}
		})

	case authrequest.EventProvideIdentifier:
		enqueue(func() {
			if err := authRequest.ProvideIdentifier(); err != nil {
				handleError(w, http.StatusBadRequest, err)
			}
		})

	case authrequest.EventUseIdentityProvider:
		enqueue(func() {
			if err := authRequest.UseIdentityProvider(); err != nil {
				handleError(w, http.StatusBadRequest, err)
			}
		})

	case authrequest.EventCallOidcCallback:
		enqueue(func() {
			if err := authRequest.CallIdentityProviderCallback(); err != nil {
				handleError(w, http.StatusBadRequest, err)
			}
		})

	case authrequest.EventReceiveConsentChallenge:
		enqueue(func() {
			if err := authRequest.ReceiveConsentChallenge(); err != nil {
				handleError(w, http.StatusBadRequest, err)
			}
		})

	case authrequest.EventApproveConsentChallenge:
		enqueue(func() {
			if err := authRequest.ApproveConsentRequest(); err != nil {
				handleError(w, http.StatusBadRequest, err)
			}
		})

	case authrequest.EventDeclineConsentChallenge:
		enqueue(func() {
			if err := authRequest.DeclineConsentRequest(); err != nil {
				handleError(w, http.StatusBadRequest, err)
			}
		})

	case authrequest.EventProvidePassword:
		enqueue(func() {
			if err := authRequest.ProvidePassword(); err != nil {
				handleError(w, http.StatusBadRequest, err)
			}
		})
	}
}
