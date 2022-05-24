package store

import (
	"context"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"go.uber.org/zap"
)

// IdentityProfile is a class that represents an IdentityProfile of a user from an Organisation.
//
// This struct is only used for storing types.IdentityProfile. This store
// maps the types.IdentityProfile to and from the store.IdentityProfile.
// This allows us to have flexibility into how we store the IdentityProfile
// and how we present it to the API.
type IdentityProfile struct {
	ID string
	Subject       string
	DisplayName   string
	FullName      string
	Email         string
	EmailVerified bool
}


// IdentityProfileStore is the store for Identities
type IdentityProfileStore interface {
	// Get an IdentityProfile
	Get(ctx context.Context, identityId string) (*types.IdentityProfile, error)
}


// NewIdentityProfileStore returns a new IdentityProfileProviderStore
func NewIdentityProfileStore(db Factory) IdentityProfileStore {
	return &identityStore{db: db}
}

// identityStore is the implementation of IdentityProfileStore
type identityStore struct {
	db Factory
}

// Make sure identityStore implements IdentityProfileStore
var _ IdentityProfileStore = &identityStore{}


// Get implements IdentityProfileStore.Get
func (i identityStore) Get(ctx context.Context, identityId string) (*types.IdentityProfile, error) {
	ctx, db, l, done, err := actionContext(ctx, i.db, identityStoreName, "get", zap.String("identity_id", identityId))
	if err != nil {
		return nil, err
	}
	defer done()

	var identity *IdentityProfile
	if err := db.WithContext(ctx).First(&identity, "id = ?", identityId).Error; err != nil {
		l.Error("failed to list identities", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	return mapIdentityProfileTo(identity), nil
}

// mapIdentityProfileTo maps a store.IdentityProfile to a types.IdentityProfile
func mapIdentityProfileTo(i *IdentityProfile) *types.IdentityProfile {
	identity := &types.IdentityProfile{
		ID            : i.ID,
		Subject       : i.Subject,
		DisplayName   : i.DisplayName,
		FullName      : i.FullName,
		Email         : i.Email,
		EmailVerified : i.EmailVerified,
	}
	return identity
}
