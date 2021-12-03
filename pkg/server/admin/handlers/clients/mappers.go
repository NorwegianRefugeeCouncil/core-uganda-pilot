package clients

import (
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/ory/hydra-client-go/models"
)

func mapToHydraClient(client types.OAuth2Client) *models.OAuth2Client {
	return &models.OAuth2Client{
		ClientID:                client.ID,
		AllowedCorsOrigins:      client.AllowedCorsOrigins,
		ClientName:              client.Name,
		ClientURI:               client.URI,
		GrantTypes:              client.GrantTypes,
		RedirectUris:            client.RedirectURIs,
		ResponseTypes:           client.ResponseTypes,
		Scope:                   client.Scope,
		TokenEndpointAuthMethod: client.TokenEndpointAuthMethod,
	}
}

func mapFromHydraClient(client *models.OAuth2Client) *types.OAuth2Client {
	return &types.OAuth2Client{
		ID:                      client.ClientID,
		AllowedCorsOrigins:      client.AllowedCorsOrigins,
		Name:                    client.ClientName,
		Secret:                  client.ClientSecret,
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
		Items: []*types.OAuth2Client{},
	}
	for _, hydraClient := range clients {
		clientList.Items = append(clientList.Items, mapFromHydraClient(hydraClient))
	}
	return &clientList
}
