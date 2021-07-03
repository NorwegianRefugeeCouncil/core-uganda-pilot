package webapp

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/nrc-no/core/pkg/apps/iam"
	"golang.org/x/oauth2"
	"net/http"
)

type Claims struct {
	Subject       string `json:"sub"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func init() {
	gob.Register(&Claims{})
}

func (s *Server) Callback(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	conf := s.oauth2Config
	sessionState := s.sessionManager.PopString(ctx, "state")
	queryState := req.URL.Query().Get("state")

	if sessionState != queryState {
		http.Error(w, "state mismatch", http.StatusInternalServerError)
		return
	}

	code := req.URL.Query().Get("code")
	if len(code) == 0 {
		http.Error(w, "code not found", http.StatusInternalServerError)
		return
	}

	tokenCtx := context.WithValue(ctx, oauth2.HTTPClient, s.HydraHTTPClient)
	token, err := conf.Exchange(tokenCtx, code)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to exchange code: %v", err).Error(), http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "no id token in response", http.StatusInternalServerError)
		return
	}

	idToken, err := s.oidcVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "failed to verify id token", http.StatusInternalServerError)
		return
	}

	var profile Claims
	if err := idToken.Claims(&profile); err != nil {
		http.Error(w, "failed to unmarshal claims", http.StatusInternalServerError)
		return
	}

	s.sessionManager.Put(ctx, "id-token", rawIDToken)
	s.sessionManager.Put(ctx, "access-token", token.AccessToken)
	s.sessionManager.Put(ctx, "refresh-token", token.RefreshToken)

	individual, err := s.IAMClient(ctx).Individuals().Get(ctx, profile.Subject)
	if err != nil {
		http.Error(w, "failed to retrieve individual", http.StatusInternalServerError)
		return
	}

	profile.FamilyName = individual.Get(iam.LastNameAttribute.ID)
	profile.GivenName = individual.Get(iam.FirstNameAttribute.ID)

	s.sessionManager.Put(ctx, "profile", profile)

	http.Redirect(w, req, "/", http.StatusSeeOther)

}
