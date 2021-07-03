package server

import (
	"context"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"strings"
)

func createOauthClient(
	ctx context.Context,
	hydraAdmin admin.ClientService,
	cli *models.OAuth2Client) error {
	_, err := hydraAdmin.CreateOAuth2Client(&admin.CreateOAuth2ClientParams{
		Body:    cli,
		Context: ctx,
	})
	if err != nil {
		if strings.Contains(err.Error(), "createOAuth2ClientConflict") {
			_, err = hydraAdmin.UpdateOAuth2Client(&admin.UpdateOAuth2ClientParams{
				ID:      cli.ClientID,
				Body:    cli,
				Context: ctx,
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
