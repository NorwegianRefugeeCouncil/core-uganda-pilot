package login

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/rest"
)

type RESTLoginClient struct {
	c *rest.Client
}

var _ LoginClient = &RESTLoginClient{}

type SetCredentialPayload struct {
	PartyID           string `json:"partyId"`
	PlaintextPassword string `json:"plaintextPassword"`
}

func (r *RESTLoginClient) SetCredentials(ctx context.Context, partyId, plaintextPassword string) error {
	return r.c.Post().Body(&SetCredentialPayload{
		PartyID:           partyId,
		PlaintextPassword: plaintextPassword,
	}).Path("/apis/login/v1/credentials").Do(ctx).Error()
}
