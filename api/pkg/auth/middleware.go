package auth

import (
	"context"
	"gopkg.in/square/go-jose.v2/jwt"
	"net/http"
)

type Roles struct {
	Roles []string `json:"roles"`
}

type ResourceAccessMap map[string]Roles

type Claims struct {
	jwt.Claims        `json:",inline"`
	PreferredUsername string            `json:"preferred_username"`
	EmailVerified     bool              `json:"email_verified"`
	Scope             string            `json:"scope"`
	DisplayName       string            `json:"display_name"`
	FullName          string            `json:"full_name"`
	Name              string            `json:"name"`
	RealmAccess       Roles             `json:"realm_access"`
	ResourceAccess    ResourceAccessMap `json:"resource_access"`
}

type UserInfo struct {
	Subject           string `json:"sub"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	FullName          string `json:"full_name"`
	DisplayName       string `json:"display_name"`
}

type AuthenticationContext struct {
	Claims          Claims
	IsAuthenticated bool
	AccessToken     string
}

const AuthenticationContextKey = "authentication_context"

func DangerouslySetDevAuthenticatedUserSubject(handler http.Handler, w http.ResponseWriter, req *http.Request) bool {
	ctx := req.Context()
	authUserSubject := req.Header.Get("X-Authenticated-User-Subject")
	if len(authUserSubject) != 0 {
		ctx = context.WithValue(ctx, "Subject", authUserSubject)
		req = req.WithContext(ctx)
		handler.ServeHTTP(w, req)
		return true
	}
	return false
}
