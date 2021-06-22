package parties

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

type ListOptions struct {
	PartyTypeID string `json:"partyTypeId" bson:"partyTypeId"`
	SearchParam string `json:"searchParam" bson:"searchParam"`
}

func (c *Client) List(ctx context.Context, listOptions ListOptions) (*PartyList, error) {
	req, err := http.NewRequest("GET", c.basePath+"/apis/v1/parties", nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	auth.SetAuthorizationHeader(ctx, req)
	qry := req.URL.Query()
	if len(listOptions.PartyTypeID) > 0 {
		qry.Set("partyTypeId", listOptions.PartyTypeID)
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
	var list PartyList
	if err := json.Unmarshal(bodyBytes, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (c *Client) Get(ctx context.Context, id string) (*Party, error) {
	req, err := http.NewRequest("GET", c.basePath+"/apis/v1/parties/"+id, nil)
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
	var v Party
	if err := json.Unmarshal(bodyBytes, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) Update(ctx context.Context, party *Party) (*Party, error) {
	bodyBytes, err := json.Marshal(party)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", c.basePath+"/apis/v1/parties/"+party.ID, bytes.NewReader(bodyBytes))
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
	var v Party
	if err := json.Unmarshal(responseBytes, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) Create(ctx context.Context, party *Party) (*Party, error) {
	bodyBytes, err := json.Marshal(party)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.basePath+"/apis/v1/parties", bytes.NewReader(bodyBytes))
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
	var v Party
	if err := json.Unmarshal(responseBytes, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
