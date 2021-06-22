package keycloak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	baseUrl      *url.URL
	clientID     string
	clientSecret string
	realmName    string
}

func NewClient(baseURL, realmName, clientID, clientSecret string) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &Client{
		baseUrl:      parsedURL,
		clientSecret: clientSecret,
		clientID:     clientID,
		realmName:    realmName,
	}, nil
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (c *Client) GetToken() (string, error) {

	reqUrl := fmt.Sprintf("%s/auth/realms/%s/protocol/openid-connect/token", c.baseUrl.String(), c.realmName)

	data := url.Values{}
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("grant_type", "client_credentials")

	encoded := data.Encode()

	req, err := http.NewRequest("POST", reqUrl, strings.NewReader(encoded))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(encoded)))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}

		return "", fmt.Errorf("unexpected status code: %d: %s", res.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	tokenResponse := TokenResponse{}
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil

}

func (c *Client) CreateUser(token string, user *User) (*User, error) {

	reqUrl := fmt.Sprintf("%s/auth/admin/realms/%s/users", c.baseUrl.String(), c.realmName)

	userBytes, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", reqUrl, bytes.NewReader(userBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
		}
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, string(body))
	}

	locationParts := strings.Split(response.Header.Get("Location"), "/")
	userID := locationParts[len(locationParts)-1]

	return c.GetUser(token, userID)

}

func (c *Client) DeleteUser(token string, id string) error {

	reqUrl := fmt.Sprintf("%s/auth/admin/realms/%s/users/%s", c.baseUrl.String(), c.realmName, id)

	req, err := http.NewRequest("DELETE", reqUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusNotFound && response.StatusCode != http.StatusNoContent {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("unexpected status code: %d", response.StatusCode)
		}
		return fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, string(body))
	}

	return nil

}

func (c *Client) GetUser(token, id string) (*User, error) {

	reqUrl := fmt.Sprintf("%s/auth/admin/realms/%s/users/%s", c.baseUrl.String(), c.realmName, id)

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	ret := User{}
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil

}

func (c *Client) GetUserByUsername(token, username string) (*User, error) {
	users, err := c.ListUsers(token, ListUserOptions{
		Username: username,
	})

	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, fmt.Errorf("user with username '%s' not found", username)
}

func (c *Client) UpdateUser(token string, user *User) (*User, error) {

	reqUrl := fmt.Sprintf("%s/auth/admin/realms/%s/users/%s", c.baseUrl.String(), c.realmName, user.ID)

	userBytes, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", reqUrl, bytes.NewReader(userBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code")
	}

	return c.GetUser(token, user.ID)

}

type ListUserOptions struct {
	BriefRepresentation bool
	Email               string
	First               *int32
	FirstName           string
	LastName            string
	Max                 *int32
	Search              string
	Username            string
}

func (c *Client) ListUsers(token string, options ListUserOptions) ([]*User, error) {

	baseUrl := *c.baseUrl
	qryParams := baseUrl.Query()
	if options.BriefRepresentation {
		qryParams.Set("briefRepresentation", "true")
	}
	if len(options.Email) != 0 {
		qryParams.Set("email", options.Email)
	}
	if options.First != nil {
		qryParams.Set("first", strconv.FormatInt(int64(*options.First), 10))
	}
	if len(options.FirstName) != 0 {
		qryParams.Set("firstName", options.FirstName)
	}
	if len(options.LastName) != 0 {
		qryParams.Set("lastName", options.LastName)
	}
	if options.Max != nil {
		qryParams.Set("max", strconv.FormatInt(int64(*options.Max), 10))
	}
	if len(options.Search) != 0 {
		qryParams.Set("search", options.Search)
	}
	if len(options.Username) != 0 {
		qryParams.Set("username", options.Username)
	}

	reqURL := fmt.Sprintf("%s/auth/admin/realms/%s/users?%s", c.baseUrl.String(), c.realmName, qryParams.Encode())

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}
		return nil, fmt.Errorf("unexpected status code: %d: %s", res.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var ret []*User
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil

}

func (c *Client) CreateClientScope(token string, scope *ClientScope) (*ClientScope, error) {

	reqUrl := fmt.Sprintf("%s/auth/admin/realms/%s/client-scopes", c.baseUrl.String(), c.realmName)

	clientScopeBytes, err := json.Marshal(scope)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", reqUrl, bytes.NewReader(clientScopeBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
		}
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, string(body))
	}

	locationParts := strings.Split(response.Header.Get("Location"), "/")
	clientScopeID := locationParts[len(locationParts)-1]

	return c.GetClientScope(token, clientScopeID)

}

func (c *Client) ListClientScopes(token string) ([]*ClientScope, error) {

	baseUrl := *c.baseUrl
	reqUrl := fmt.Sprintf("%s/auth/admin/realms/%s/client-scopes", baseUrl.String(), c.realmName)

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var ret []*ClientScope
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil

}

func (c *Client) GetClientScope(token, id string) (*ClientScope, error) {

	reqUrl := fmt.Sprintf("%s/auth/admin/realms/%s/client-scopes/%s", c.baseUrl.String(), c.realmName, id)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var ret ClientScope
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil

}

