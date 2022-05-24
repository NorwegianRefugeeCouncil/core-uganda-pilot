package types

// Identity represents a user including the claims from the auth request
type Identity struct {
	ID            string `json:"id"`
	Subject       string `json:"subject"`
	DisplayName   string `json:"displayName"`
	FullName      string `json:"fullName"`
	Email         string `json:"email"`
	EmailVerified bool `json:"emailVerified"`
}
