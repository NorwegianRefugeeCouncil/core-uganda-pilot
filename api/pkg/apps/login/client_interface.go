package login

import "context"

type Interface interface {
	Login() Client
}

type Client interface {
	SetCredentials(ctx context.Context, partyId, plaintextPassword string) error
}
