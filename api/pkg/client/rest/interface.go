package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	v1 "github.com/nrc-no/core/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/api/pkg/exceptions"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Interface interface {
	Verb(verb string) *Request
	Get() *Request
	Put() *Request
	Post() *Request
	Delete() *Request
}

type RESTClient struct {
	base   *url.URL
	client *http.Client
}

var _ Interface = &RESTClient{}

func NewRestClient(baseURL *url.URL, client *http.Client) *RESTClient {
	base := *baseURL
	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}
	if client == nil {
		client = http.DefaultClient
	}
	return &RESTClient{
		base:   &base,
		client: client,
	}
}

func NewRequest(c *RESTClient) *Request {
	return nil
}

func (r *RESTClient) Verb(verb string) *Request {
	return NewRequest(r).Verb(verb)
}

func (r *RESTClient) Get() *Request {
	return r.Verb("GET")
}

func (r *RESTClient) Put() *Request {
	return r.Verb("PUT")
}

func (r *RESTClient) Post() *Request {
	return r.Verb("POST")
}

func (r *RESTClient) Delete() *Request {
	return r.Verb("DELETE")
}

type Request struct {
	c        *RESTClient
	verb     string
	name     string
	resource string
	err      error
	params   url.Values
	headers  http.Header
	body     io.Reader
}

func (r *Request) Verb(verb string) *Request {
	r.verb = verb
	return r
}

func (r *Request) Resource(resource string) *Request {
	if r.err != nil {
		return r
	}
	if len(r.resource) != 0 {
		r.err = fmt.Errorf("resource already set to %s. cannot change to %s", r.resource, resource)
		return r
	}
	r.resource = resource
	return r
}

func (r *Request) Name(name string) *Request {
	if r.err != nil {
		return r
	}
	if len(r.name) != 0 {
		r.err = fmt.Errorf("name already set to %s. cannot change to %s", r.name, name)
		return r
	}
	r.name = name
	return r
}

func (r *Request) Param(paramName, s string) *Request {
	if r.err != nil {
		return r
	}
	return r.setParam(paramName, s)
}

func (r *Request) setParam(paramName, value string) *Request {
	if r.params == nil {
		r.params = make(url.Values)
	}
	r.params[paramName] = append(r.params[paramName], value)
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
		data, err := json.Marshal(obj)
		if err != nil {
			r.err = err
			return r
		}
		r.SetHeader("Content-Type", "application/json")
		r.body = bytes.NewReader(data)
		return r
	}
	return r
}

func (r *Request) URL() *url.URL {
	p := ""
	if len(r.resource) != 0 {
		p = path.Join(p, strings.ToLower(r.resource))
	}
	if len(r.name) != 0 {
		p = path.Join(p, r.name)
	}
	finalURL := &url.URL{}
	if r.c.base != nil {
		*finalURL = *r.c.base
	}

	query := url.Values{}
	for key, values := range r.params {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	finalURL.RawQuery = query.Encode()
	return finalURL
}

type Result struct {
	body        []byte
	contentType string
	err         error
	statusCode  int
}

func (r *Request) Do(ctx context.Context) Result {
	var result Result
	client := r.c.client
	if client == nil {
		client = http.DefaultClient
	}

	url := r.URL().String()
	req, err := http.NewRequest(r.verb, url, r.body)
	if err != nil {
		result.err = err
		return result
	}
	req = req.WithContext(ctx)
	req.Header = r.headers

	resp, err := client.Do(req)
	if err != nil {
		result.err = err
		return result
	}

	result.statusCode = resp.StatusCode
	result.contentType = resp.Header.Get("Content-Type")
	result.body = result.body
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		result.err = err
	}
	result.body = bodyBytes

	if result.statusCode < http.StatusOK || result.statusCode > 299 {
		if result.contentType == "application/json" && len(result.body) != 0 {
			var status = v1.Status{}
			if err := json.Unmarshal(result.body, &status); err != nil {
				result.err = errors.New(string(result.body))
			}
			result.err = exceptions.NewStatusError(status)
			return result
		}
		result.err = errors.New(string(result.body))
	}

	return result
}

func (r Result) Into(obj interface{}) error {
	if r.err != nil {
		return r.err
	}

	if len(r.body) == 0 {
		return fmt.Errorf("0-length response with status code %d and content-type %s", r.statusCode, r.contentType)
	}

	if err := DecodeErrorStatus(r.body); err != nil {
		r.err = err
		return err
	}

	return json.Unmarshal(r.body, obj)
}

func (r Result) Error() error {
	if r.err == nil || len(r.body) == 0 {
		return r.err
	}
	if err := DecodeErrorStatus(r.body); err != nil {
		return err
	}
	return r.err
}

func DecodeErrorStatus(body []byte) error {
	findKind := struct {
		APIVersion string `json:"apiVersion,omitempty"`
		Kind       string `json:"kind,omitempty"`
	}{}
	if err := json.Unmarshal(body, &findKind); err != nil {
		return err
	}
	if findKind.APIVersion == "meta/v1" && findKind.Kind == "Status" {
		var status = v1.Status{}
		if err := json.Unmarshal(body, &status); err != nil {
			return err
		}
		if status.Status != v1.StatusSuccess {
			return exceptions.NewStatusError(status)
		}
	}
	return nil
}
