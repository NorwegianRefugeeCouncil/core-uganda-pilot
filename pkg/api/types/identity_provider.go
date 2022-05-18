package types

import (
	"database/sql/driver"
	"encoding/json"
)

// IdentityProvider represents an Organization trusted Identity Provider
type IdentityProvider struct {
	// ID of the IdentityProvider
	ID string `json:"id"`
	// Name of the IdentityProvider
	Name string `json:"name"`
	// OrganizationID owning this IdentityProvider
	OrganizationID string `json:"organizationId"`
	// Domain OIDC issuer
	Domain string `json:"domain"`
	// ClientID is the OAuth2 client id
	ClientID string `json:"clientId"`
	// ClientSecret is the OAuth2 client secret
	ClientSecret string `json:"clientSecret"`
	// EmailDomain is the email domain "nrc.no" bound to this IdentityProvider
	// TODO: add unique constraint for email domains
	// TODO: add support for multiple email domains for a single IdentityProvider
	EmailDomain string `json:"emailDomain"`

	Scopes string `json:"scopes"`

	ClaimMappings ClaimMappings `json:"claimMappings"`
}

type ClaimMappings struct {
	Version string `json:"Version"`
	Mappings map[string]string `json:"Mappings"`
}

// IdentityProviderList represents a list of IdentityProvider
type IdentityProviderList struct {
	Items []*IdentityProvider `json:"items"`
}

func (c ClaimMappings) Value() (driver.Value, error) {
	j, err := json.Marshal(c)
	return j, err
}
