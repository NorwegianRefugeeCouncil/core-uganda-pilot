package cases

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/cases/api"
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

func (c *Client) List(ctx context.Context, listOptions ListOptions) (*api.CaseList, error) {
	req, err := http.NewRequest("GET", c.basePath+"/apis/v1/cases", nil)
	if err != nil {
		return nil, err
	}

	qry := req.URL.Query()
	if len(listOptions.PartyID) > 0 {
		qry.Set("partyId", listOptions.PartyID)
	}
	if len(listOptions.CaseTypeID) > 0 {
		qry.Set("caseTypeId", listOptions.CaseTypeID)
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
	var list api.CaseList
	if err := json.Unmarshal(bodyBytes, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (c *Client) Get(ctx context.Context, id string) (*api.Case, error) {
	req, err := http.NewRequest("GET", c.basePath+"/apis/v1/cases/"+id, nil)
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
	var v api.Case
	if err := json.Unmarshal(bodyBytes, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) Update(ctx context.Context, kase *api.Case) (*api.Case, error) {
	bodyBytes, err := json.Marshal(kase)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", c.basePath+"/apis/v1/cases/"+kase.ID, bytes.NewReader(bodyBytes))
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
	var v api.Case
	if err := json.Unmarshal(responseBytes, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) Create(ctx context.Context, kase *api.Case) (*api.Case, error) {
	bodyBytes, err := json.Marshal(kase)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.basePath+"/apis/v1/cases", bytes.NewReader(bodyBytes))
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
	var v api.Case
	if err := json.Unmarshal(responseBytes, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
