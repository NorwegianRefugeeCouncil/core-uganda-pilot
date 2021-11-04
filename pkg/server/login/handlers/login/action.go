package login

import (
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	loginstore "github.com/nrc-no/core/pkg/server/login/store"
	"github.com/nrc-no/core/pkg/store"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"sync"
)

func handleError(w http.ResponseWriter, status int, err error) {
	logrus.WithError(err).Errorf("login error")
	w.Write([]byte(err.Error()))
	w.WriteHeader(status)
}

func handleAuthRequestAction(
	sessionStore sessions.Store,
	idpStore store.IdentityProviderStore,
	orgStore store.OrganizationStore,
	loginStore loginstore.Interface,
	hydraAdmin admin.ClientService,
	selfURL string,
) func(action string, pathParameters map[string]string, requestParameters url.Values) http.HandlerFunc {

	return func(action string, pathParameters map[string]string, requestParameters url.Values) http.HandlerFunc {

		return func(w http.ResponseWriter, req *http.Request) {

			ctx := req.Context()

			logger := logrus.
				WithField("server", "core-login").
				WithField("name", "login_action")

			logger.
				WithField("action", action).
				WithField("path_parameters", pathParameters).
				WithField("request_parameters", requestParameters).
				Tracef("received login action")

			logger.Tracef("getting user session")

			// getting user session
			userSession, done := getUserSession(w, req, sessionStore)
			if done {
				return
			}

			// cache login request
			var _loginRequest *models.LoginRequest
			getLoginRequest := func(loginChallenge string) (*models.LoginRequest, error) {

				if _loginRequest != nil {
					logger.Tracef("returning cached login request")
					return _loginRequest, nil
				}

				logger.Tracef("getting login request")
				loginRequestResp, err := hydraAdmin.GetLoginRequest(&admin.GetLoginRequestParams{
					Context:        ctx,
					LoginChallenge: loginChallenge,
				})
				if err != nil {
					logger.WithError(err).Errorf("failed to get login request")
					return nil, err
				}
				_loginRequest = loginRequestResp.Payload
				return _loginRequest, nil
			}

			// cache consent request
			var _consentRequest *models.ConsentRequest
			getConsentRequest := func(consentChallenge string) (*models.ConsentRequest, error) {
				if _consentRequest != nil {
					logger.Tracef("returning cached consent request")
					return _consentRequest, nil
				}
				logger.Tracef("getting consent request")
				consentRequestResp, err := hydraAdmin.GetConsentRequest(&admin.GetConsentRequestParams{
					Context:          ctx,
					ConsentChallenge: consentChallenge,
				})
				if err != nil {
					logger.WithError(err).Errorf("failed to get consent request")
					return nil, err
				}
				_consentRequest = consentRequestResp.Payload
				return _consentRequest, nil
			}

			wg := sync.WaitGroup{}
			var queue = make(chan func(), 100)
			defer close(queue)
			go func() {
				for f := range queue {
					f()
					wg.Done()
				}
			}()

			var enqueue = func(fn func()) {
				wg.Add(1)
				queue <- fn
			}

			authHandlers := authrequest.Handlers{
				OnFailed: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
					w.Write([]byte("error"))
				},
				OnLoginRequested:               handleLoginRequested(w, req, userSession, enqueue, getLoginRequest),
				OnRefreshingIdentity:           handleRefreshingIdentity(w, req, userSession, idpStore, selfURL, enqueue, getLoginRequest),
				OnPromptingForIdentifier:       handlePromptingForIdentifier(w, req, userSession, enqueue),
				OnValidatingIdentifier:         handleValidatingIdentifier(w, req, userSession, enqueue),
				OnFindingAuthMethod:            handleFindingAuthMethod(w, req, userSession, enqueue),
				OnPromptingForIdentityProvider: handlePromptingForIdentityProvider(w, req, userSession, enqueue, idpStore, ctx, orgStore),
				OnUseIdentityProvider:          handleUseIdentityProvider(w, req, userSession, enqueue, getLoginRequest, pathParameters, idpStore, ctx, selfURL),
				OnAwaitingIDPCallback:          handleAwaitingIDPCallback(w, req, userSession, enqueue),
				OnPerformingAuthCodeExchange:   handlePerformingAuthCodeExchange(w, req, userSession, enqueue, idpStore, ctx, selfURL),
				OnAuthCodeExchangeSucceeded:    handleAuthCodeExchangeSucceeded(w, req, userSession, enqueue, idpStore, ctx, loginStore),
				OnAwaitingConsentChallenge:     handleAwaitingConsentChallenge(w, req, hydraAdmin, userSession, enqueue),
				OnReceivedConsentChallenge:     handleReceivedConsentChallenge(w, req, userSession, enqueue, requestParameters, getConsentRequest),
				OnPresentingConsentChallenge:   handlePresentingConsentChallenge(w, req, userSession, enqueue, getConsentRequest),
				OnConsentRequestApproved:       handleConsentRequestApproved(w, req, userSession, enqueue, getConsentRequest, hydraAdmin, ctx),
				OnConsentRequestDeclined:       handleConsentRequestDenied(w, req, userSession, enqueue, hydraAdmin, ctx),
				OnApproved:                     handleApproved(w, req, userSession, enqueue),
				OnDeclined:                     handleDeclined(w, req, userSession, enqueue),
			}

			authRequest := getAuthRequest(action, authHandlers, userSession)

			logger.WithField("action", action).Tracef("dispatching login action")
			dispatchAction(w, action, enqueue, authRequest)

			wg.Wait()

			logger.Tracef("done dispatching login actions")

		}

	}
}

func getAuthRequest(action string, authHandlers authrequest.Handlers, userSession *sessions.Session) *authrequest.AuthRequest {
	var authRequest *authrequest.AuthRequest
	if action == authrequest.EventRequestLogin {
		prevAuthRequest := authrequest.CreateOrRestore(userSession, authHandlers)
		if prevAuthRequest.State() == authrequest.StateAccepted {
			return prevAuthRequest
		} else {
			authRequest = authrequest.NewAuthRequest(authHandlers)
		}
	} else {
		authRequest = authrequest.CreateOrRestore(userSession, authHandlers)
		if authRequest == nil {
			authRequest = authrequest.NewAuthRequest(authHandlers)
		}
	}
	return authRequest
}

func getUserSession(w http.ResponseWriter, req *http.Request, sessionStore sessions.Store) (*sessions.Session, bool) {
	userSession, err := sessionStore.Get(req, "login-session")
	if err != nil {
		logrus.WithError(err).Warnf("failed to restore session")
		userSession, err := sessionStore.New(req, "login-session")
		if err != nil {
			logrus.WithError(err).Warnf("failed to create new session session")
			handleError(w, http.StatusBadRequest, err)
			return nil, true
		}
		return userSession, false
	}
	return userSession, false
}
