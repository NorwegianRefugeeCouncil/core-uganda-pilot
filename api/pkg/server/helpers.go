package server

import (
	"context"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"net/http"
	"strings"
)

func createOauthClient(
	ctx context.Context,
	hydraAdmin admin.ClientService,
	httpClient *http.Client,
	cli *models.OAuth2Client) error {
	_, err := hydraAdmin.CreateOAuth2Client(&admin.CreateOAuth2ClientParams{
		Body:       cli,
		Context:    ctx,
		HTTPClient: httpClient,
	})
	if err != nil {
		if strings.Contains(err.Error(), "createOAuth2ClientConflict") {
			_, err = hydraAdmin.UpdateOAuth2Client(&admin.UpdateOAuth2ClientParams{
				ID:         cli.ClientID,
				Body:       cli,
				Context:    ctx,
				HTTPClient: httpClient,
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
