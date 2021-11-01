package types

type Oauth2Client struct {
	ID                      string   `json:"id"`
	ClientName              string   `json:"clientName"`
	URI                     string   `json:"uri"`
	GrantTypes              []string `json:"grantTypes"`
	ResponseTypes           []string `json:"responseTypes"`
	Scope                   string   `json:"scope"`
	RedirectURIs            []string `json:"redirectUris"`
	AllowedCorsOrigins      []string `json:"allowedCorsOrigins"`
	TokenEndpointAuthMethod string   `json:"tokenEndpointAuthMethod"`
}

type Oauth2ClientList struct {
	Items []*Oauth2Client `json:"items"`
}
