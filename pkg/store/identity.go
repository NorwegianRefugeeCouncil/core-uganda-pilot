package store

import (
	"context"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"go.uber.org/zap"
)

// Identity is a class that represents an Identity of a user from an Organisation.
//
// This struct is only used for storing types.Identity. This store
// maps the types.Identity to and from the store.Identity.
// This allows us to have flexibility into how we store the Identity
// and how we present it to the API.
type Identity struct {
	ID string
	Subject       string
	DisplayName   string
	FullName      string
	Email         string
	EmailVerified bool
}


// IdentityStore is the store for Identities
type IdentityStore interface {
	// Get an Identity
	Get(ctx context.Context, identityId string) (*types.Identity, error)
}


// NewIdentityStore returns a new IdentityProviderStore
func NewIdentityStore(db Factory) IdentityStore {
	return &identityStore{db: db}
}

// identityStore is the implementation of IdentityStore
type identityStore struct {
	db Factory
}

// Make sure identityStore implements IdentityStore
var _ IdentityStore = &identityStore{}


// Get implements IdentityStore.Get
func (i identityStore) Get(ctx context.Context, identityId string) (*types.Identity, error) {
	ctx, db, l, done, err := actionContext(ctx, i.db, idpStoreName, "get", zap.String("identity_id", identityId))
	if err != nil {
		return nil, err
	}
	defer done()

	var identity *Identity
	if err := db.WithContext(ctx).First(&identity, "id = ?", identityId).Error; err != nil {
		l.Error("failed to list identities", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	return mapIdentityTo(identity), nil
}

// mapIdentityTo maps a store.Identity to a types.Identity
func mapIdentityTo(i *Identity) *types.Identity {
	identity := &types.Identity{
		ID            : i.ID,
		Subject       : i.Subject,
		DisplayName   : i.DisplayName,
		FullName      : i.FullName,
		Email         : i.Email,
		EmailVerified : i.EmailVerified,
	}
	return identity
}
