package rest

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/watch"
	"mime"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Request struct {
	c            *RESTClient
	verb         string
	pathPrefix   string
	timeout      time.Duration
	headers      http.Header
	params       url.Values
	resource     string
	resourceName string
	err          error
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
	if c.Client != nil {
		timeout = c.Client.Timeout
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

	return r
}

func (r *Request) Verb(verb string) *Request {
	r.verb = verb
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

func (r *Request) Resource(resource string) *Request {
	if r.err != nil {
		return r
	}
	if len(r.resource) > 0 {
		r.err = fmt.Errorf("resource already set to %q, cannot change to %q", r.resource, resource)
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
	if len(r.resourceName) != 0 {
		r.err = fmt.Errorf("resource name already set to %q, cannot change to %q", r.resourceName, resourceName)
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
	return r.setParam(paramName, s)
}

func (r *Request) VersionedParams(obj runtime.Object, codec runtime.ParameterCodec) *Request {
	return r.SpecificallyVersionedParam(obj, codec, r.c.content.GroupVersion)
}

func (r *Request) SpecificallyVersionedParam(obj runtime.Object, codec runtime.ParameterCodec, version schema.GroupVersion) *Request {
	if r.err != nil {
		return r
	}
	params, err := codec.EncodeParameters(obj, version)
	if err != nil {
		r.err = err
		return r
	}
	for k, v := range params {
		for _, value := range v {
			r.setParam(k, value)
		}
	}
	return r
}

func (r *Request) setParam(paramName, value string) *Request {
	if r.params == nil {
		r.params = make(url.Values)
	}
	r.params[paramName] = append(r.params[paramName], value)
	return r
}

func (r *Request) Timeout(d time.Duration) *Request {
	if r.err != nil {
		return r
	}
	r.timeout = d
	return r
}

func (r *Request) Body(body interface{}) *Request {
	if r.err != nil {
		return r
	}
	switch t := body.(type) {
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
		r.SetHeader("Content-Type", r.c.content.ContentType)
	default:
		r.err = fmt.Errorf("unknown type for body: %+v", body)
	}
	return r
}

func (r *Request) URL() *url.URL {
	p := r.pathPrefix
	if len(r.resource) != 0 {
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

	url := r.URL().String()

	url = strings.Replace(url, "https://", "wss://", -1)
	url = strings.Replace(url, "http://", "ws://", -1)

	req, err := http.NewRequest(r.verb, url, r.body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header = r.headers
	client := r.c.Client
	if client == nil {
		client = http.DefaultClient
	}
	conn, res, err := websocket.DefaultDialer.DialContext(ctx, url, r.headers)
	if err != nil {
		return nil, Result{err: err}.Error()
	}
	if res.StatusCode != http.StatusSwitchingProtocols {
		if result := r.transformResponse(res, r.verb); result.err != nil {
			return nil, result.Error()
		}
		return nil, fmt.Errorf("for request %s, got status: %v", url, res.StatusCode)
	}

	decoder, err := r.c.content.Negotiator.Decoder(res.Header.Get("Content-Type"), nil)
	if err != nil {
		return nil, err
	}

	ws := NewWebSocketWatcher(conn, decoder)
	return ws, nil
}

type WebSocketWatcher struct {
	conn    *websocket.Conn
	decoder runtime.Decoder
	result  chan watch.Event
	done    chan struct{}
	lock    sync.Mutex
}

func NewWebSocketWatcher(conn *websocket.Conn, decoder runtime.Decoder) *WebSocketWatcher {
	ws := &WebSocketWatcher{
		conn:    conn,
		decoder: decoder,
		result:  make(chan watch.Event),
		done:    make(chan struct{}),
	}
	go ws.receive()
	return ws
}

func (ws *WebSocketWatcher) ResultChan() <-chan watch.Event {
	return ws.result
}

func (ws *WebSocketWatcher) Stop() {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	select {
	case <-ws.done:
	default:
		close(ws.done)
		ws.conn.Close()
	}
}

func (ws *WebSocketWatcher) receive() {
	defer close(ws.result)
	defer ws.Stop()
	for {
		_, reader, err := ws.conn.NextReader()
		if err != nil {
			switch err {
			case io.EOF:
			case io.ErrUnexpectedEOF:
				logrus.Warnf("unexpected EOF during watch: %v", err)
			default:
				if net.IsProbableEOF(err) || net.IsTimeout(err) {
					logrus.Errorf("unable to decode event from the watch stream: %v", err)
				} else {
					select {
					case <-ws.done:
					case ws.result <- watch.Event{
						Type:   watch.Error,
						Object: &errors.NewInternalError(err).ErrStatus,
					}:
					}
				}
			}
			return
		}
		var buf []byte
		if _, err := reader.Read(buf); err != nil {
			logrus.Errorf("unable to read: %v", err)
			return
		}
		we := metav1.WatchEvent{}
		decoded, _, err := ws.decoder.Decode(buf, nil, &we)
		if err != nil {
			logrus.Errorf("unable to decode event: %v", err)
			return
		}
		if decoded != &we {
			logrus.Errorf("unable to decode to metav1.Event")
			return
		}
		obj, _, err := ws.decoder.Decode(we.Object.Raw, nil, nil)
		if err != nil {
			logrus.Errorf("unable to decode object: %v", err)
			return
		}
		ev := watch.Event{
			Type:   watch.EventType(we.Type),
			Object: obj,
		}
		select {
		case <-ws.done:
		case ws.result <- ev:
		}
	}
}

func (r *Request) Do(ctx context.Context) Result {
	var result Result
	err := r.request(ctx, func(req *http.Request, resp *http.Response) {
		result = r.transformResponse(resp, req.Method)
	})
	if err != nil {
		return Result{err: err}
	}
	return result
}

func (r *Request) request(ctx context.Context, fn func(*http.Request, *http.Response)) error {
	if r.err != nil {
		return r.err
	}

	client := r.c.Client
	if client == nil {
		client = http.DefaultClient
	}

	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
		defer cancel()
	}

	reqUrl := r.URL().String()
	req, err := http.NewRequest(r.verb, reqUrl, r.body)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	req.Header = r.headers

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	fn(req, resp)
	return nil
}

func (r *Request) transformResponse(resp *http.Response, method string) Result {
	var body []byte
	if resp.Body != nil {
		data, err := ioutil.ReadAll(resp.Body)
		switch err.(type) {
		case nil:
			body = data
		default:
			return Result{
				err: fmt.Errorf("unexpected error while reading response body: %v", err),
			}
		}
	}

	var decoder runtime.Decoder
	contentType := resp.Header.Get("Content-Type")
	if len(contentType) == 0 {
		contentType = r.c.content.ContentType
	}
	if len(contentType) > 0 {
		var err error
		mediaType, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			return Result{err: errors.NewInternalError(err)}
		}
		decoder, err = r.c.content.Negotiator.Decoder(mediaType, params)
		if err != nil {
			switch {
			case resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusPartialContent:
				return Result{err: r.transformUnstructuredResponseError(resp, method, body)}
			}
			return Result{
				body:        body,
				contentType: contentType,
				statusCode:  resp.StatusCode,
			}
		}
	}

	switch {
	case resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusPartialContent:
		retryAfter, _ := retryAfterSeconds(resp)
		err := r.newUnstructuredResponseError(body, isTextResponse(resp), resp.StatusCode, method, retryAfter)
		return Result{
			body:        body,
			contentType: contentType,
			statusCode:  resp.StatusCode,
			decoder:     decoder,
			err:         err,
		}
	}

	return Result{
		body:        body,
		contentType: contentType,
		statusCode:  resp.StatusCode,
		decoder:     decoder,
	}

}

type Result struct {
	err         error
	body        []byte
	contentType string
	statusCode  int
	decoder     runtime.Decoder
}

func (r Result) Into(obj runtime.Object) error {
	if r.err != nil {
		return r.Error()
	}
	if r.decoder == nil {
		return fmt.Errorf("serializer for %s doesn't exist", r.contentType)
	}
	if len(r.body) == 0 {
		return fmt.Errorf("0-length response with status code: %d and content type: %s", r.statusCode, r.contentType)
	}

	out, _, err := r.decoder.Decode(r.body, nil, obj)
	if err != nil {
		return err
	}

	switch t := out.(type) {
	case *metav1.Status:
		if t.Status != metav1.StatusSuccess {
			return errors.FromObject(t)
		}
	}
	return nil
}

func (r Result) Error() error {
	if r.err == nil || !errors.IsUnexpectedServerError(r.err) || len(r.body) == 0 || r.decoder == nil {
		return r.err
	}
	out, _, err := r.decoder.Decode(r.body, &schema.GroupVersionKind{Version: "v1"}, nil)
	if err != nil {
		return r.err
	}
	switch t := out.(type) {
	case *metav1.Status:
		if t.Status == metav1.StatusFailure {
			return errors.FromObject(t)
		}
	}
	return r.err
}

// maxUnstructuredResponseTextBytes is an upper bound on how much output to include in the unstructured error.
const maxUnstructuredResponseTextBytes = 2048

func (r *Request) transformUnstructuredResponseError(resp *http.Response, method string, body []byte) error {
	if body == nil && resp.Body != nil {
		if data, err := ioutil.ReadAll(&io.LimitedReader{R: resp.Body, N: maxUnstructuredResponseTextBytes}); err == nil {
			body = data
		}
	}
	retryAfter, _ := retryAfterSeconds(resp)
	return r.newUnstructuredResponseError(body, isTextResponse(resp), resp.StatusCode, method, retryAfter)
}

func (r *Request) newUnstructuredResponseError(body []byte, isTextResponse bool, statusCode int, method string, retryAfter int) error {
	message := "unknown"
	if isTextResponse {
		message = strings.TrimSpace(string(body))
	}
	var groupResource schema.GroupResource
	if len(r.resource) > 0 {
		groupResource.Group = r.c.content.GroupVersion.Group
		groupResource.Resource = r.resource
	}
	return errors.NewGenericServerResponse(
		statusCode,
		method,
		groupResource,
		r.resourceName,
		message,
		retryAfter,
		true,
	)
}

func retryAfterSeconds(resp *http.Response) (int, bool) {
	if h := resp.Header.Get("Retry-After"); len(h) > 0 {
		if i, err := strconv.Atoi(h); err == nil {
			return i, true
		}
	}
	return 0, false
}

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

// NameMayNotBe specifies strings that cannot be used as names specified as path segments (like the REST API or etcd store)
var NameMayNotBe = []string{".", ".."}

// NameMayNotContain specifies substrings that cannot be used in names specified as path segments (like the REST API or etcd store)
var NameMayNotContain = []string{"/", "%"}

// IsValidPathSegmentName validates the name can be safely encoded as a path segment
func IsValidPathSegmentName(name string) []string {
	for _, illegalName := range NameMayNotBe {
		if name == illegalName {
			return []string{fmt.Sprintf(`may not be '%s'`, illegalName)}
		}
	}

	var errors []string
	for _, illegalContent := range NameMayNotContain {
		if strings.Contains(name, illegalContent) {
			errors = append(errors, fmt.Sprintf(`may not contain '%s'`, illegalContent))
		}
	}

	return errors
}

// IsValidPathSegmentPrefix validates the name can be used as a prefix for a name which will be encoded as a path segment
// It does not check for exact matches with disallowed names, since an arbitrary suffix might make the name valid
func IsValidPathSegmentPrefix(name string) []string {
	var errors []string
	for _, illegalContent := range NameMayNotContain {
		if strings.Contains(name, illegalContent) {
			errors = append(errors, fmt.Sprintf(`may not contain '%s'`, illegalContent))
		}
	}

	return errors
}

// ValidatePathSegmentName validates the name can be safely encoded as a path segment
func ValidatePathSegmentName(name string, prefix bool) []string {
	if prefix {
		return IsValidPathSegmentPrefix(name)
	}
	return IsValidPathSegmentName(name)
}
