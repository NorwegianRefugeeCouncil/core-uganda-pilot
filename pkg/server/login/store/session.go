package store

import (
	"time"
)

type Session struct {
	ID              string
	Active          bool
	ExpiresAt       time.Time
	AuthenticatedAt time.Time
	IssuedAt        time.Time
	IdentityID      string
	Identity        Identity
}

type SessionAuthMethodKind string

const (
	SessionAuthMethodPassword = "password"
	SessionAuthMethodOidc     = "oidc"
)

type SessionAuthMethod struct {
	ID          string
	SessionID   string
	Session     Session
	Method      SessionAuthMethodKind
	CompletedAt time.Time
}
