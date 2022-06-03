package store

type IdentityProfile struct {
	ID                 string
	IdentityProviderID string
	Subject            string
	DisplayName        string
	FullName           string
	Email              string
	EmailVerified      bool
}
