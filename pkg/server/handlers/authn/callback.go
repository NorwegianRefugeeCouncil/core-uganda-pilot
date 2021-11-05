package authn

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"net/http"
)

func (h *Handler) Callback(redirectURL string, sessionKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := logging.WithOperation(req.Context(), "callback")
		l := logging.NewLogger(ctx)

		l.Debug("beginning oidc callback")

		l.Debug("getting user session")
		userSession, err := getSession(w, req, h.sessionStore, sessionKey)
		if err != nil {
			l.With(zap.Error(err)).Error("failed to retrieve session from store")
			clearSession(w, req, h.sessionStore, sessionKey)
			utils.ErrorResponse(w, meta.NewInternalServerError(err))
			return
		}

		l.Debug("getting state from session")
		stateFromSessionIntf, ok := userSession.Values[constants.SessionState]
		if !ok {
			l.Warn("state was not present in session")
			clearSession(w, req, h.sessionStore, sessionKey)
			utils.ErrorResponse(w, meta.NewInternalServerError(fmt.Errorf("no state found in session: unable to read state from session values")))
			return
		}

		l.Debug("saving user session")
		delete(userSession.Values, "state")
		if err := userSession.Save(req, w); err != nil {
			l.Error("failed to save state into session", zap.Error(err))
			clearSession(w, req, h.sessionStore, sessionKey)
			utils.ErrorResponse(w, meta.NewInternalServerError(err))
			return
		}

		l.Debug("casting state to string")
		stateFromSession, ok := stateFromSessionIntf.(string)
		if !ok {
			l.Warn("state from session was not a string", zap.Error(err), zap.Any("unexpected_type", stateFromSessionIntf))
			clearSession(w, req, h.sessionStore, sessionKey)
			utils.ErrorResponse(w, meta.NewInternalServerError(fmt.Errorf("no state found in session: unable to get string of session state value")))
			return
		}

		if len(stateFromSession) == 0 {
			l.Warn("empty session state on oauth callback", zap.Error(err))
			clearSession(w, req, h.sessionStore, sessionKey)
			utils.ErrorResponse(w, meta.NewInternalServerError(fmt.Errorf("state in session was found to be empty")))
			return
		}

		l.Debug("getting state from query")
		stateFromQuery := req.URL.Query().Get(constants.QueryParamState)

		l.Debug("comparing query and session states")
		if stateFromQuery != stateFromSession {
			l.Warn("auth state mismatch", zap.String("expected", stateFromSession), zap.String("actual", stateFromQuery))
			clearSession(w, req, h.sessionStore, sessionKey)
			utils.ErrorResponse(w, meta.NewUnauthorized(fmt.Sprintf("state mismatch")))
			return
		}

		l.Debug("getting authorization code from query")
		authorizationCode := req.URL.Query().Get(constants.QueryParamCode)
		if len(authorizationCode) == 0 {
			l.Warn("authorization code not present in query")
			clearSession(w, req, h.sessionStore, sessionKey)
			utils.ErrorResponse(w, meta.NewUnauthorized(fmt.Sprintf("authorization code not present in query")))
			return
		}

		l.Debug("exchanging authorization code")
		tokenFromExchange, err := h.oauth2Config.Exchange(ctx, authorizationCode)
		if err != nil {
			l.Error("failed to exchange authorization code", zap.Error(err))
			clearSession(w, req, h.sessionStore, sessionKey)
			utils.ErrorResponse(w, meta.NewUnauthorized(fmt.Sprintf("failed to exchange authorization code")))
			return
		}

		l.Debug("verifying token")
		rawIDToken, idToken, err := verifyToken(ctx, h.tokenVerifier, tokenFromExchange)
		if err != nil {
			l.Error("failed to verify token", zap.Error(err))
			clearSession(w, req, h.sessionStore, sessionKey)
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("unmarshaling claims")
		var userProfile Claims
		if err := idToken.Claims(&userProfile); err != nil {
			l.Error("failed to unmarshal claims from id token", zap.Error(err))
			clearSession(w, req, h.sessionStore, sessionKey)
			utils.ErrorResponse(w, meta.NewInternalServerError(err))
			return
		}

		l.Debug("saving new session info")
		userSession.Values[constants.SessionIDToken] = rawIDToken
		userSession.Values[constants.SessionRefreshToken] = tokenFromExchange.RefreshToken
		userSession.Values[constants.SessionAccessToken] = tokenFromExchange.AccessToken
		userSession.Values[constants.SessionTokenExpiry] = tokenFromExchange.Expiry
		userSession.Values[constants.SessionTokenType] = tokenFromExchange.TokenType
		userSession.Values[constants.SessionProfile] = userProfile
		if err := userSession.Save(req, w); err != nil {
			l.Warn("failed to save user session", zap.Error(err))
			clearSession(w, req, h.sessionStore, sessionKey)
			utils.ErrorResponse(w, meta.NewInternalServerError(err))
			return
		}

		redirectUri := userSession.Values[constants.SessionDesiredURL].(string)

		w.Header().Set("Location", redirectUri)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return

	}
}

func verifyToken(ctx context.Context, verifier *oidc.IDTokenVerifier, tokenFromExchange *oauth2.Token) (string, *oidc.IDToken, error) {
	rawIDTokenIntf := tokenFromExchange.Extra("id_token")
	if rawIDTokenIntf == nil {
		logrus.Warnf("id token not present in token")
		return "", nil, meta.NewUnauthorized(fmt.Sprintf("id token not present in token exchange response"))
	}
	rawIDToken, ok := rawIDTokenIntf.(string)
	if !ok {
		logrus.Warnf("id token in response was not a string but was: %T", rawIDTokenIntf)
		return "", nil, meta.NewUnauthorized(fmt.Sprintf("id token in exchange response was not a string"))
	}
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		logrus.Warnf("failed to validate ID token: %v", err)
		return "", nil, meta.NewUnauthorized(fmt.Sprintf("invalid id token"))
	}
	return rawIDToken, idToken, nil
}

func (h *Handler) RestfulCallback(sessionKey, redirectURL string) restful.RouteFunction {
	return func(r *restful.Request, response *restful.Response) {
		h.Callback(redirectURL, sessionKey)(response.ResponseWriter, r.Request)
	}
}

type Claims struct {
	Subject       string `json:"sub"`
	DisplayName   string `json:"display_name"`
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}
