package login

import (
	"errors"
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	loginstore "github.com/nrc-no/core/pkg/server/login/store"
	"github.com/nrc-no/core/pkg/server/login/templates"
	"github.com/nrc-no/core/pkg/store"
	"github.com/nrc-no/core/pkg/utils/pointers"
	"github.com/nrc-no/core/pkg/utils/sets"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"sync"
)

func handleError(w http.ResponseWriter, status int, err error) {
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

			// getting user session
			userSession, done := getUserSession(w, req, sessionStore)
			if done {
				return
			}

			var _loginRequest *models.LoginRequest
			getLoginRequest := func(loginChallenge string) (*models.LoginRequest, error) {
				if _loginRequest != nil {
					return _loginRequest, nil
				}
				loginRequestResp, err := hydraAdmin.GetLoginRequest(&admin.GetLoginRequestParams{
					Context:        ctx,
					LoginChallenge: loginChallenge,
				})
				if err != nil {
					return nil, err
				}
				_loginRequest = loginRequestResp.Payload
				return _loginRequest, nil
			}

			var _consentRequest *models.ConsentRequest
			getConsentRequest := func(consentChallenge string) (*models.ConsentRequest, error) {
				if _consentRequest != nil {
					return _consentRequest, nil
				}
				consentRequestResp, err := hydraAdmin.GetConsentRequest(&admin.GetConsentRequestParams{
					Context:          ctx,
					ConsentChallenge: consentChallenge,
				})
				if err != nil {
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

			var addToQueue = func(fn func()) {
				wg.Add(1)
				queue <- fn
			}

			authHandlers := authrequest.Handlers{
				OnFailed: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
					w.Write([]byte("error"))
				},
				OnLoginRequested: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {

					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
					}

					// getting login request
					loginChallenge := req.URL.Query().Get("login_challenge")
					loginRequest, err := getLoginRequest(loginChallenge)
					if err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					authRequest.LoginChallenge = loginChallenge

					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					if loginRequest.Skip != nil && *loginRequest.Skip {
						addToQueue(func() {
							err := authRequest.SkipLoginRequest()
							if err != nil {
								addToQueue(func() {
									_ = authRequest.Fail(err)
								})
							}
						})
						return
					} else {
						addToQueue(func() {
							err := authRequest.PerformLogin()
							if err != nil {
								addToQueue(func() {
									_ = authRequest.Fail(err)
								})
							}
						})
						return
					}
				},
				OnPromptingForIdentifier: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}
					_ = templates.Template.ExecuteTemplate(w, "login_subject", map[string]interface{}{
						"Error": authRequest.IdentifierError,
					})
					return
				},
				OnValidatingIdentifier: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					// Parsing form data
					if err := req.ParseForm(); err != nil {

						addToQueue(func() {
							if err := authRequest.FailIdentifierValidation(err); err != nil {
								addToQueue(func() {
									_ = authRequest.Fail(err)
								})
							}
						})
						return
					}
					q := req.Form

					// Retrieving identifier from user form data
					email := q.Get("email")
					emailDomain, err := getEmailDomain(email)
					if err != nil {
						addToQueue(func() {
							if err := authRequest.FailIdentifierValidation(err); err != nil {
								addToQueue(func() {
									_ = authRequest.Fail(err)
								})
							}
						})
						return
					}

					authRequest.Identifier = email
					authRequest.EmailDomain = emailDomain

					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					addToQueue(func() {
						if err := authRequest.SucceedIdentifierValidation(); err != nil {
							addToQueue(func() {
								if err := authRequest.FailIdentifierValidation(err); err != nil {
									addToQueue(func() {
										_ = authRequest.Fail(err)
									})
								}
							})
						}
					})

				},
				OnFindingAuthMethod: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					addToQueue(func() {
						if err := authRequest.UseOidcAuth(); err != nil {
							addToQueue(func() {
								_ = authRequest.Fail(err)
							})
							return
						}
					})

				},
				OnPromptingForIdentityProvider: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					// Retrieving suitable identity provider for given identifier
					identityProviders, err := idpStore.FindForEmailDomain(ctx, authRequest.EmailDomain, store.IdentityProviderListOptions{
						ReturnClientSecret: false,
					})
					if err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}
					if len(identityProviders) == 0 {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					// ensuring a single organization match the given email domain
					organizationIDs := sets.NewString()
					for _, idp := range identityProviders {
						organizationIDs.Insert(idp.OrganizationID)
					}
					if len(organizationIDs) > 1 {
						addToQueue(func() {
							_ = authRequest.Fail(errors.New("email address domain conflict"))
						})
						return
					}

					// get organization
					organization, err := orgStore.Get(ctx, organizationIDs.List()[0])
					if err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					// prompt choosing identity provider
					err = templates.Template.ExecuteTemplate(w, "login_idp", map[string]interface{}{
						"OrganizationName":  organization.Name,
						"IdentityProviders": identityProviders,
					})

				},
				OnUseIdentityProvider: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					loginRequest, err := getLoginRequest(authRequest.LoginChallenge)
					if err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					identityProviderID := pathParameters["identityProviderId"]
					// getting identity provider
					idp, err := idpStore.Get(ctx, identityProviderID, store.IdentityProviderGetOptions{ReturnClientSecret: true})
					if err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}
					authRequest.IdentityProviderId = identityProviderID

					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					// creating state variable
					stateVar, err := createStateVariable()
					if err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}
					authRequest.StateVariable = stateVar

					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					// Getting Identity Provider Client Config
					oauth2Config, _, _, err := getOauthProvider(ctx, idp, selfURL, loginRequest)
					if err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					authCodeURL := oauth2Config.AuthCodeURL(stateVar)
					http.Redirect(w, req, authCodeURL, http.StatusSeeOther)

				},
				OnAwaitingIDPCallback: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}
					// noop
				},
				OnPerformingAuthCodeExchange: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					// getting identity provider
					idp, err := idpStore.Get(ctx, authRequest.IdentityProviderId, store.IdentityProviderGetOptions{ReturnClientSecret: true})
					if err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					// Getting Identity Provider Client Config
					oauth2Config, _, verifier, err := getOauthProvider(ctx, idp, selfURL, nil)
					if err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					// getting state from query
					stateFromQuery := req.URL.Query().Get("state")
					if len(stateFromQuery) == 0 {
						logrus.Warnf("state not found in callback query parameter")
						err := errors.New("state not found in response")
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					// getting authorization code from query
					authorizationCodeFromQuery := req.URL.Query().Get("code")
					if len(authorizationCodeFromQuery) == 0 {
						logrus.Warnf("auth code not found in callback query parameter")
						err := errors.New("code not found in response")
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					// exchange authorization code
					tokenFromExchange, err := oauth2Config.Exchange(req.Context(), authorizationCodeFromQuery)
					if err != nil {
						logrus.Warnf("failed to perform authorization code exchange: %v", err)
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					// getting id token from exchange
					rawIDTokenIntf := tokenFromExchange.Extra("id_token")
					if rawIDTokenIntf == nil {
						logrus.Warnf("id token not present in token: %v", err)
						var err = errors.New("id token not present in token exchange response")
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					// converting id token to string
					rawIDToken, ok := rawIDTokenIntf.(string)
					if !ok {
						logrus.Warnf("id token in response was not a string but was: %T", rawIDTokenIntf)
						var err = errors.New("id token in exchange response was not a string")
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					// verifying id token
					idToken, err := verifier.Verify(req.Context(), rawIDToken)
					if err != nil {
						logrus.Warnf("failed to verify ID token: %v", err)
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					authRequest.IDToken = rawIDToken
					authRequest.AccessToken = tokenFromExchange.AccessToken
					authRequest.RefreshToken = tokenFromExchange.RefreshToken

					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					var userProfile authrequest.Claims
					if err := idToken.Claims(&userProfile); err != nil {
						logrus.WithError(err).Warnf("failed to unmarshal claims from ID token")
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					authRequest.Claims = userProfile

					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					addToQueue(func() {
						if err := authRequest.SucceedAuthCodeExchange(); err != nil {
							addToQueue(func() {
								_ = authRequest.Fail(err)
							})
							return
						}
					})

				},
				OnAuthCodeExchangeSucceeded: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {

					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					// getting identity provider
					idp, err := idpStore.Get(ctx, authRequest.IdentityProviderId, store.IdentityProviderGetOptions{ReturnClientSecret: false})
					if err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					identifier, err := loginStore.FindOidcIdentifier(authRequest.Claims.Subject, idp.Domain)
					if err != nil {
						if meta.ReasonForError(err) == meta.StatusReasonNotFound {
							newIdentity, err := loginStore.CreateOidcIdentity(
								idp.Domain,
								authRequest.Claims.Subject,
								authRequest.AccessToken,
								authRequest.RefreshToken,
								authRequest.IDToken)
							if err != nil {
								logrus.Warnf("failed to create new oidc identity: %v", err)
								addToQueue(func() {
									_ = authRequest.Fail(err)
								})
								return
							}
							identifier = newIdentity.Credentials[0].Identifiers[0]
						} else {
							logrus.Warnf("could not retrieve oidc identifier for user: %v", err)
							addToQueue(func() {
								_ = authRequest.Fail(err)
							})
							return
						}
					}

					authRequest.Identity = identifier.Credential.Identity

					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					addToQueue(func() {
						if err := authRequest.AcceptLoginRequest(); err != nil {
							addToQueue(func() {
								_ = authRequest.Fail(err)
							})
							return
						}
					})

				},
				OnAcceptLoginRequest: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {

					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					acceptResp, err := hydraAdmin.AcceptLoginRequest(&admin.AcceptLoginRequestParams{
						Body: &models.AcceptLoginRequest{
							Acr:         "",
							Context:     nil,
							Remember:    false,
							RememberFor: 0,
							Subject:     pointers.String(authRequest.Identity.ID),
						},
						LoginChallenge: authRequest.LoginChallenge,
						Context:        req.Context(),
					})
					if err != nil {
						logrus.Warnf("could not accept login request: %v", err)
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					http.Redirect(w, req, *acceptResp.Payload.RedirectTo, http.StatusTemporaryRedirect)

				},
				OnAwaitingConsentChallenge: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}
				},
				OnReceivedConsentChallenge: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {

					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					// getting consent request
					consentChallenge := requestParameters.Get("consent_challenge")
					consentRequest, err := getConsentRequest(consentChallenge)
					if err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}
					authRequest.ConsentChallenge = consentChallenge

					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					if consentRequest.Skip {
						addToQueue(func() {
							if err := authRequest.SkipConsentRequest(); err != nil {
								addToQueue(func() {
									_ = authRequest.Fail(err)
								})
								return
							}
						})
					} else {
						addToQueue(func() {
							if err := authRequest.PresentConsentChallenge(); err != nil {
								addToQueue(func() {
									_ = authRequest.Fail(err)
								})
								return
							}
						})
					}
				},
				OnPresentingConsentChallenge: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					consentRequest, err := getConsentRequest(authRequest.ConsentChallenge)
					if err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					// prompt choosing identity provider
					err = templates.Template.ExecuteTemplate(w, "challenge", map[string]interface{}{
						"Scopes":     consentRequest.RequestedScope,
						"ClientName": consentRequest.Client.ClientName,
					})

				},
				OnConsentRequestApproved: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					consentRequest, err := getConsentRequest(authRequest.ConsentChallenge)
					if err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}
					if err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}
					resp, err := hydraAdmin.AcceptConsentRequest(&admin.AcceptConsentRequestParams{
						Context:          ctx,
						ConsentChallenge: authRequest.ConsentChallenge,
						Body: &models.AcceptConsentRequest{
							GrantAccessTokenAudience: consentRequest.RequestedAccessTokenAudience,
							GrantScope:               consentRequest.RequestedScope,
						},
					})
					if err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					authRequest.PostConsentURL = *resp.Payload.RedirectTo

					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					addToQueue(func() {
						if err := authRequest.Accept(); err != nil {
							addToQueue(func() {
								_ = authRequest.Fail(err)
							})
						}
					})

				},
				OnConsentRequestDeclined: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					resp, err := hydraAdmin.RejectConsentRequest(&admin.RejectConsentRequestParams{
						Context:          ctx,
						ConsentChallenge: authRequest.ConsentChallenge,
						Body:             &models.RejectRequest{},
					})
					if err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}
					authRequest.PostConsentURL = *resp.Payload.RedirectTo

					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}

					addToQueue(func() {
						if err := authRequest.Decline(); err != nil {
							addToQueue(func() {
								_ = authRequest.Fail(err)
							})
						}
					})

				},
				OnApproved: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}
					http.Redirect(w, req, authRequest.PostConsentURL, http.StatusTemporaryRedirect)
				},
				OnDeclined: func(authRequest *authrequest.AuthRequest, evt *fsm.Event) {
					if err := authRequest.Save(w, req, userSession); err != nil {
						addToQueue(func() {
							_ = authRequest.Fail(err)
						})
						return
					}
					http.Redirect(w, req, authRequest.PostConsentURL, http.StatusTemporaryRedirect)
				},
			}

			var authRequest *authrequest.AuthRequest
			if action == authrequest.EventRequestLogin {
				authRequest = authrequest.NewAuthRequest(authHandlers)
			} else {
				authRequest = authrequest.CreateOrRestore(userSession, authHandlers)
				if authRequest == nil {
					authRequest = authrequest.NewAuthRequest(authHandlers)
				}
			}

			switch action {
			case authrequest.EventRequestLogin:
				addToQueue(func() {
					if err := authRequest.RequestLogin(); err != nil {
						handleError(w, http.StatusBadRequest, err)
					}
				})

			case authrequest.EventProvideIdentifier:
				addToQueue(func() {
					if err := authRequest.ProvideIdentifier(); err != nil {
						handleError(w, http.StatusBadRequest, err)
					}
				})

			case authrequest.EventUseIdentityProvider:
				addToQueue(func() {
					if err := authRequest.UseIdentityProvider(); err != nil {
						handleError(w, http.StatusBadRequest, err)
					}
				})

			case authrequest.EventCallOidcCallback:
				addToQueue(func() {
					if err := authRequest.CallIdentityProviderCallback(); err != nil {
						handleError(w, http.StatusBadRequest, err)
					}
				})

			case authrequest.EventReceiveConsentChallenge:
				addToQueue(func() {
					if err := authRequest.ReceiveConsentChallenge(); err != nil {
						handleError(w, http.StatusBadRequest, err)
					}
				})

			case authrequest.EventApproveConsentChallenge:
				addToQueue(func() {
					if err := authRequest.ApproveConsentRequest(); err != nil {
						handleError(w, http.StatusBadRequest, err)
					}
				})

			case authrequest.EventDeclineConsentChallenge:
				addToQueue(func() {
					if err := authRequest.DeclineConsentRequest(); err != nil {
						handleError(w, http.StatusBadRequest, err)
					}
				})

			case authrequest.EventProvidePassword:
				addToQueue(func() {
					if err := authRequest.ProvidePassword(); err != nil {
						handleError(w, http.StatusBadRequest, err)
					}
				})
			}

			wg.Wait()

		}

	}
}

func getUserSession(w http.ResponseWriter, req *http.Request, sessionStore sessions.Store) (*sessions.Session, bool) {
	userSession, err := sessionStore.Get(req, "login-session")
	if err != nil {
		logrus.WithError(err).Warnf("failed to get login request")
		handleError(w, http.StatusBadRequest, err)
		return nil, true
	}
	return userSession, false
}
