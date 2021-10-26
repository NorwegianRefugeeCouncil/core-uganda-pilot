package store

import (
	"context"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/types"
	uuid "github.com/satori/go.uuid"
)

type IdentityProvider struct {
	ID             string
	OrganizationID string
	Organization   Organization
	Kind           types.IdentityProviderKind
	Domain         string
	ClientID       string
	ClientSecret   string
}

type IdentityProviderListOptions struct {
	ReturnClientSecret bool
}

type IdentityProviderCreateOptions struct {
	ReturnClientSecret bool
}

type IdentityProviderUpdateOptions struct {
	ReturnClientSecret bool
}

type IdentityProviderStore interface {
	List(ctx context.Context, organizationID string, options IdentityProviderListOptions) ([]*types.IdentityProvider, error)
	Create(ctx context.Context, identityProvidr *types.IdentityProvider, options IdentityProviderCreateOptions) (*types.IdentityProvider, error)
	Update(ctx context.Context, identityProvidr *types.IdentityProvider, options IdentityProviderUpdateOptions) (*types.IdentityProvider, error)
}

func NewIdentityProviderStore(db Factory) IdentityProviderStore {
	return &identityProviderStore{db: db}
}

type identityProviderStore struct {
	db Factory
}

var _ IdentityProviderStore = &identityProviderStore{}

func (i identityProviderStore) List(ctx context.Context, organizationID string, options IdentityProviderListOptions) ([]*types.IdentityProvider, error) {
	db, err := i.db.Get()
	if err != nil {
		return nil, err
	}

	var idps []*IdentityProvider
	if err := db.WithContext(ctx).Find(&idps, "organization_id = ?", organizationID).Error; err != nil {
		return nil, meta.NewInternalServerError(err)
	}
	return mapList(idps, options.ReturnClientSecret), nil
}

func (i identityProviderStore) Create(ctx context.Context, identityProvidr *types.IdentityProvider, options IdentityProviderCreateOptions) (*types.IdentityProvider, error) {
	db, err := i.db.Get()
	if err != nil {
		return nil, err
	}

	idp := mapFrom(identityProvidr)
	idp.ID = uuid.NewV4().String()
	if err := db.WithContext(ctx).Create(idp).Error; err != nil {
		return nil, meta.NewInternalServerError(err)
	}
	return mapTo(idp, options.ReturnClientSecret), nil
}

func (i identityProviderStore) Update(ctx context.Context, identityProvidr *types.IdentityProvider, options IdentityProviderUpdateOptions) (*types.IdentityProvider, error) {
	db, err := i.db.Get()
	if err != nil {
		return nil, err
	}

	idp := mapFrom(identityProvidr)
	updates := map[string]interface{}{
		"client_id": identityProvidr.ClientID,
		"domain":    identityProvidr.Domain,
		"kind":      identityProvidr.Kind,
	}
	if len(identityProvidr.ClientSecret) != 0 {
		updates["client_secret"] = identityProvidr.ClientSecret
	}
	result := db.WithContext(ctx).Model(&IdentityProvider{}).Where("id = ?", idp.ID).UpdateColumns(updates)
	if err := result.Error; err != nil {
		return nil, meta.NewInternalServerError(err)
	}
	if result.RowsAffected == 0 {
		return nil, meta.NewNotFound(meta.GroupResource{
			Group:    "nrc.no",
			Resource: "identityproviders",
		}, identityProvidr.ID)
	}
	return mapTo(idp, options.ReturnClientSecret), nil
}

func mapList(i []*IdentityProvider, keepClientSecrets bool) []*types.IdentityProvider {
	var result []*types.IdentityProvider
	for _, provider := range i {
		result = append(result, mapTo(provider, keepClientSecrets))
	}
	return result
}

func mapTo(i *IdentityProvider, keepClientSecret bool) *types.IdentityProvider {
	result := &types.IdentityProvider{
		ID:             i.ID,
		OrganizationID: i.OrganizationID,
		Kind:           i.Kind,
		Domain:         i.Domain,
		ClientID:       i.ClientID,
	}
	if keepClientSecret {
		result.ClientSecret = i.ClientSecret
	}
	return result
}

func mapFrom(i *types.IdentityProvider) *IdentityProvider {
	return &IdentityProvider{
		ID:             i.ID,
		OrganizationID: i.OrganizationID,
		Kind:           i.Kind,
		Domain:         i.Domain,
		ClientID:       i.ClientID,
		ClientSecret:   i.ClientSecret,
	}
}