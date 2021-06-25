package iam

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/square/go-jose.v2/jwt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type RESTConfig struct {
	Scheme       string
	Host         string
	Token        string
	RefreshToken string
	IssuerURL    string
	ClientID     string
	ClientSecret string
	tokenLock    sync.Mutex
}

type RESTClient struct {
	config     *RESTConfig
	httpClient *http.Client
}

func (r *RESTClient) Verb(verb string) *Request {
	return r.NewRequest().Verb(verb)
}

func (r *RESTClient) Get() *Request {
	return r.Verb("GET")
}

func (r *RESTClient) Post() *Request {
	return r.Verb("POST")
}

func (r *RESTClient) Put() *Request {
	return r.Verb("PUT")
}

func (r *RESTClient) Delete() *Request {
	return r.Verb("DELETE")
}

func (r *RESTClient) NewRequest() *Request {
	return &Request{}
}

type Request struct {
	c       *RESTClient
	err     error
	verb    string
	path    string
	body    io.Reader
	params  url.Values
	headers http.Header
}

func (r *Request) Verb(verb string) *Request {
	if r.err != nil {
		return r
	}
	r.verb = verb
	return r
}

func (r *Request) Path(path string) *Request {
	if r.err != nil {
		return r
	}
	r.path = path
	return r
}

func (r *Request) Body(body interface{}) *Request {
	if r.err != nil {
		return r
	}
	switch b := body.(type) {
	case io.Reader:
		r.body = b
	case []byte:
		r.body = bytes.NewReader(b)
	case string:
		r.body = strings.NewReader(b)
	default:
		bs, err := json.Marshal(body)
		if err != nil {
			r.err = err
			return r
		}
		r.body = bytes.NewReader(bs)
	}
	return r
}

type UrlValuer interface {
	MarshalQueryParameters() (url.Values, error)
	UnmarshalQueryParameters(values url.Values) error
}

func (r *Request) WithParams(params interface{}) *Request {
	if r.err != nil {
		return r
	}

	switch p := params.(type) {
	case url.Values:
		r.params = url.Values{}
		for key, values := range p {
			for _, value := range values {
				r.params.Add(key, value)
			}
		}
	case UrlValuer:
		values, err := p.MarshalQueryParameters()
		if err != nil {
			r.err = err
			return r
		}
		return r.WithParams(values)
	default:
		r.err = fmt.Errorf("invalid params parameter")
		return r
	}

	return r
}

func (r *Request) WithHeader(key, value string) *Request {
	if r.err != nil {
		return r
	}
	if r.headers == nil {
		r.headers = http.Header{}
	}
	r.headers.Add(key, value)
	return r
}

func (r *Request) refreshToken() error {

	r.c.config.tokenLock.Lock()
	defer r.c.config.tokenLock.Unlock()

	if len(r.c.config.Token) == 0 {
		return nil
	}

	refreshToken := r.c.config.RefreshToken
	if len(refreshToken) == 0 {
		return nil
	}

	accessToken, err := jwt.ParseSigned(r.c.config.Token)
	if err != nil {
		return err
	}

	var claims jwt.Claims
	if err := accessToken.Claims(&claims); err != nil {
		return err
	}

	if err := claims.Validate(jwt.Expected{
		Time: time.Now(),
	}); err == nil {
		return nil
	}

	issuerURL := r.c.config.IssuerURL
	tokenURL := fmt.Sprintf("%s/oauth/token", issuerURL)

	payload := url.Values{}
	payload.Set("grant_type", "refresh_token")
	payload.Set("refresh_token", r.c.config.RefreshToken)
	payload.Set("client_id", r.c.config.ClientID)
	payload.Set("client_secret", r.c.config.ClientSecret)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(payload.Encode()))
	if err != nil {
		return err
	}

	httpClient := r.c.httpClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var tokenResponse struct {
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
	}

	if err := json.Unmarshal(bodyBytes, &tokenResponse); err != nil {
		return err
	}

	r.c.config.Token = tokenResponse.AccessToken
	if len(tokenResponse.RefreshToken) > 0 {
		r.c.config.RefreshToken = tokenResponse.AccessToken
	}

	return nil
}

func (r *Request) Do(ctx context.Context) *Response {

	if r.err != nil {
		return &Response{err: r.err}
	}

	if err := r.refreshToken(); err != nil {
		return &Response{err: err}
	}

	u := url.URL{}
	u.Scheme = r.c.config.Scheme
	u.Host = r.c.config.Host
	u.Path = r.path

	qry := u.Query()
	for k, values := range r.params {
		for _, value := range values {
			qry.Add(k, value)
		}
	}
	u.RawQuery = qry.Encode()

	req, err := http.NewRequestWithContext(ctx, r.verb, u.String(), r.body)
	if err != nil {
		return &Response{err: err}
	}

	for key, values := range r.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	if len(req.Header.Get("Accept")) == 0 {
		req.Header.Set("Accept", "application/json")
	}

	if len(req.Header.Get("Content-Type")) == 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	if len(req.Header.Get("Authorization")) == 0 && len(r.c.config.Token) > 0 {
		req.Header.Set("Authorization", "Bearer "+r.c.config.Token)
	}

	httpClient := r.c.httpClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return &Response{err: err}
	}

	if res.StatusCode < 200 || res.StatusCode > 399 {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return &Response{err: fmt.Errorf("unexpected status code: %d", res.StatusCode)}
		}
		return &Response{
			err:  fmt.Errorf("unexpected status code: %d. response: %s", res.StatusCode, string(bodyBytes)),
			body: bodyBytes,
		}
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &Response{err: err}
	}

	return &Response{
		body: bodyBytes,
	}

}

type Response struct {
	err  error
	body []byte
}

func (r *Response) Into(into interface{}) error {
	if r.err != nil {
		return r.err
	}
	if len(r.body) == 0 {
		return fmt.Errorf("0-length body")
	}
	if err := json.Unmarshal(r.body, &into); err != nil {
		return err
	}
	return nil
}

type ClientSet struct {
	c *RESTClient
}

func (c ClientSet) Parties() PartyClient {
	return &RESTPartyClient{
		c: c.c,
	}
}

func (c ClientSet) PartyTypes() PartyTypeClient {
	return &RESTPartyTypeClient{
		c: c.c,
	}
}

func (c ClientSet) Relationships() RelationshipClient {
	return &RESTRelationshipClient{
		c: c.c,
	}
}

func (c ClientSet) RelationshipTypes() RelationshipTypeClient {
	return &RESTRelationshipTypeClient{
		c: c.c,
	}
}

func (c ClientSet) Attributes() AttributeClient {
	return &RESTAttributeClient{
		c: c.c,
	}
}

func (c ClientSet) Teams() TeamClient {
	return &RESTTeamClient{
		c: c.c,
	}
}

func (c ClientSet) Organizations() OrganizationClient {
	return &RESTOrganizationClient{
		c: c.c,
	}
}

func (c ClientSet) Staff() StaffClient {
	return &RESTStaffClient{
		c: c.c,
	}
}

func (c ClientSet) Memberships() MembershipClient {
	return &RESTMembershipClient{
		c: c.c,
	}
}

func (c ClientSet) Individuals() IndividualClient {
	return &RESTIndividualClient{
		c: c.c,
	}
}

var _ Interface = &ClientSet{}
