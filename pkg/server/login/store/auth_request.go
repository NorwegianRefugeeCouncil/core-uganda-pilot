package store

import (
	"github.com/nrc-no/core/pkg/store"
	"time"
)

type AuthRequest struct {
	// ID is the unique id of the auth request
	ID string
	// CreatedAt is the time the AuthnzBouncer request was created
	CreatedAt time.Time
	// ClientID is the client id requesting authentication
	ClientID string
	// Identifier is the user-provided identifier
	Identifier string
	// IdentityProviderID is the upstream identity provider id
	IdentityProviderID *string
	// IdentityProvider is the upstream identity provider
	IdentityProvider *store.IdentityProvider
	// Scope is the client requested scope
	Scope string
	// LoginChallenge is the hydra login challenge
	LoginChallenge *string
	// ConsentChallenge is the hydra consent challenge
	ConsentChallenge *string
	// State is the oidc state variable
	State *string
	// IdentityID is the found identity ID
	IdentityID *string
	// Identity is the found identity
	Identity Identity
}
