package types

type IdentityProviderKind string

const (
	IdentityProviderOIDC        = "oidc"
	IdentityProviderCredentials = "credentials"
)

type IdentityProvider struct {
	ID             string               `json:"id"`
	OrganizationID string               `json:"organizationId"`
	Kind           IdentityProviderKind `json:"kind"`
	Domain         string               `json:"domain"`
	ClientID       string               `json:"clientId"`
	ClientSecret   string               `json:"clientSecret"`
}

func (l *IdentityProvider) ClearClientSecret() {
	l.ClientSecret = ""
}

type IdentityProviderList struct {
	Items []*IdentityProvider `json:"items"`
}

func (l *IdentityProviderList) ClearClientSecrets() {
	for _, item := range l.Items {
		item.ClearClientSecret()
	}
}
