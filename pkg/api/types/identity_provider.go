package types

type IdentityProviderKind string

type IdentityProvider struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	OrganizationID string `json:"organizationId"`
	Domain         string `json:"domain"`
	ClientID       string `json:"clientId"`
	ClientSecret   string `json:"clientSecret"`
	EmailDomain    string `json:"emailDomain"`
}

type IdentityProviderList struct {
	Items []*IdentityProvider `json:"items"`
}
