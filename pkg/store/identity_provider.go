package store

import (
	"context"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type IdentityProvider struct {
	ID             string
	OrganizationID string
	Organization   Organization
	Domain         string
	ClientID       string
	ClientSecret   string
	EmailDomain    string
	Name           string
}

type IdentityProviderListOptions struct {
	ReturnClientSecret bool
}

type IdentityProviderGetOptions struct {
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
	Get(ctx context.Context, identityProviderId string, options IdentityProviderGetOptions) (*types.IdentityProvider, error)
	FindForEmailDomain(ctx context.Context, emailDomain string, options IdentityProviderListOptions) ([]*types.IdentityProvider, error)
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
	ctx, db, l, done, err := actionContext(ctx, i.db, "identity_provider", "list", zap.String("organization_id", organizationID))
	if err != nil {
		return nil, err
	}
	defer done()

	l.Debug("listing identity providers")
	var idps []*IdentityProvider
	if err := db.WithContext(ctx).Find(&idps, "organization_id = ?", organizationID).Error; err != nil {
		l.Error("failed to list identity providers", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}
	if idps == nil {
		idps = []*IdentityProvider{}
	}

	l.Debug("successfully listed identity providers")
	return mapList(idps, options.ReturnClientSecret), nil
}

func (i identityProviderStore) Get(ctx context.Context, identityProviderId string, options IdentityProviderGetOptions) (*types.IdentityProvider, error) {
	ctx, db, l, done, err := actionContext(ctx, i.db, "identity_provider", "get", zap.String("identity_provider_id", identityProviderId))
	if err != nil {
		return nil, err
	}
	defer done()

	l.Debug("listing identity providers")
	var idp *IdentityProvider
	if err := db.WithContext(ctx).First(&idp, "id = ?", identityProviderId).Error; err != nil {
		l.Error("failed to list identity providers", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	l.Debug("successfully listed identity providers")
	return mapTo(idp, options.ReturnClientSecret), nil
}

func (i identityProviderStore) FindForEmailDomain(ctx context.Context, emailDomain string, options IdentityProviderListOptions) ([]*types.IdentityProvider, error) {
	ctx, db, l, done, err := actionContext(ctx, i.db, "identity_provider", "find_for_email_domain", zap.String("email_domain", emailDomain))
	if err != nil {
		return nil, err
	}
	defer done()

	l.Debug("finding identity providers for email domain")
	var idps []*IdentityProvider
	if err := db.WithContext(ctx).Find(&idps, "email_domain = ?", emailDomain).Error; err != nil {
		l.Error("failed to find identity providers for email domain", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	l.Debug("successfully found identity providers for email domain", zap.Int("count", len(idps)))
	return mapList(idps, options.ReturnClientSecret), nil
}

func (i identityProviderStore) Create(ctx context.Context, identityProvidr *types.IdentityProvider, options IdentityProviderCreateOptions) (*types.IdentityProvider, error) {
	ctx, db, l, done, err := actionContext(ctx, i.db, "identity_provider", "create")
	if err != nil {
		return nil, err
	}
	defer done()

	idp := mapFrom(identityProvidr)
	idp.ID = uuid.NewV4().String()

	l.Debug("creating identity provider")
	if err := db.WithContext(ctx).Create(idp).Error; err != nil {
		l.Error("failed to create identity provider", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	l.Debug("successfully created identity provider")
	return mapTo(idp, options.ReturnClientSecret), nil
}

func (i identityProviderStore) Update(ctx context.Context, identityProvider *types.IdentityProvider, options IdentityProviderUpdateOptions) (*types.IdentityProvider, error) {
	ctx, db, l, done, err := actionContext(ctx, i.db, "identity_provider", "create", zap.String("identity_provider_id", identityProvider.ID))
	if err != nil {
		return nil, err
	}
	defer done()

	idp := mapFrom(identityProvider)
	updates := map[string]interface{}{
		"name":         identityProvider.Name,
		"client_id":    identityProvider.ClientID,
		"domain":       identityProvider.Domain,
		"email_domain": identityProvider.EmailDomain,
	}
	if len(identityProvider.ClientSecret) != 0 {
		updates["client_secret"] = identityProvider.ClientSecret
	}

	l.Debug("updating identity provider")
	result := db.WithContext(ctx).Model(&IdentityProvider{}).Where("id = ?", idp.ID).UpdateColumns(updates)
	if err := result.Error; err != nil {
		l.Error("failed to update identity provider", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}
	if result.RowsAffected == 0 {
		l.Error("identity provider not found")
		return nil, meta.NewNotFound(meta.GroupResource{
			Group:    "nrc.no",
			Resource: "identityproviders",
		}, identityProvider.ID)
	}

	l.Debug("successfully updated identity provider")
	return mapTo(idp, options.ReturnClientSecret), nil
}

func mapList(i []*IdentityProvider, keepClientSecrets bool) []*types.IdentityProvider {
	var result []*types.IdentityProvider
	for _, provider := range i {
		result = append(result, mapTo(provider, keepClientSecrets))
	}
	if result == nil {
		result = []*types.IdentityProvider{}
	}
	return result
}

func mapTo(i *IdentityProvider, keepClientSecret bool) *types.IdentityProvider {
	result := &types.IdentityProvider{
		ID:             i.ID,
		OrganizationID: i.OrganizationID,
		Domain:         i.Domain,
		ClientID:       i.ClientID,
		EmailDomain:    i.EmailDomain,
		Name:           i.Name,
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
		Domain:         i.Domain,
		ClientID:       i.ClientID,
		ClientSecret:   i.ClientSecret,
		EmailDomain:    i.EmailDomain,
		Name:           i.Name,
	}
}
