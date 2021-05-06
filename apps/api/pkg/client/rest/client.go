package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

var (
	// NameMayNotBe specifies strings that cannot be used as names specified as path segments
	NameMayNotBe = []string{".", ".."}

	// NameMayNotContain specifies substrings that cannot be used in names specified as path segments
	NameMayNotContain = []string{"/", "%"}
)

type Interface interface {
	Verb(verb string) *Request
	Post() *Request
	Put() *Request
	Get() *Request
	Delete() *Request
}

type RESTClient struct {
	base             *url.URL
	client           *http.Client
	content          ContentConfig
	versionedAPIPath string
}

type ContentConfig struct {
	ContentType       string
	AcceptContentType string
	Group             string
	Version           string
}

var DefaultContentConfig = ContentConfig{
	ContentType:       "application/json",
	AcceptContentType: "application/json",
	Group:             "",
	Version:           "",
}

type Config struct {
	ContentConfig
	APIPath string
	Host    string
}

func NewRESTClient(baseURL *url.URL, versionedAPIPath string, config ContentConfig, httpClient *http.Client) *RESTClient {
	if len(config.ContentType) == 0 {
		config.ContentType = "application/json"
	}
	base := *baseURL
	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}
	base.RawQuery = ""
	base.Fragment = ""

	return &RESTClient{
		base:             &base,
		versionedAPIPath: versionedAPIPath,
		client:           httpClient,
		content:          config,
	}
}

func RESTClientFor(config *Config) (*RESTClient, error) {

	if len(config.Group) == 0 {
		return nil, fmt.Errorf("group is required")
	}
	if len(config.Version) == 0 {
		return nil, fmt.Errorf("group is required")
	}

	versionedAPIPath := path.Join(config.APIPath, config.Group, config.Version)
	hostURL, err := url.Parse(config.Host)
	if err != nil {
		return nil, err
	}

	httpClient := http.DefaultClient

	restClient := NewRESTClient(hostURL, versionedAPIPath, config.ContentConfig, httpClient)
	return restClient, nil

}

func (c *RESTClient) Verb(verb string) *Request {
	return NewRequest(c).Verb(verb)
}

func (c *RESTClient) Post() *Request {
	return c.Verb("POST")
}

func (c *RESTClient) Put() *Request {
	return c.Verb("PUT")
}

func (c *RESTClient) Delete() *Request {
	return c.Verb("DELETE")
}

func (c *RESTClient) Get() *Request {
	return c.Verb("GET")
}

type Request struct {
	c            *RESTClient
	pathPrefix   string
	timeout      time.Duration
	headers      http.Header
	verb         string
	err          error
	subpath      string
	resource     string
	resourceName string
	params       url.Values
	body         io.Reader
}

func NewRequest(c *RESTClient) *Request {
	var pathPrefix string
	if c.base != nil {
		pathPrefix = path.Join("/", c.base.Path, c.versionedAPIPath)
	} else {
		pathPrefix = path.Join("/", c.versionedAPIPath)
	}

	var timeout time.Duration
	if c.client != nil {
		timeout = c.client.Timeout
	}

	r := &Request{
		c:          c,
		pathPrefix: pathPrefix,
		timeout:    timeout,
	}

	switch {
	case len(c.content.AcceptContentType) > 0:
		r.SetHeader("Accept", c.content.AcceptContentType)
	case len(c.content.ContentType) > 0:
		r.SetHeader("Accept", c.content.ContentType+", */*")
	}

	switch {
	case len(c.content.ContentType) > 0:
		r.SetHeader("Content-Type", c.content.ContentType)
	default:
		r.SetHeader("Content-Type", "application/json")
	}

	return r

}

func NewRequestWithClient(base *url.URL, versionedAPIPath string, config ContentConfig, client *http.Client) *Request {
	return NewRequest(&RESTClient{
		base:             base,
		versionedAPIPath: versionedAPIPath,
		content:          config,
		client:           client,
	})
}

func (r *Request) Verb(verb string) *Request {
	r.verb = verb
	return r
}

func (r *Request) Prefix(segments ...string) *Request {
	if r.err != nil {
		return r
	}
	r.pathPrefix = path.Join(r.pathPrefix, path.Join(segments...))
	return r
}

func (r *Request) Suffix(segments ...string) *Request {
	if r.err != nil {
		return r
	}
	r.subpath = path.Join(r.subpath, path.Join(segments...))
	return r
}

func (r *Request) Resource(resource string) *Request {
	if r.err != nil {
		return r
	}
	if len(r.resource) != 0 {
		r.err = fmt.Errorf("resource already set to %s, cannot change to %s", r.resource, resource)
		return r
	}
	if msgs := IsValidPathSegmentName(resource); len(msgs) != 0 {
		r.err = fmt.Errorf("invalid resource %q: %v", resource, msgs)
		return r
	}
	r.resource = resource
	return r
}

