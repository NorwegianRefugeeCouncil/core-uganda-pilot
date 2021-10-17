package bla

type Party struct {
	ID string
}

type User struct {
	ID            string
	PartyID       string
	PasswordHash  string
	Email         string
	EmailVerified bool
}

type UserIdentity struct {
	ID         string
	UserID     string
	ProviderID string
	Attributes map[string]string
}

type IdentityProvider struct {
	ID string
}
