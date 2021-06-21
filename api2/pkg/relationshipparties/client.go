package relationshipparties

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"io/ioutil"
	"net/http"
)

type Client struct {
	BasePath string
}

func NewClient(basePath string) *Client {
	return &Client{
		BasePath: basePath,
	}
}

func (c *Client) PickParty(ctx context.Context, pickPartyOptions PickPartyOptions) (*parties.PartyList, error) {
	req, err := http.NewRequest("GET", c.BasePath+"/apis/v1/relationshipparties/picker", nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	q := req.URL.Query()

	if len(pickPartyOptions.PartyTypeID) != 0 {
		q.Set("partyTypeId", pickPartyOptions.PartyTypeID)
	}
	if len(pickPartyOptions.SearchParam) != 0 {
		q.Set("searchParam", pickPartyOptions.SearchParam)
	}
	req.URL.RawQuery = q.Encode()

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
	var returnList parties.PartyList
	if err := json.Unmarshal(bodyBytes, &returnList); err != nil {
		return nil, err
	}
	return &returnList, nil
}
