package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nrc-no/core/internal/validation"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type RESTConfig struct {
	Scheme     string
	Host       string
	HTTPClient *http.Client
	Headers    http.Header
}

type Client struct {
	config *RESTConfig
}

func NewClient(config *RESTConfig) *Client {
	return &Client{
		config: config,
	}
}

func (r *Client) Verb(verb string) *Request {
	return r.NewRequest().Verb(verb)
}

func (r *Client) Get() *Request {
	return r.Verb("GET")
}

func (r *Client) Post() *Request {
	return r.Verb("POST")
}

func (r *Client) Put() *Request {
	return r.Verb("PUT")
}

func (r *Client) Delete() *Request {
	return r.Verb("DELETE")
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
	default:
		vp := reflect.New(reflect.TypeOf(p))
		vp.Elem().Set(reflect.ValueOf(p))
		intf := vp.Interface()

		valuer, ok := intf.(UrlValuer)
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

	l := logrus.WithField("scheme", r.c.config.Scheme).WithField("host", r.c.config.Host).WithField("path", r.path)

	if r.err != nil {
		l.WithError(r.err).Errorf("could not send request")
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
		l.WithError(r.err).Errorf("could not create http request")
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

	if len(req.Header.Get("Accept")) == 0 {
		req.Header.Set("Accept", "application/json")
	}

	if len(req.Header.Get("Content-Type")) == 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	if xAuthUserSubject, ok := ctx.Value("Subject").(string); ok {
		req.Header.Set("X-Authenticated-User-Subject", xAuthUserSubject)
	}

	httpClient := r.c.config.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	res, err := httpClient.Do(req)
	if err != nil {
		l.WithError(err).Errorf("could not send http request")
		return &Response{err: err}
	}

	l = l.WithField("statusCode", res.StatusCode)

	if res.StatusCode < 200 || res.StatusCode > 399 {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			l.Errorf("error status code received")
			return &Response{err: fmt.Errorf("unexpected status code")}
		} else {
			l.WithField("response", string(bodyBytes)).Errorf("unexpected status code")
		}

		if len(bodyBytes) == 0 {
			return &Response{
				err: &validation.Status{
					Status:  validation.Failure,
					Code:    http.StatusInternalServerError,
					Message: "Unexpected error",
					Errors:  nil,
				},
				body: bodyBytes,
			}
		}

		var status validation.Status
		if err := json.Unmarshal(bodyBytes, &status); err != nil {
			return &Response{
				err: &validation.Status{
					Status:  validation.Failure,
					Code:    http.StatusInternalServerError,
					Message: "Unexpected error",
					Errors:  nil,
				},
				body: bodyBytes,
			}
		}

		return &Response{
			err:  &status,
			body: bodyBytes,
		}
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		l.WithError(r.err).Errorf("failed to read response body")
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

func (r *Response) Error() error {
	return r.err
}
