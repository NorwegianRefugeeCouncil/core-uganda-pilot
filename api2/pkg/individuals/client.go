package individuals

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

type ListOptions struct {
	PartyTypes []string `json:"partyTypes" bson:"partyTypes"`
}

func (c *Client) List(ctx context.Context, listOptions ListOptions) (*IndividualList, error) {
	req, err := http.NewRequest("GET", c.basePath+"/apis/v1/individuals", nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	qry := req.URL.Query()
	if len(listOptions.PartyTypes) > 0 {
		for _, partyType := range listOptions.PartyTypes {
			qry.Add("partyTypes", partyType)
		}
	}
	req.URL.RawQuery = qry.Encode()
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
	var list IndividualList
	if err := json.Unmarshal(bodyBytes, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (c *Client) Get(ctx context.Context, id string) (*Individual, error) {
	req, err := http.NewRequest("GET", c.basePath+"/apis/v1/individuals/"+id, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
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
	var v Individual
	if err := json.Unmarshal(bodyBytes, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) Update(ctx context.Context, individual *Individual) (*Individual, error) {
	bodyBytes, err := json.Marshal(individual)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", c.basePath+"/apis/v1/individuals/"+individual.ID, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
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
	var v Individual
	if err := json.Unmarshal(responseBytes, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) Create(ctx context.Context, individual *Individual) (*Individual, error) {
	bodyBytes, err := json.Marshal(individual)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.basePath+"/apis/v1/individuals", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
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
	var v Individual
	if err := json.Unmarshal(responseBytes, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