func (c *Client) UpdateClientScope(token string, clientScope *ClientScope) (*ClientScope, error) {

	reqUrl := fmt.Sprintf("%s/auth/admin/realms/%s/client-scopes/%s",
		c.baseUrl.String(), c.realmName, clientScope.ID)

	bodyBytes, err := json.Marshal(clientScope)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", reqUrl, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusNoContent {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}
		return nil, fmt.Errorf("unexpected status code: %d: %s", res.StatusCode, string(body))
	}

	return c.GetClientScope(token, clientScope.ID)

}

func (c *Client) CreateProtocolMapper(token, clientScopeID string, protocolMapper *ProtocolMapper) (*ProtocolMapper, error) {

	reqUrl := fmt.Sprintf("%s/auth/admin/realms/%s/client-scopes/%s/protocol-mappers/models", c.baseUrl.String(), c.realmName, clientScopeID)

	protocolMapperBytes, err := json.Marshal(protocolMapper)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", reqUrl, bytes.NewReader(protocolMapperBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
		}
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, string(body))
	}

	locationParts := strings.Split(response.Header.Get("Location"), "/")
	protocolMapperID := locationParts[len(locationParts)-1]

	return c.GetProtocolMapper(token, clientScopeID, protocolMapperID)

}

func (c *Client) GetProtocolMapper(token, clientScopeID, id string) (*ProtocolMapper, error) {

	reqUrl := fmt.Sprintf("%s/auth/admin/realms/%s/client-scopes/%s/protocol-mappers/models/%s", c.baseUrl.String(), c.realmName, clientScopeID, id)

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var ret ProtocolMapper
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil

}

func (c *Client) UpdateProtocolMapper(token string, clientScopeID string, protocolMapper *ProtocolMapper) (*ProtocolMapper, error) {

	reqUrl := fmt.Sprintf("%s/auth/admin/realms/%s/client-scopes/%s/protocol-mappers/models/%s", c.baseUrl.String(), c.realmName, clientScopeID, protocolMapper.ID)

	bodyBytes, err := json.Marshal(protocolMapper)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", reqUrl, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusNoContent {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}
		return nil, fmt.Errorf("unexpected status code: %d: %s", res.StatusCode, string(body))
	}

	return c.GetProtocolMapper(token, clientScopeID, protocolMapper.ID)

}

func (c *Client) DeleteProtocolMapper(token, clientScopeID, ID string) error {

	reqUrl := fmt.Sprintf("%s/auth/admin/realms/%s/client-scopes/%s/protocol-mappers/models/%s", c.baseUrl.String(), c.realmName, clientScopeID, ID)

	req, err := http.NewRequest("DELETE", reqUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}
		return fmt.Errorf("unexpected status code: %d: %s", res.StatusCode, string(body))
	}

	return nil
}

type ClientScope struct {
	Attributes      map[string]string `json:"attributes,omitempty"`
	Description     string            `json:"description,omitempty"`
	ID              string            `json:"id,omitempty"`
	Name            string            `json:"name,omitempty"`
	Protocol        string            `json:"protocol,omitempty"`
	ProtocolMappers []*ProtocolMapper `json:"protocolMappers,omitempty"`
}

type ProtocolMapper struct {
	Configuration  map[string]string `json:"config,omitempty"`
	ID             string            `json:"id,omitempty"`
	Name           string            `json:"name,omitempty"`
	Protocol       string            `json:"protocol,omitempty"`
	ProtocolMapper string            `json:"protocolMapper,omitempty"`
}

type ProtocolMapperAttribute string

const (
	TokenClaimName                      ProtocolMapperAttribute = "claim.name"
	TokenClaimNameLabel                 ProtocolMapperAttribute = "tokenClaimName.label"
	TokenClaimNameTooltip               ProtocolMapperAttribute = "tokenClaimName.tooltip"
	JsonType                            ProtocolMapperAttribute = "jsonType.label"
	JsonTypeTooltip                     ProtocolMapperAttribute = "jsonType.tooltip"
	IncludeInAccessToken                ProtocolMapperAttribute = "access.token.claim"
	IncludeInAccessTokenLabel           ProtocolMapperAttribute = "includeInAccessToken.label"
	IncludeInAccessTokenTooltip         ProtocolMapperAttribute = "includeInAccessTokenTooltip.label"
	IncludeInIDToken                    ProtocolMapperAttribute = "id.token.claim"
	IncludeInIDTokenLabel               ProtocolMapperAttribute = "includeInIdToken.label"
	IncludeInIDTokenTooltip             ProtocolMapperAttribute = "includeInIdToken.tooltip"
	IncludeInAccessTokenResponse        ProtocolMapperAttribute = "includeInAccessTokenResponse.label"
	IncludeInAccessTokenResponseTooltip ProtocolMapperAttribute = "includeInAccessTokenResponse.tooltip"
	UserAttribute                       ProtocolMapperAttribute = "user.attribute"
	IncludeInUserInfo                   ProtocolMapperAttribute = "userinfo.token.claim"
	IncludeInUserInfoLabel              ProtocolMapperAttribute = "includeInUserInfo.label"
	IncludeInUserInfoTooltip            ProtocolMapperAttribute = "includeInUserInfo.tooltip"
)

