package store

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

// IdentityProvider is a class that represents an Organization IdentityProvider.
// Organizations can register multiple IdentityProvider that they trust, through
// which their staff will be allowed to log in.
// For example, this could be an Okta, Auth0, GitHub identity provider.
//
// This struct is only used for storing types.IdentityProvider. This store
// maps the types.IdentityProvider to and from the store.IdentityProvider.
// This allows us to have flexibility into how we store the IdentityProvider
// and how we present it to the API.
type IdentityProvider struct {

	// ID is the Identity provider ID
	ID string

	// OrganizationID is the OrganizationID that owns this IdentityProvider
	OrganizationID string

	// Organization is the Organization that owns this IdentityProvider.
	// We include this field here because gorm.DB creates the Foreign Key
	// constraints when this field is present.
	Organization Organization

	// Domain represents the OIDC Issuer for this Identity Provider
	Domain string

	// ClientID represents the OAuth2 ClientID for this IdentityProvider
	// For example, an Okta ClientID
	ClientID string

	// ClientSecret represents the OAuth2 ClientSecret for the IdentityProvider
	// For example, an Okta ClientSecret
	ClientSecret string

	// EmailDomain represents the Domain of email addresses that should be
	// using this identity provider.
	// e.g. if EmailDomain = "my-org.com", then johndoe@my-org.com would use
	// this IdentityProvider
	EmailDomain string

	// Name of this IdentityProvider
	Name string

	Scopes string

	Claim Claim
}

type Claim struct {
	Version string
	Mappings map[string]string
}

func (c Claim) Value() (driver.Value, error) {
	j, err := json.Marshal(c)
	return j, err
}

func (c *Claim) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	m, ok := i.(map[string]interface{})

	if !ok {
		return errors.New("type assertion .(map[string]interface{}) failed")
	}

	versionIntf, ok := m["Version"]
	if ok {
		versionStr, ok := versionIntf.(string)
		if !ok {
			return fmt.Errorf("version is not a string")
		}
		c.Version = versionStr
	}

	claimMappingsInft, ok := m["Mappings"]

	c.Mappings = map[string]string{}

	if ok {
		claimMapping, ok := claimMappingsInft.(map[string]interface{})
		if ok {

			for key, elementIntf := range claimMapping {
				claimStr, ok := elementIntf.(string)
				if !ok {
					return fmt.Errorf("claim mapping is not a string")
				}
				c.Mappings[key] = claimStr
			}
		}
	}

	return nil
}

// IdentityProviders represent a list of IdentityProvider
type IdentityProviders []*IdentityProvider

// IdentityProviderListOptions represent the options when Listing IdentityProviders
type IdentityProviderListOptions struct {
	// ReturnClientSecret will include the IdentityProvider.ClientSecret in the response
	ReturnClientSecret bool
}

// IdentityProviderGetOptions represent the options when getting an IdentityProvider
type IdentityProviderGetOptions struct {
	// ReturnClientSecret will include the IdentityProvider.ClientSecret in the response
	ReturnClientSecret bool
}

// IdentityProviderCreateOptions represent the options when creating an IdentityProvider
type IdentityProviderCreateOptions struct {
	// ReturnClientSecret will include the IdentityProvider.ClientSecret in the response
	ReturnClientSecret bool
}

// IdentityProviderUpdateOptions represent the options when updating an IdentityProvider
type IdentityProviderUpdateOptions struct {
	// ReturnClientSecret will include the IdentityProvider.ClientSecret in the response
	ReturnClientSecret bool
}

// IdentityProviderStore is the store for IdentityProviders
type IdentityProviderStore interface {
	// List identity providers
	List(ctx context.Context, organizationID string, options IdentityProviderListOptions) ([]*types.IdentityProvider, error)
	// Get an IdentityProvider
	Get(ctx context.Context, identityProviderId string, options IdentityProviderGetOptions) (*types.IdentityProvider, error)
	// FindForEmailDomain finds IdentityProviders for the given email domain. Should only return one because we should not allow two
	// identityProvider to use the same email domain
	// TODO make this method return a single IdentityProvider
	FindForEmailDomain(ctx context.Context, emailDomain string, options IdentityProviderListOptions) ([]*types.IdentityProvider, error)
	// Create an IdentityProvider
	Create(ctx context.Context, identityProvider *types.IdentityProvider, options IdentityProviderCreateOptions) (*types.IdentityProvider, error)
	// Update an IdentityProvider
	Update(ctx context.Context, identityProvider *types.IdentityProvider, options IdentityProviderUpdateOptions) (*types.IdentityProvider, error)
}

// NewIdentityProviderStore returns a new IdentityProviderStore
func NewIdentityProviderStore(db Factory) IdentityProviderStore {
	return &identityProviderStore{db: db}
}

// identityProviderStore is the implementation of IdentityProviderStore
type identityProviderStore struct {
	db Factory
}

// Make sure identityProviderStore implements IdentityProviderStore
var _ IdentityProviderStore = &identityProviderStore{}

