package store

type IdentityState string

const (
	IdentityStateActive   IdentityState = "active"
	IdentityStateInactive IdentityState = "inactive"
)

type Identity struct {
	ID          string
	State       IdentityState
	Credentials []*Credential
}
