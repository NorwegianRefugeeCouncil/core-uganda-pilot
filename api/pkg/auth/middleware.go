package auth

import (
	"gopkg.in/square/go-jose.v2/jwt"
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
	FamilyName        string            `json:"family_name"`
	GivenName         string            `json:"given_name"`
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
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
}

type AuthenticationContext struct {
	Claims          Claims
	IsAuthenticated bool
	AccessToken     string
}

const AuthenticationContextKey = "authentication_context"
