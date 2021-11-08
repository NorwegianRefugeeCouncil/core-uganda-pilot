package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/mimetypes"
	"github.com/nrc-no/core/pkg/logging"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type Config struct {
	Scheme     string
	Host       string
	HTTPClient *http.Client
	Headers    http.Header
}

type Client struct {
	config *Config
}

func NewClient(config *Config) *Client {
	return &Client{
		config: config,
	}
}

func (r *Client) Verb(verb string) *Request {
	return r.NewRequest().Verb(verb)
}

func (r *Client) Get() *Request {
	return r.Verb(http.MethodGet)
}

func (r *Client) Post() *Request {
	return r.Verb(http.MethodPost)
}

func (r *Client) Put() *Request {
	return r.Verb(http.MethodPut)
}

func (r *Client) Delete() *Request {
	return r.Verb(http.MethodDelete)
}

func (r *Client) NewRequest() *Request {
	return &Request{
		c: r,
	}
}

type Request struct {
	c       *Client
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

type URLValuer interface {
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
	default:
		vp := reflect.New(reflect.TypeOf(p))
		vp.Elem().Set(reflect.ValueOf(p))
		intf := vp.Interface()

		valuer, ok := intf.(URLValuer)
		if ok {
			values, err := valuer.MarshalQueryParameters()
			if err != nil {
				r.err = err
				return r
			}
			return r.WithParams(values)
		}
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

type TokenIntrospectionResponse struct {
	Active bool `json:"active"`
}

func (r *Request) Do(ctx context.Context) *Response {

	l := logging.NewLogger(ctx).With(
		zap.String("scheme", r.c.config.Scheme),
		zap.String("host", r.c.config.Host),
		zap.String("path", r.path),
		zap.Any("params", r.params),
		zap.String("verb", r.verb))

	if r.err != nil {
		l.Error("could not send request", zap.Error(r.err))
		return &Response{err: r.err}
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
		l.Error("could not create http request", zap.Error(err))
		return &Response{err: err}
	}

	for key, values := range r.c.config.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	for key, values := range r.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	if len(req.Header.Get("Accept")) == 0 && req.Method != http.MethodDelete {
		req.Header.Set("Accept", mimetypes.ApplicationJson)
	}

	if len(req.Header.Get("Content-Type")) == 0 && r.body != nil {
		req.Header.Set("Content-Type", mimetypes.ApplicationJson)
	}

	l = l.With(zap.Any("headers", req.Header))

	httpClient := r.c.config.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	l.Debug("sending request")

	res, err := httpClient.Do(req)
	if err != nil {
		l.Error("could not send http request", zap.Error(err))
		return &Response{err: err}
	}

	l = l.With(zap.Int("status_code", res.StatusCode))

	if res.StatusCode < 200 || res.StatusCode > 399 {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			l.Error("error status code received")
			return &Response{err: fmt.Errorf("unexpected status code: %d", res.StatusCode)}
		}
		l = l.With(zap.String("response", string(bodyBytes)))
		l.Error("unexpected status code")

		status := meta.Status{
			Status:  meta.StatusFailure,
			Code:    int32(res.StatusCode),
			Message: "Unexpected error",
		}

		if len(bodyBytes) == 0 {
			return &Response{
				err: &meta.StatusError{
					ErrStatus: status,
				},
				body: bodyBytes,
			}
		}

		if err := json.Unmarshal(bodyBytes, &status); err != nil {
			return &Response{
				err: &meta.StatusError{
					ErrStatus: meta.Status{
						Status:  meta.StatusFailure,
						Code:    int32(res.StatusCode),
						Message: string(bodyBytes),
					},
				},
				body: bodyBytes,
			}
		}

		return &Response{
			err: &meta.StatusError{
				ErrStatus: status,
			},
			body: bodyBytes,
		}
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		l.Error("failed to read response body", zap.Error(err))
		return &Response{err: err}
	}

	return &Response{
		response: res,
		body:     bodyBytes,
	}
}

type Response struct {
	err      error
	response *http.Response
	body     []byte
}

func (r *Response) Into(into interface{}) error {
	if r.err != nil {
		return r.err
	}
	if into == nil {
		return nil
	}
	if len(r.body) == 0 {
		return fmt.Errorf("0-length body")
	}
	return json.Unmarshal(r.body, &into)
}

func (r *Response) Raw() ([]byte, *http.Response, error) {
	return r.body, r.response, r.err
}

func (r *Response) Error() error {
	return r.err
}
