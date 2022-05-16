package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
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

	Claim Claim `json:"claim"`
}

type Claim struct {
	Version string `json:"Version"`
	Mappings map[string]string `json:"Mappings"`
}

// IdentityProviderList represents a list of IdentityProvider
type IdentityProviderList struct {
	Items []*IdentityProvider `json:"items"`
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