func (r *Request) Name(resourceName string) *Request {
	if r.err != nil {
		return r
	}
	if len(resourceName) == 0 {
		r.err = fmt.Errorf("resource name may not be empty")
		return r
	}
	if len(r.resourceName) > 0 {
		r.err = fmt.Errorf("resource name already set to %s, cannot change to %s", r.resourceName, resourceName)
		return r
	}
	if msgs := IsValidPathSegmentName(resourceName); len(msgs) != 0 {
		r.err = fmt.Errorf("invalid resource name %q: %v", resourceName, msgs)
		return r
	}
	r.resourceName = resourceName
	return r
}

func (r *Request) Param(paramName, s string) *Request {
	if r.err != nil {
		return r
	}
	if r.params == nil {
		r.params = make(url.Values)
	}
	r.params[paramName] = append(r.params[paramName], s)
	return r
}

func (r *Request) SetHeader(key string, values ...string) *Request {
	if r.headers == nil {
		r.headers = http.Header{}
	}
	r.headers.Del(key)
	for _, value := range values {
		r.headers.Add(key, value)
	}
	return r
}

func (r *Request) Body(obj interface{}) *Request {
	if r.err != nil {
		return r
	}
	switch t := obj.(type) {
	case string:
		data, err := ioutil.ReadFile(t)
		if err != nil {
			r.err = err
			return r
		}
		r.body = bytes.NewReader(data)
	case []byte:
		r.body = bytes.NewReader(t)
	case io.Reader:
		r.body = t
	default:
		bodyBytes, err := json.Marshal(obj)
		if err != nil {
			r.err = err
			return r
		}
		r.body = bytes.NewReader(bodyBytes)
	}
	return r
}

func (r *Request) URL() *url.URL {
	p := r.pathPrefix
	if len(r.resource) > 0 {
		p = path.Join(p, strings.ToLower(r.resource))
	}
	if len(r.resourceName) != 0 {
		p = path.Join(p, r.resourceName)
	}
	finalUrl := &url.URL{}
	if r.c.base != nil {
		*finalUrl = *r.c.base
	}
	finalUrl.Path = p

	query := url.Values{}
	for key, values := range r.params {
		for _, value := range values {
			query.Add(key, value)
		}
	}

	if r.timeout != 0 {
		query.Set("timeout", r.timeout.String())
	}

	finalUrl.RawQuery = query.Encode()
	return finalUrl

}

func (r *Request) request(ctx context.Context, fn func(r *http.Request, response *http.Response)) error {

	if r.err != nil {
		logrus.Errorf("error in request: %v", r.err)
		return r.err
	}

	client := r.c.client
	if client == nil {
		client = http.DefaultClient
	}

	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
		defer cancel()
	}

	reqURL := r.URL().String()
	req, err := http.NewRequest(r.verb, reqURL, r.body)
	if err != nil {
		return err
	}
	req.WithContext(ctx)
	req.Header = r.headers

	logrus.Infof("%s request to %s", req.Method, req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	fn(req, resp)
	return nil

}

func (r *Request) Do(ctx context.Context) Result {
	var result Result
	err := r.request(ctx, func(req *http.Request, res *http.Response) {
		result = r.transformResponse(req, res)
	})
	if err != nil {
		return Result{err: err}
	}
	return result
}

func (r *Request) transformResponse(req *http.Request, res *http.Response) Result {
	var body []byte
	if res.Body != nil {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return Result{err: err}
		}
		body = data
	}
	contentType := res.Header.Get("Content-Type")
	if len(contentType) == 0 {
		contentType = r.c.content.ContentType
	}

	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusPartialContent {
		return Result{
			err: fmt.Errorf("unexpected server response with status code %d and content type %v", res.StatusCode, res.Header.Get("Content-Type")),
		}
	}

	return Result{
		body:        body,
		contentType: contentType,
		statusCode:  res.StatusCode,
	}

}

type Result struct {
	body        []byte
	contentType string
	err         error
	statusCode  int
}

func (r Result) Raw() ([]byte, error) {
	return r.body, r.err
}

func (r Result) Error() error {
	return r.err
}

func (r Result) StatusCode() int {
	return r.statusCode
}

func IsValidPathSegmentName(name string) []string {
	var errors []string
	for _, illegalName := range NameMayNotBe {
		if name == illegalName {
			errors = append(errors, fmt.Sprintf("name may not be '%s'", illegalName))
		}
	}

	for _, illegalContent := range NameMayNotContain {
		if strings.Contains(name, illegalContent) {
			errors = append(errors, fmt.Sprintf("name may not container '%s'", illegalContent))
		}
	}
	return errors
}

func (r Result) Into(obj interface{}) error {
	if r.err != nil {
		return r.Error()
	}
	if len(r.body) == 0 {
		return fmt.Errorf("0-length response with status code: %d and content-type %s", r.statusCode, r.contentType)
	}
	if err := json.Unmarshal(r.body, obj); err != nil {
		return err
	}
	return nil
}
