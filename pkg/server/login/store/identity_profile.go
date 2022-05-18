package store

type IdentityProfile struct {
	ID                 string
	IdentityProviderID string
	Claims             map[string]string `json:"Claims"`
}
