package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/nrc-no/core/apps/api/pkg/api/patch"
	v1 "github.com/nrc-no/core/apps/api/pkg/apis/core/v1"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	watch2 "github.com/nrc-no/core/apps/api/pkg/client/rest/watch"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/util/exceptions"
	"github.com/nrc-no/core/apps/api/pkg/watch"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/types"
	"mime"
	"net/http"
	"net/url"
	"path"
	"reflect"
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
	Patch(pt patch.PatchType) *Request
}

type RESTClient struct {
	base             *url.URL
	client           *http.Client
	content          ClientContentConfig
	versionedAPIPath string
}

type ContentConfig struct {
	ContentType          string
	AcceptContentType    string
	GroupVersion         *schema.GroupVersion
	NegotiatedSerializer runtime.NegotiatedSerializer
}

var DefaultContentConfig = ContentConfig{
	ContentType:       "application/json",
	AcceptContentType: "application/json",
}

type Config struct {
	APIPath string
	Host    string
	ContentConfig
}

type ClientContentConfig struct {
	AcceptContentTypes string
	ContentType        string
	GroupVersion       schema.GroupVersion
	Negotiator         runtime.ClientNegotiator
}

func NewRESTClient(
	baseURL *url.URL,
	versionedAPIPath string,
	config ClientContentConfig,
	httpClient *http.Client,
) *RESTClient {
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

	if len(config.GroupVersion.Group) == 0 {
		return nil, fmt.Errorf("group is required")
	}
	if len(config.GroupVersion.Version) == 0 {
		return nil, fmt.Errorf("group is required")
	}
	var gv schema.GroupVersion
	if config.GroupVersion != nil {
		gv = *config.GroupVersion
	}

	versionedAPIPath := path.Join(config.APIPath, config.GroupVersion.Group, config.GroupVersion.Version)
	hostURL, err := url.Parse(config.Host)
	if err != nil {
		return nil, err
	}

	httpClient := http.DefaultClient

	clientContent := ClientContentConfig{
		AcceptContentTypes: config.AcceptContentType,
		ContentType:        config.ContentType,
		GroupVersion:       gv,
		Negotiator:         runtime.NewClientNegotiator(config.NegotiatedSerializer, gv),
	}

	restClient := NewRESTClient(
		hostURL,
		versionedAPIPath,
		clientContent,
		httpClient)

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

func (c *RESTClient) Patch(pt patch.PatchType) *Request {
	return c.Verb("PATCH")
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
	namespaceSet bool
	namespace    string
	subresource  string
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
	case len(c.content.AcceptContentTypes) > 0:
		r.SetHeader("Accept", c.content.AcceptContentTypes)
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

func NewRequestWithClient(base *url.URL, versionedAPIPath string, config ClientContentConfig, client *http.Client) *Request {
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

func (r *Request) Namespace(namespace string) *Request {
	if r.err != nil {
		return r
	}
	if r.namespaceSet {
		r.err = fmt.Errorf("namespace already set to %q, cannot change to %q", r.namespace, namespace)
		return r
	}
	r.namespaceSet = true
	r.namespace = namespace
	return r
}

// SubResource sets a sub-resource path which can be multiple segments after the resource
// name but before the suffix.
func (r *Request) SubResource(subresources ...string) *Request {
	if r.err != nil {
		return r
	}
	subresource := path.Join(subresources...)
	if len(r.subresource) != 0 {
		r.err = fmt.Errorf("subresource already set to %q, cannot change to %q", r.subresource, subresource)
		return r
	}
	for _, s := range subresources {
		if msgs := IsValidPathSegmentName(s); len(msgs) != 0 {
			r.err = fmt.Errorf("invalid subresource %q: %v", s, msgs)
			return r
		}
	}
	r.subresource = subresource
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

func (r *Request) VersionedParams(obj runtime.Object, codec runtime.ParameterCodec) *Request {
	return r.SpecificallyVersionedParams(obj, codec, r.c.content.GroupVersion)
}

func (r *Request) SpecificallyVersionedParams(obj runtime.Object, codec runtime.ParameterCodec, version schema.GroupVersion) *Request {
	if r.err != nil {
		return r
	}
	params, err := codec.EncodeParameters(obj, version)
	if err != nil {
		r.err = err
		return r
	}
	for k, v := range params {
		if r.params == nil {
			r.params = make(url.Values)
		}
		r.params[k] = append(r.params[k], v...)
	}
	return r
}

// Timeout makes the request use the given duration as an overall timeout for the
// request. Additionally, if set passes the value as "timeout" parameter in URL.
func (r *Request) Timeout(d time.Duration) *Request {
	if r.err != nil {
		return r
	}
	r.timeout = d
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
	case runtime.Object:
		if reflect.ValueOf(t).IsNil() {
			return r
		}
		encoder, err := r.c.content.Negotiator.Encoder(r.c.content.ContentType, nil)
		if err != nil {
			r.err = err
			return r
		}
		data, err := runtime.Encode(encoder, t)
		if err != nil {
			r.err = err
			return r
		}
		r.body = bytes.NewReader(data)
		r.SetHeader("Content-Type", "application/json")
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

func (r *Request) Watch(ctx context.Context) (watch.Interface, error) {
	if r.err != nil {
		return nil, r.err
	}

	reqURL := r.URL().String()

	reqURL = strings.Replace(reqURL, "https://", "wss://", -1)
	reqURL = strings.Replace(reqURL, "http://", "ws://", -1)

	c, resp, err := websocket.DefaultDialer.DialContext(ctx, reqURL, r.headers)
	if err != nil {
		return nil, err
	}

	contentType := resp.Header.Get("Content-Type")
	decoder, err := r.c.content.Negotiator.Decoder(contentType, nil)
	if err != nil {
		r.err = err
		return nil, err
	}

	ctx, stop := context.WithCancel(ctx)

	watcher := watch2.NewWebSocketWatcher(ctx, stop, decoder, func() ([]byte, error) {
		_, message, err := c.ReadMessage()
		return message, err
	})

	return watcher, nil

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

	var decoder runtime.Decoder
	contentType := res.Header.Get("Content-Type")
	if len(contentType) == 0 {
		contentType = r.c.content.ContentType
	}
	if len(contentType) > 0 {
		var err error
		mediaType, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			return Result{err: exceptions.NewInternalError(err)}
		}
		decoder, err = r.c.content.Negotiator.Decoder(mediaType, params)
		if err != nil {
			switch {
			case res.StatusCode < http.StatusOK || res.StatusCode > http.StatusPartialContent:
				return Result{err: r.transformUnstructuredResponseError(res, req, body)}
			}
			return Result{
				body:        body,
				contentType: contentType,
				statusCode:  res.StatusCode,
			}
		}
	}

	switch {
	case res.StatusCode < http.StatusOK || res.StatusCode > http.StatusPartialContent:
		err := r.newUnstructuredResponseError(body, isTextResponse(res), res.StatusCode, req.Method)
		return Result{
			body:        body,
			contentType: contentType,
			statusCode:  res.StatusCode,
			decoder:     decoder,
			err:         err,
		}
	}

	decoder, err := r.c.content.Negotiator.Decoder(contentType, nil)
	if err != nil {
		return Result{
			body:        body,
			contentType: contentType,
			err:         fmt.Errorf("unable to get decoder for response: %v", err),
			statusCode:  res.StatusCode,
			decoder:     decoder,
		}
	}

	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusPartialContent {
		return Result{
			body:        body,
			contentType: contentType,
			statusCode:  res.StatusCode,
			decoder:     decoder,
			err:         fmt.Errorf("unexpected server response with status code %d and content type %v", res.StatusCode, res.Header.Get("Content-Type")),
		}
	}

	return Result{
		body:        body,
		contentType: contentType,
		statusCode:  res.StatusCode,
		decoder:     decoder,
	}

}

type Result struct {
	body        []byte
	contentType string
	err         error
	statusCode  int
	decoder     runtime.Decoder
}

func (r Result) Raw() ([]byte, error) {
	return r.body, r.err
}

func (r Result) Error() error {

	if r.err == nil || len(r.body) == 0 || r.decoder == nil {
		return r.err
	}

	gvk := v1.SchemeGroupVersion.WithKind("")
	out, _, err := r.decoder.Decode(r.body, &gvk, nil)
	if err != nil {
		logrus.Infof("body was not decodable (unable to check for Status): %v", err)
		return r.err
	}

	switch t := out.(type) {
	case *metav1.Status:
		if t.Status == metav1.StatusFailure {
			return exceptions.FromObject(t)
		}
	}

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

func (r Result) Into(obj runtime.Object) error {
	if r.err != nil {
		err := r.Error()

		return err
	}

	if len(r.body) == 0 {
		return fmt.Errorf("0-length response with status code: %d and content-type %s", r.statusCode, r.contentType)
	}

	out, _, err := r.decoder.Decode(r.body, nil, obj)
	if err != nil {
		return err
	}

	switch t := out.(type) {
	case *metav1.Status:
		if t.Status != metav1.StatusSuccess {
			return exceptions.FromObject(t)
		}
	}

	return nil
}

func (r Result) Get() (runtime.Object, error) {

	if r.err != nil {
		return nil, r.Error()
	}
	if r.decoder == nil {
		return nil, fmt.Errorf("serializer for %s doesn't exist", r.contentType)
	}

	out, _, err := r.decoder.Decode(r.body, nil, nil)
	if err != nil {
		return nil, err
	}
	switch t := out.(type) {
	case *metav1.Status:
		if t.Status != metav1.StatusSuccess {
			return nil, exceptions.FromObject(t)
		}
	}
	return out, nil

}

const maxUnstructuredResponseTextBytes = 2048

func (r *Request) transformUnstructuredResponseError(resp *http.Response, req *http.Request, body []byte) error {
	if body == nil && resp.Body != nil {
		if data, err := ioutil.ReadAll(&io.LimitedReader{R: resp.Body, N: maxUnstructuredResponseTextBytes}); err == nil {
			body = data
		}
	}
	return r.newUnstructuredResponseError(body, isTextResponse(resp), resp.StatusCode, req.Method)
}

func (r *Request) newUnstructuredResponseError(body []byte, isTextResponse bool, statusCode int, method string) error {
	// cap the amount of output we create
	if len(body) > maxUnstructuredResponseTextBytes {
		body = body[:maxUnstructuredResponseTextBytes]
	}

	message := "unknown"
	if isTextResponse {
		message = strings.TrimSpace(string(body))
	}
	var groupResource schema.GroupResource
	if len(r.resource) > 0 {
		groupResource.Group = r.c.content.GroupVersion.Group
		groupResource.Resource = r.resource
	}
	return exceptions.NewGenericServerResponse(
		statusCode,
		method,
		groupResource,
		r.resourceName,
		message,
		true,
	)
}

// isTextResponse returns true if the response appears to be a textual media type.
func isTextResponse(resp *http.Response) bool {
	contentType := resp.Header.Get("Content-Type")
	if len(contentType) == 0 {
		return true
	}
	media, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	return strings.HasPrefix(media, "text/")
}