type ProtocolMapperValueType = string

const (
	JsonTypeBoolean ProtocolMapperValueType = "boolean"
	JsonTypeString  ProtocolMapperValueType = "string"
	JsonTypeLong    ProtocolMapperValueType = "long"
	JsonTypeInt     ProtocolMapperValueType = "int"
	JsonTypeJson    ProtocolMapperValueType = "JSON"
)

func (p *ProtocolMapper) GetAttribute(attribute ProtocolMapperAttribute) (string, bool) {
	if p.Configuration == nil {
		return "", false
	}
	value, ok := p.Configuration[string(attribute)]
	return value, ok
}

func (p *ProtocolMapper) HasAttribute(attribute ProtocolMapperAttribute) bool {
	_, ok := p.Configuration[string(attribute)]
	return ok
}

func (p *ProtocolMapper) SetAttribute(attribute ProtocolMapperAttribute, value string) {
	if p.Configuration == nil {
		p.Configuration = map[string]string{}
	}
	p.Configuration[string(attribute)] = value
}

func (p *ProtocolMapper) GetJSONType() (*ProtocolMapperValueType, bool) {
	if value, ok := p.GetAttribute(JsonType); ok {
		return &value, true
	}
	return nil, false
}

func (p *ProtocolMapper) SetJSONType(valueType ProtocolMapperValueType) {
	p.SetAttribute(JsonType, valueType)
}

func (p *ProtocolMapper) SetClaimName(claimName string) {
	p.SetAttribute(TokenClaimName, claimName)
}

func (p *ProtocolMapper) GetClaimName() string {
	if value, ok := p.GetAttribute(TokenClaimName); ok {
		return value
	}
	return ""
}
func (p *ProtocolMapper) SetUserAttribute(userAttribute string) {
	p.SetAttribute(UserAttribute, userAttribute)
}

func (p *ProtocolMapper) GetUserAttribute() string {
	if value, ok := p.GetAttribute(UserAttribute); ok {
		return value
	}
	return ""
}

type ProtocolMapperProtocol string

const (
	OpenIDProtocol ProtocolMapperProtocol = "openid-connect"
)

func (p *ProtocolMapper) SetProtocol(protocol ProtocolMapperProtocol) {
	p.Protocol = string(protocol)
}

func (p *ProtocolMapper) GetProtocol() ProtocolMapperProtocol {
	return ProtocolMapperProtocol(p.Protocol)
}

type User struct {
	Access                     map[string]bool     `json:"access,omitempty"`
	Attributes                 map[string][]string `json:"attributes,omitempty"`
	ClientConsents             []UserConsent       `json:"clientConsents,omitempty"`
	ClientRoles                map[string][]string `json:"clientRoles,omitempty"`
	CreatedTimestamp           int64               `json:"createdTimestamp,omitempty"`
	Credentials                []Credential        `json:"credentials,omitempty"`
	DisableableCredentialTypes []string            `json:"disableableCredentialTypes,omitempty"`
	Email                      string              `json:"email,omitempty"`
	EmailVerified              bool                `json:"emailVerified,omitempty"`
	Enabled                    bool                `json:"enabled,omitempty"`
	FederatedIdentities        []FederatedIdentity `json:"federatedIdentities,omitempty"`
	FederationLink             string              `json:"federationLink,omitempty"`
	FirstName                  string              `json:"firstName,omitempty"`
	Groups                     []string            `json:"groups,omitempty"`
	ID                         string              `json:"id,omitempty"`
	LastName                   string              `json:"lastName,omitempty"`
	NotBefore                  int32               `json:"notBefore,omitempty"`
	Origin                     string              `json:"origin,omitempty"`
	RealmRoles                 []string            `json:"realmRoles,omitempty"`
	RequiredActions            []string            `json:"requiredActions,omitempty"`
	Self                       string              `json:"self,omitempty"`
	ServiceAccountClientID     string              `json:"serviceAccountClientId,omitempty"`
	Username                   string              `json:"username,omitempty"`
}

type FederatedIdentity struct {
	IdentityProvider string `json:"identityProvider,omitempty"`
	UserID           string `json:"userId,omitempty"`
	Username         string `json:"userName,omitempty"`
}

type Credential struct {
	CreatedDate    int64  `json:"createdDate,omitempty"`
	CredentialData string `json:"credentialData,omitempty"`
	ID             string `json:"id,omitempty"`
	Priority       string `json:"priority,omitempty"`
	SecretData     string `json:"secretData,omitempty"`
	Temporary      bool   `json:"temporary,omitempty"`
	Type           string `json:"type,omitempty"`
	UserLabel      string `json:"userLabel,omitempty"`
	Value          string `json:"value,omitempty"`
}

type UserConsent struct {
	ClientID            string   `json:"clientId,omitempty"`
	CreatedDate         int64    `json:"createdDate,omitempty"`
	GrantedClientScopes []string `json:"grantedClientScopes,omitempty"`
	LastUpdateDate      int64    `json:"lastUpdateDate,omitempty"`
}
