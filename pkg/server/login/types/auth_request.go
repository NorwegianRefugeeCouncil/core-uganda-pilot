package types

import "time"

type Claims struct {
	UserID            string
	Username          string
	PreferredUsername string
	Email             string
	EmailVerified     bool
	Groups            []string
}

type PKCE struct {
	CodeChallenge       string
	CodeChallengeMethod string
}

type AuthRequest struct {
	ID                  string
	ClientID            string
	ResponseType        string
	Scopes              []string
	RedirectURI         string
	Nonce               string
	State               string
	ForceApprovalPrompt bool
	Expiry              time.Time
	LoggedIn            bool
	Claims              Claims
	PKCE                PKCE
}
