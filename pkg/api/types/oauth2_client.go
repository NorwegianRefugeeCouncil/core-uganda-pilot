package types

// OAuth2Client represents an OAuth2 client registered on Core.
// Each application that needs to authenticate users, or to authenticate itself,
// needs a dedicated OAuth2Client
type OAuth2Client struct {
	// ID of the OAuth2Client (Client ID)
	ID string `json:"id"`
	// Name of the OAuth2Client
	Name string `json:"clientName"`
	// Secret of the OAuth2Client
	Secret string `json:"clientSecret"`
	// URI is the application main URI
	URI string `json:"uri"`
	// GrantTypes of the OAuth2Client
	// see https://oauth.net/2/grant-types/
	GrantTypes []string `json:"grantTypes"`
	// ResponseTypes of the OAuth2Client
	// see https://openid.net/specs/oauth-v2-multiple-response-types-1_0.html
	ResponseTypes []string `json:"responseTypes"`
	// Scope is the OAuth2 scope
	// see https://oauth.net/2/scope/
	Scope string `json:"scope"`
	// RedirectURIs accepted by this OAuth2Client
	// see https://www.oauth.com/oauth2-servers/redirect-uris/
	RedirectURIs []string `json:"redirectUris"`
	// AllowedCorsOrigins for this OAuth2Client
	AllowedCorsOrigins []string `json:"allowedCorsOrigins"`
	// TokenEndpointAuthMethod fot this OAuth2Client
	TokenEndpointAuthMethod string `json:"tokenEndpointAuthMethod"`
}

// Oauth2ClientList represents a list of OAuth2Client
type Oauth2ClientList struct {
	Items []*OAuth2Client `json:"items"`
}
