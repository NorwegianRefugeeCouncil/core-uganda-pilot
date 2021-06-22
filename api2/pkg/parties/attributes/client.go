package attributes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/auth"
	"io/ioutil"
	"net/http"
)

type Client struct {
	basePath string
}

func NewClient(basePath string) *Client {
	return &Client{
		basePath: basePath,
	}
}

func (c *Client) List(ctx context.Context) (*AttributeList, error) {
	req, err := http.NewRequest("GET", c.basePath+"/apis/v1/attributes", nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	auth.SetAuthorizationHeader(ctx, req)
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		body, bodyErr := ioutil.ReadAll(res.Body)
		if bodyErr == nil {
			return nil, fmt.Errorf("unexpected status code %d: %s", res.StatusCode, string(body))
		}
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var list AttributeList
	if err := json.Unmarshal(bodyBytes, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (c *Client) Get(ctx context.Context, id string) (*Attribute, error) {
	req, err := http.NewRequest("GET", c.basePath+"/apis/v1/attributes/"+id, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	auth.SetAuthorizationHeader(ctx, req)
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var v Attribute
	if err := json.Unmarshal(bodyBytes, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) Update(ctx context.Context, attribute *Attribute) (*Attribute, error) {
	bodyBytes, err := json.Marshal(attribute)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", c.basePath+"/apis/v1/attributes/"+attribute.ID, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	auth.SetAuthorizationHeader(ctx, req)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	responseBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var v Attribute
	if err := json.Unmarshal(responseBytes, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) Create(ctx context.Context, attribute *Attribute) (*Attribute, error) {
	bodyBytes, err := json.Marshal(attribute)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.basePath+"/apis/v1/attributes", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	auth.SetAuthorizationHeader(ctx, req)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	responseBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var v Attribute
	if err := json.Unmarshal(responseBytes, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
