package clients

import (
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/ory/hydra-client-go/models"
)

func mapToHydraClient(client types.Oauth2Client) *models.OAuth2Client {
	return &models.OAuth2Client{
		ClientID:                client.ID,
		AllowedCorsOrigins:      client.AllowedCorsOrigins,
		ClientName:              client.ClientName,
		ClientURI:               client.URI,
		GrantTypes:              client.GrantTypes,
		RedirectUris:            client.RedirectURIs,
		ResponseTypes:           client.ResponseTypes,
		Scope:                   client.Scope,
		TokenEndpointAuthMethod: client.TokenEndpointAuthMethod,
	}
}

func mapFromHydraClient(client *models.OAuth2Client) *types.Oauth2Client {
	return &types.Oauth2Client{
		ID:                      client.ClientID,
		AllowedCorsOrigins:      client.AllowedCorsOrigins,
		ClientName:              client.ClientName,
		ClientSecret:            client.ClientSecret,
		URI:                     client.ClientURI,
		GrantTypes:              client.GrantTypes,
		RedirectURIs:            client.RedirectUris,
		ResponseTypes:           client.ResponseTypes,
		Scope:                   client.Scope,
		TokenEndpointAuthMethod: client.TokenEndpointAuthMethod,
	}
}

func mapFromHydraClients(clients []*models.OAuth2Client) *types.Oauth2ClientList {
	clientList := types.Oauth2ClientList{
		Items: []*types.Oauth2Client{},
	}
	for _, hydraClient := range clients {
		clientList.Items = append(clientList.Items, mapFromHydraClient(hydraClient))
	}
	return &clientList
}
