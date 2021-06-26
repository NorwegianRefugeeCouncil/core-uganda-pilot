package login

import "context"

type Interface interface {
	Login() LoginClient
}

type LoginClient interface {
	SetCredentials(ctx context.Context, partyId, plaintextPassword string) error
}
