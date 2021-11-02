package store

import (
	"time"
)

type CredentialKind string

const (
	PasswordCredential CredentialKind = "password"
	OidcCredential     CredentialKind = "oidc"
)

type Credential struct {
	ID                  string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	IdentityID          string
	Identity            *Identity
	Kind                CredentialKind
	HashedPassword      *string
	Issuer              *string
	Identifiers         []*CredentialIdentifier
	InitialAccessToken  *string
	InitialRefreshToken *string
	InitialIdToken      *string
}
