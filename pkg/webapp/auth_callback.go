package webapp

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/iam"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
)

type Claims struct {
	Subject       string `json:"sub"`
	DisplayName   string `json:"display_name"`
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func init() {
	gob.Register(&Claims{})
}

func (s *Server) Callback(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	conf := s.privateOauth2Config

	session, err := s.sessionManager.Get(req)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get session")
		s.Error(w, err)
		return
	}

	stateIntf, ok := session.Values["state"]
	if !ok {
		s.Error(w, fmt.Errorf("no state found in session: unable to read state from session values"))
		return
	}

	delete(session.Values, "state")
	session.Save(req, w)

	stateStr, ok := stateIntf.(string)
	if !ok {
		s.Error(w, fmt.Errorf("no state found in session: unable to get string of session state value"))
		return
	}

	sessionState := stateStr

	queryState := req.URL.Query().Get("state")

	if sessionState != queryState {
		err := errors.New("state mismatch")
		s.Error(w, err)
		return
	}

	code := req.URL.Query().Get("code")
	if len(code) == 0 {
		err := errors.New("code not found")
		s.Error(w, err)
		return
	}

	tokenCtx := context.WithValue(ctx, oauth2.HTTPClient, s.HydraHTTPClient)
	token, err := conf.Exchange(tokenCtx, code)
	if err != nil {
		err := fmt.Errorf("failed to exchange code: %v", err)
		s.Error(w, err)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		err := fmt.Errorf("no id token in response")
		s.Error(w, err)
		return
	}

	idToken, err := s.oidcVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		err := fmt.Errorf("failed to verify id token: %v", err)
		s.Error(w, err)
		return
	}

	var profile Claims
	if err := idToken.Claims(&profile); err != nil {
		err := fmt.Errorf("failed to unmarshal claim: %v", err)
		s.Error(w, err)
		return
	}

	session.Values["id-token"] = rawIDToken
	session.Values["access-token"] = token.AccessToken
	session.Values["refresh-token"] = token.RefreshToken

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	individual, err := iamClient.Individuals().Get(ctx, profile.Subject)
	if err != nil {
		err := fmt.Errorf("failed to retrieve individual: %v", err)
		s.Error(w, err)
		return
	}

	profile.DisplayName = individual.Get(iam.DisplayNameAttribute.ID)
	profile.FullName = individual.Get(iam.FullNameAttribute.ID)

	session.Values["profile"] = profile
	if err := session.Save(req, w); err != nil {
		logrus.WithError(err).Errorf("failed to save session")
		s.Error(w, err)
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)

}
