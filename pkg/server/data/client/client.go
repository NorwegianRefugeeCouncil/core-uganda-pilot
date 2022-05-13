package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type client struct {
	client *http.Client
}

func (c *client) Method(method string) *request {
	return &request{
		client: c,
		method: method,
	}
}

func (c *client) Post() *request {
	return c.Method("POST")
}

func (c *client) Get() *request {
	return c.Method("GET")
}

func (c *client) Put() *request {
	return c.Method("PUT")
}

func (c *client) Delete() *request {
	return c.Method("DELETE")
}

type request struct {
	err         error
	method      string
	client      *client
	queryParams map[string]string
	headers     map[string]string
	body        []byte
	url         string
}

func (r *request) URL(url string) *request {
	r.url = url
	return r
}

func (r *request) WithMethod(method string) *request {
	r.method = method
	return r
}

func (r *request) WithQueryParam(key, value string) *request {
	if r.err != nil {
		return r
	}
	if r.queryParams == nil {
		r.queryParams = make(map[string]string)
	}
	r.queryParams[key] = value
	return r
}

func (r *request) WithHeader(key, value string) *request {
	if r.err != nil {
		return r
	}
	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	r.headers[key] = value
	return r
}

func (r *request) WithBody(body interface{}) *request {
	if r.err != nil {
		return r
	}
	switch b := body.(type) {
	case []byte:
		r.body = b
	case string:
		r.body = []byte(b)
	default:
		jsonBytes, err := json.Marshal(b)
		if err != nil {
			r.err = err
			return r
		}
		r.body = jsonBytes
	}
	return r
}

func (r *request) Do(ctx context.Context) *response {
	if r.err != nil {
		return &response{err: r.err}
	}
	var body []byte
	if r.body != nil {
		body = r.body
	}
	req, err := http.NewRequestWithContext(ctx, r.method, r.url, bytes.NewBuffer(body))
	if err != nil {
		return &response{err: r.err}
	}
	if req.Header == nil {
		req.Header = make(http.Header)
	}
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}
	if len(r.body) > 0 && (r.headers == nil || r.headers["Content-Type"] == "") {
		req.Header.Set("Content-Type", "application/json")
	}
	if r.method == "POST" || r.method == "PUT" {
		if r.headers["Accept"] == "" {
			req.Header.Set("Accept", "application/json")
		}
	}
	for k, v := range r.queryParams {
		q := req.URL.Query()
		q.Add(k, v)
		req.URL.RawQuery = q.Encode()
	}
	resp, err := r.client.client.Do(req)
	if err != nil {
		return &response{err: r.err}
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			r.err = err
		}
	}()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return &response{err: r.err}
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err := fmt.Errorf("Error occured with status code %d: %s\n", resp.StatusCode, string(body))
		return &response{
			err:        err,
			statusCode: resp.StatusCode,
			body:       body,
			url:        resp.Request.URL.String(),
			headers:    resp.Header,
		}
	}
	return &response{
		statusCode: resp.StatusCode,
		body:       body,
		headers:    resp.Header,
		url:        resp.Request.URL.String(),
	}

}

type response struct {
	err        error
	statusCode int
	headers    http.Header
	body       []byte
	url        string
}

func (r *response) Into(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	if r.body == nil {
		return nil
	}
	return json.Unmarshal(r.body, v)
}
