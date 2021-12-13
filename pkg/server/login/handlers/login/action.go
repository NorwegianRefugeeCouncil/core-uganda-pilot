package login

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	loginstore "github.com/nrc-no/core/pkg/server/login/store"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"net/http"
	"net/url"
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
) func(action string, pathParameters map[string]string, requestParameters url.Values) http.HandlerFunc {

	return func(action string, pathParameters map[string]string, requestParameters url.Values) http.HandlerFunc {

		return func(w http.ResponseWriter, req *http.Request) {

			ctx := req.Context()

			l := logging.NewLogger(ctx).With(
				zap.Any("path_parameters", pathParameters),
				zap.Any("request_parameters", requestParameters))

			l.Debug("received request for login action")

			l.Debug("getting user session")
			userSession, done := getUserSession(w, req, sessionStore)
			if done {
				return
			}

			// cache login request
			var _loginRequest *models.LoginRequest
			getLoginRequest := func(loginChallenge string) (*models.LoginRequest, error) {

				if _loginRequest != nil {
					l.Debug("returning cached login request")
					return _loginRequest, nil
				}

				l.Debug("getting login request")
				loginRequestResp, err := hydraAdmin.GetLoginRequest(&admin.GetLoginRequestParams{
					Context:        ctx,
					LoginChallenge: loginChallenge,
				})
				if err != nil {
					l.Error("failed to get login request", zap.Error(err))
					return nil, err
				}
				_loginRequest = loginRequestResp.Payload
				return _loginRequest, nil
			}

			// cache consent request
			var _consentRequest *models.ConsentRequest
			getConsentRequest := func(consentChallenge string) (*models.ConsentRequest, error) {
				if _consentRequest != nil {
					l.Debug("using cached consent request")
					return _consentRequest, nil
				}

				l.Debug("getting consent request")
				consentRequestResp, err := hydraAdmin.GetConsentRequest(&admin.GetConsentRequestParams{
					Context:          ctx,
					ConsentChallenge: consentChallenge,
				})
				if err != nil {
					l.Error("failed go get consent request", zap.Error(err))
					return nil, err
				}
				_consentRequest = consentRequestResp.Payload
				return _consentRequest, nil
			}

			var events []string

			var authRequest *authrequest.AuthRequest

			dispatch := func(evt string) {
				events = append(events, evt)
			}

			//goland:noinspection ALL
			reqScheme := "http://"
			if req.TLS != nil || req.Header.Get("X-Forwarded-Proto") == "https" {
				reqScheme = "https://"
			}
			selfURL := reqScheme + req.Host

			loginRequestedHandler := handleLoginRequested(ctx, req.URL.Query(), dispatch, getLoginRequest)
			refreshingIdentityHandler := handleRefreshingIdentity(ctx, idpStore, selfURL, dispatch, getLoginRequest)
			promptingForIdentifierHandler := handlePromptingForIdentifier(w, userSession, req)
			validatingIdentifierHandler := handleValidatingIdentifier(ctx, req, dispatch)
			findingAuthMethodHandler := handleFindingAuthMethod(ctx, dispatch)
			promptingForIdentityProviderHandler := handlePromptingForIdentityProvider(ctx, w, req, userSession, idpStore, orgStore)
			useIdentityProviderHandler := handleUseIdentityProvider(w, req, getLoginRequest, pathParameters, idpStore, selfURL)
			awaitingIdpCallbackHandler := handleAwaitingIDPCallback()
			performingAuthCodeExchangeHandler := handlePerformingAuthCodeExchange(req, dispatch, idpStore, selfURL)
			authCodeExchangeHandler := handleAuthCodeExchangeSucceeded(req, dispatch, idpStore, loginStore)
			awaitingConsentChallengeHandler := handleAwaitingConsentChallenge(w, req, hydraAdmin)
			receivedConsentChallengeHandler := handleReceivedConsentChallenge(ctx, dispatch, requestParameters, getConsentRequest)
			presentingConsentChallengeHandler := handlePresentingConsentChallenge(w, req, userSession, getConsentRequest)
			consentRequestApprovedHandler := handleConsentRequestApproved(ctx, dispatch, getConsentRequest, hydraAdmin)
			consentRequestDeclinedHandler := handleConsentRequestDenied(ctx, dispatch, hydraAdmin)
			approvedHandler := handleApproved(w, req)
			declinedHandler := handleDeclined(w, req)

			callbacks := map[string]fsm.Callback{
				authrequest.StateFailed: func(evt *fsm.Event) {
					utils.ErrorResponse(w, meta.NewUnauthorized("not authorized"))
					return
				},
				authrequest.StateLoginRequested: func(evt *fsm.Event) {
					if err := loginRequestedHandler(authRequest, evt); err != nil {
						dispatch(authrequest.EventFail)
					}
				},
				authrequest.StateRefreshingIdentity: func(evt *fsm.Event) {
					if err := refreshingIdentityHandler(authRequest, evt); err != nil {
						dispatch(authrequest.EventFail)
					}
				},
				authrequest.StatePromptingForIdentifier: func(evt *fsm.Event) {
					if err := promptingForIdentifierHandler(authRequest, evt); err != nil {
						dispatch(authrequest.EventFail)
					}
				},
				authrequest.StateValidatingIdentifier: func(evt *fsm.Event) {
					if err := validatingIdentifierHandler(authRequest, evt); err != nil {
						dispatch(authrequest.EventFail)
					}
				},
				authrequest.StateFindingAuthMethod: func(evt *fsm.Event) {
					if err := findingAuthMethodHandler(authRequest, evt); err != nil {
						dispatch(authrequest.EventFail)
					}
				},
				authrequest.StatePromptingForIdentityProvider: func(evt *fsm.Event) {
					if err := promptingForIdentityProviderHandler(authRequest, evt); err != nil {
						dispatch(authrequest.EventFail)
					}
				},
				authrequest.EventUseIdentityProvider: func(evt *fsm.Event) {
					if err := useIdentityProviderHandler(authRequest, evt); err != nil {
						dispatch(authrequest.EventFail)
					}
				},
				authrequest.StateAwaitingIdpCallback: func(evt *fsm.Event) {
					if err := awaitingIdpCallbackHandler(authRequest, evt); err != nil {
						dispatch(authrequest.EventFail)
					}
				},
				authrequest.StatePerformingAuthCodeExchange: func(evt *fsm.Event) {
					if err := performingAuthCodeExchangeHandler(authRequest, evt); err != nil {
						dispatch(authrequest.EventFail)
					}
				},
				authrequest.StateAuthCodeExchangeSucceeded: func(evt *fsm.Event) {
					if err := authCodeExchangeHandler(authRequest, evt); err != nil {
						dispatch(authrequest.EventFail)
					}
				},
				authrequest.StateAwaitingConsentChallenge: func(evt *fsm.Event) {
					if err := awaitingConsentChallengeHandler(authRequest, evt); err != nil {
						dispatch(authrequest.EventFail)
					}
				},
				authrequest.StateReceivedConsentChallenge: func(evt *fsm.Event) {
					if err := receivedConsentChallengeHandler(authRequest, evt); err != nil {
						dispatch(authrequest.EventFail)
					}
				},
				authrequest.StatePresentingConsent: func(evt *fsm.Event) {
					if err := presentingConsentChallengeHandler(authRequest, evt); err != nil {
						dispatch(authrequest.EventFail)
					}
				},
				authrequest.StateConsentRequestApproved: func(evt *fsm.Event) {
					if err := consentRequestApprovedHandler(authRequest, evt); err != nil {
						dispatch(authrequest.EventFail)
					}
				},
				authrequest.StateConsentRequestDeclined: func(evt *fsm.Event) {
					if err := consentRequestDeclinedHandler(authRequest, evt); err != nil {
						dispatch(authrequest.EventFail)
					}
				},
				authrequest.StateAccepted: func(evt *fsm.Event) {
					if err := approvedHandler(authRequest, evt); err != nil {
						dispatch(authrequest.EventFail)
					}
				},
				authrequest.StateDeclined: func(evt *fsm.Event) {
					if err := declinedHandler(authRequest, evt); err != nil {
						dispatch(authrequest.EventFail)
					}
				},
			}

			authRequest = getAuthRequest(action, callbacks, userSession)
			dispatch(action)

			i := -1
			for {
				i++
				if i > len(events)-1 {
					break
				}
				evt := events[i]

				l.Info("dispatching action",
					zap.String("action", evt),
					zap.String("state", authRequest.State()))

				if err := authRequest.Event(evt); err != nil {
					l.Error("failed to dispatch action",
						zap.Error(err),
						zap.String("action", evt),
						zap.String("state", authRequest.State()))
					if authRequest.State() == authrequest.StateFailed {
						utils.ErrorResponse(w, meta.NewBadRequest("invalid action"))
						break
					}
					dispatch(authrequest.EventFail)

				}
				l.Info("entered state", zap.String("state", authRequest.State()))
			}
			if err := authRequest.Save(w, req, userSession); err != nil {
				l.Error("failed to save session", zap.Error(err))
				dispatch(authrequest.EventFail)
			}
		}
	}
}

func getAuthRequest(action string, authHandlers fsm.Callbacks, userSession *sessions.Session) *authrequest.AuthRequest {
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
		if cookieErr, ok := err.(securecookie.MultiError); ok {
			if !cookieErr.IsDecode() {
				logrus.WithError(err).Errorf("failed to retrieve user session: %s", err)
				handleError(w, http.StatusBadRequest, err)
				return nil, true
			}
		}
		if err := userSession.Save(req, w); err != nil {
			logrus.WithError(err).Errorf("failed to clear user session: %s", err)
			handleError(w, http.StatusBadRequest, err)
			return nil, true
		}
	}
	return userSession, false
}
