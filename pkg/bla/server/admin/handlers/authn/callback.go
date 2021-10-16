package authn

import (
	"encoding/gob"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/bla/constants"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) Callback() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		userSession, err := h.sessionStore.Get(req, constants.SessionKey)
		if err != nil {
			logrus.WithError(err).Error("failed to retrieve session from store")
			utils.ErrorResponse(w, meta.NewInternalServerError(err))
			return
		}

		stateFromSessionIntf, ok := userSession.Values[constants.SessionState]
		if !ok {
			logrus.Warnf("state was not present in session")
			utils.ErrorResponse(w, meta.NewInternalServerError(fmt.Errorf("no state found in session: unable to read state from session values")))
			return
		}

		delete(userSession.Values, "state")
		if err := userSession.Save(req, w); err != nil {
			logrus.WithError(err).Errorf("failed to save state into session")
			utils.ErrorResponse(w, meta.NewInternalServerError(err))
			return
		}

		stateFromSession, ok := stateFromSessionIntf.(string)
		if !ok {
			logrus.WithError(err).Warnf("state from session was not a string but was %T", stateFromSession)
			utils.ErrorResponse(w, meta.NewInternalServerError(fmt.Errorf("no state found in session: unable to get string of session state value")))
			return
		}

		if len(stateFromSession) == 0 {
			logrus.WithError(err).Warnf("empty session state on oauth callback")
			utils.ErrorResponse(w, meta.NewInternalServerError(fmt.Errorf("state in session was found to be empty")))
			return
		}

		stateFromQuery := req.URL.Query().Get(constants.QueryParamState)

		if stateFromQuery != stateFromSession {
			logrus.WithError(err).Warnf("auth state mismatch: expected '%s' got '%s'", stateFromSession, stateFromQuery)
			utils.ErrorResponse(w, meta.NewUnauthorized(fmt.Sprintf("state mismatch")))
			return
		}

		authorizationCode := req.URL.Query().Get(constants.QueryParamCode)
		if len(authorizationCode) == 0 {
			logrus.WithError(err).Warnf("authorization code not present in query")
			utils.ErrorResponse(w, meta.NewUnauthorized(fmt.Sprintf("authorization code not present in query")))
			return
		}

		tokenFromExchange, err := h.oauth2Config.Exchange(ctx, authorizationCode)
		if err != nil {
			logrus.WithError(err).Warnf("failed to exchange authorization code: %v", err)
			utils.ErrorResponse(w, meta.NewUnauthorized(fmt.Sprintf("failed to exchange authorization code")))
			return
		}

		rawIDTokenIntf := tokenFromExchange.Extra("id_token")
		if rawIDTokenIntf == nil {
			logrus.Warnf("id token not present in token: %v", err)
			utils.ErrorResponse(w, meta.NewUnauthorized(fmt.Sprintf("id token not present in token exchange response")))
			return
		}

		rawIDToken, ok := rawIDTokenIntf.(string)
		if !ok {
			logrus.Warnf("id token in response was not a string but was: %T", rawIDTokenIntf)
			utils.ErrorResponse(w, meta.NewUnauthorized(fmt.Sprintf("id token in exchange response was not a string")))
			return
		}

		idToken, err := h.tokenVerifier.Verify(ctx, rawIDToken)
		if err != nil {
			logrus.Warnf("failed to validate ID token: %v", err)
			utils.ErrorResponse(w, meta.NewUnauthorized(fmt.Sprintf("invalid id token")))
			return
		}

		var userProfile Claims
		if err := idToken.Claims(&userProfile); err != nil {
			logrus.Warnf("failed to unmarshal claims from ID token: %v", err)
			utils.ErrorResponse(w, meta.NewInternalServerError(err))
			return
		}

		userSession.Values[constants.SessionIDToken] = rawIDToken
		userSession.Values[constants.SessionRefreshToken] = tokenFromExchange.RefreshToken
		userSession.Values[constants.SessionProfile] = userProfile

		if err := userSession.Save(req, w); err != nil {
			logrus.Warnf("failed to save user session: %v", err)
			utils.ErrorResponse(w, meta.NewInternalServerError(err))
			return
		}

		desiredURL, ok := userSession.Values[constants.SessionDesiredURL]
		if ok {
			if desiredUrlStr, ok := desiredURL.(string); ok {
				delete(userSession.Values, constants.SessionDesiredURL)
				_ = userSession.Save(req, w)
				w.Header().Set("Location", desiredUrlStr)
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}
		}

		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return

	}
}

func init() {
	gob.Register(&Claims{})
}

func (h *Handler) RestfulCallback(request *restful.Request, response *restful.Response) {
	handler := h.Callback()
	handler(response.ResponseWriter, request.Request)
}

type Claims struct {
	Subject       string `json:"sub"`
	DisplayName   string `json:"display_name"`
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}