// List implements IdentityProviderStore.List
func (i identityProviderStore) List(ctx context.Context, organizationID string, options IdentityProviderListOptions) ([]*types.IdentityProvider, error) {
	ctx, db, l, done, err := actionContext(ctx, i.db, idpStoreName, "list", zap.String("organization_id", organizationID))
	if err != nil {
		return nil, err
	}
	defer done()

	var storeIdps IdentityProviders
	if err := db.WithContext(ctx).Find(&storeIdps, "organization_id = ?", organizationID).Error; err != nil {
		l.Error("failed to list identity providers", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	return mapIdentityProviderList(storeIdps, options.ReturnClientSecret), nil
}

// Get implements IdentityProviderStore.Get
func (i identityProviderStore) Get(ctx context.Context, identityProviderId string, options IdentityProviderGetOptions) (*types.IdentityProvider, error) {
	ctx, db, l, done, err := actionContext(ctx, i.db, idpStoreName, "get", zap.String("identity_provider_id", identityProviderId))
	if err != nil {
		return nil, err
	}
	defer done()

	var storeIdp *IdentityProvider
	if err := db.WithContext(ctx).First(&storeIdp, "id = ?", identityProviderId).Error; err != nil {
		l.Error("failed to list identity providers", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	return mapIdentityProviderTo(storeIdp, options.ReturnClientSecret), nil
}

// FindForEmailDomain implements IdentityProviderStore.FindForEmailDomain
func (i identityProviderStore) FindForEmailDomain(ctx context.Context, emailDomain string, options IdentityProviderListOptions) ([]*types.IdentityProvider, error) {
	ctx, db, l, done, err := actionContext(ctx, i.db, idpStoreName, "find_for_email_domain", zap.String("email_domain", emailDomain))
	if err != nil {
		return nil, err
	}
	defer done()

	var storeIdps IdentityProviders
	if err := db.WithContext(ctx).Find(&storeIdps, "email_domain = ?", emailDomain).Error; err != nil {
		l.Error("failed to find identity providers for email domain", zap.Error(err), zap.String("email_domain", emailDomain))
		return nil, meta.NewInternalServerError(err)
	}

	return mapIdentityProviderList(storeIdps, options.ReturnClientSecret), nil
}

// Create implements IdentityProviderStore.Create
func (i identityProviderStore) Create(ctx context.Context, identityProvider *types.IdentityProvider, options IdentityProviderCreateOptions) (*types.IdentityProvider, error) {
	ctx, db, l, done, err := actionContext(ctx, i.db, idpStoreName, "create")
	if err != nil {
		return nil, err
	}
	defer done()

	idp := mapIdentityProviderFrom(identityProvider)
	idp.ID = uuid.NewV4().String()

	if err := db.WithContext(ctx).Create(idp).Error; err != nil {
		l.Error("failed to create identity provider", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	return mapIdentityProviderTo(idp, options.ReturnClientSecret), nil
}

// Update implements IdentityProviderStore.Update
func (i identityProviderStore) Update(ctx context.Context, identityProvider *types.IdentityProvider, options IdentityProviderUpdateOptions) (*types.IdentityProvider, error) {
	ctx, db, l, done, err := actionContext(ctx, i.db, idpStoreName, "update", zap.String("identity_provider_id", identityProvider.ID))
	if err != nil {
		return nil, err
	}
	defer done()

	storeIdp := mapIdentityProviderFrom(identityProvider)

	// these are the only fields allowed to be updated on an IdentityProvider
	updates := map[string]interface{}{
		"name":         	identityProvider.Name,
		"client_id":    	identityProvider.ClientID,
		"domain":       	identityProvider.Domain,
		"email_domain": 	identityProvider.EmailDomain,
		"scopes":           identityProvider.Scopes,
		"claim":            identityProvider.Claim,
	}
	if len(identityProvider.ClientSecret) != 0 {
		updates["client_secret"] = identityProvider.ClientSecret
	}

	result := db.WithContext(ctx).Model(&IdentityProvider{}).Where("id = ?", storeIdp.ID).UpdateColumns(updates)
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

	return mapIdentityProviderTo(storeIdp, options.ReturnClientSecret), nil
}

// mapIdentityProviderList maps a list of store.IdentityProvider to a list of types.IdentityProvider
func mapIdentityProviderList(i IdentityProviders, keepClientSecrets bool) []*types.IdentityProvider {
	var result []*types.IdentityProvider
	for _, provider := range i {
		result = append(result, mapIdentityProviderTo(provider, keepClientSecrets))
	}
	if result == nil {
		result = []*types.IdentityProvider{}
	}
	return result
}

// mapIdentityProviderTo maps a store.IdentityProvider to a types.IdentityProvider
func mapIdentityProviderTo(i *IdentityProvider, keepClientSecret bool) *types.IdentityProvider {
	claim := &types.Claim{
		Version: i.Claim.Version,
		Mappings: i.Claim.Mappings,
	}
	result := &types.IdentityProvider{
		ID:             i.ID,
		OrganizationID: i.OrganizationID,
		Domain:         i.Domain,
		ClientID:       i.ClientID,
		EmailDomain:    i.EmailDomain,
		Name:           i.Name,
		Scopes:         i.Scopes,
		Claim:          *claim,
	}
	if keepClientSecret {
		result.ClientSecret = i.ClientSecret
	}
	return result
}

// mapIdentityProviderFrom maps a types.IdentityProvider into a store.IdentityProvider
func mapIdentityProviderFrom(i *types.IdentityProvider) *IdentityProvider {
	claim := &Claim{
		Version: i.Claim.Version,
		Mappings: i.Claim.Mappings,
	}
	return &IdentityProvider{
		ID:             i.ID,
		OrganizationID: i.OrganizationID,
		Domain:         i.Domain,
		ClientID:       i.ClientID,
		ClientSecret:   i.ClientSecret,
		EmailDomain:    i.EmailDomain,
		Name:           i.Name,
		Scopes:         i.Scopes,
		Claim:          *claim,
	}
}
