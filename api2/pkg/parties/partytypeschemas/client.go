package partytypeschemas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
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

func (c *Client) List(ctx context.Context) (*partytypes.PartyTypeSchemaList, error) {
	req, err := http.NewRequest("GET", c.basePath+"/apis/v1/partytypeschemas", nil)
	if err != nil {
		return nil, err
	}

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
	var list partytypes.PartyTypeSchemaList
	if err := json.Unmarshal(bodyBytes, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (c *Client) Get(ctx context.Context, id string) (*partytypes.PartyTypeSchema, error) {
	req, err := http.NewRequest("GET", c.basePath+"/apis/v1/partytypeschemas/"+id, nil)
	if err != nil {
		return nil, err
	}
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
	var v partytypes.PartyTypeSchema
	if err := json.Unmarshal(bodyBytes, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) Update(ctx context.Context, partyTypeSchema *partytypes.PartyTypeSchema) (*partytypes.PartyTypeSchema, error) {
	bodyBytes, err := json.Marshal(partyTypeSchema)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", c.basePath+"/apis/v1/partytypeschemas/"+partyTypeSchema.ID, bytes.NewReader(bodyBytes))
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
	var v partytypes.PartyTypeSchema
	if err := json.Unmarshal(responseBytes, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) Create(ctx context.Context, partyTypeSchema *partytypes.PartyTypeSchema) (*partytypes.PartyTypeSchema, error) {
	bodyBytes, err := json.Marshal(partyTypeSchema)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.basePath+"/apis/v1/partytypeschemas", bytes.NewReader(bodyBytes))
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
	var v partytypes.PartyTypeSchema
	if err := json.Unmarshal(responseBytes, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
