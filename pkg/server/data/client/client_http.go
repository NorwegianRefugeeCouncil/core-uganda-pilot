package client

import (
	"context"
	"net/http"
	"strconv"

	"github.com/nrc-no/core/pkg/server/data/api"
)

type httpClient struct {
	baseURL string
	client  *client
}

type HTTPClient interface {
	CreateTable(ctx context.Context, request api.Table) (api.Table, error)
	GetRecord(ctx context.Context, request api.GetRecordRequest) (api.Record, error)
	GetRecords(ctx context.Context, request api.GetRecordsRequest) (api.RecordList, error)
	GetTables(ctx context.Context, request api.GetTablesRequest) (api.GetTablesResponse, error)
	GetTable(ctx context.Context, request api.GetTableRequest) (api.Table, error)
	GetChanges(ctx context.Context, request api.GetChangesRequest) (api.Changes, error)
	PutRecord(ctx context.Context, request api.PutRecordRequest) (api.Record, error)
}

func NewClient(baseURL string) HTTPClient {
	return &httpClient{
		client:  &client{client: http.DefaultClient},
		baseURL: baseURL,
	}
}

func (c *httpClient) GetRecord(ctx context.Context, request api.GetRecordRequest) (api.Record, error) {
	var response api.Record
	err := c.client.Get().
		URL(c.baseURL+"/apis/data.nrc.no/v1/tables/"+request.TableName+"/records/"+request.RecordID).
		WithHeader("Accept", "application/json").
		WithQueryParam("revision", request.Revision.String()).
		Do(ctx).Into(&response)
	return response, err
}

func (c *httpClient) GetTables(ctx context.Context, _ api.GetTablesRequest) (api.GetTablesResponse, error) {
	var response api.GetTablesResponse
	err := c.client.Get().
		URL(c.baseURL+"/apis/data.nrc.no/v1/tables").
		WithHeader("Accept", "application/json").
		Do(ctx).Into(&response)
	return response, err
}

func (c *httpClient) GetTable(ctx context.Context, request api.GetTableRequest) (api.Table, error) {
	var response api.Table
	err := c.client.Get().
		URL(c.baseURL+"/apis/data.nrc.no/v1/tables/"+request.TableName).
		WithHeader("Accept", "application/json").
		Do(ctx).Into(&response)
	return response, err
}

func (c *httpClient) GetRecords(ctx context.Context, request api.GetRecordsRequest) (api.RecordList, error) {
	var response api.RecordList
	err := c.client.Get().
		URL(c.baseURL+"/apis/data.nrc.no/v1/tables/"+request.TableName+"/records").
		WithHeader("Accept", "application/json").
		WithQueryParam("revisions", strconv.FormatBool(request.Revisions)).
		Do(ctx).Into(&response)
	return response, err
}

func (c *httpClient) GetChanges(ctx context.Context, request api.GetChangesRequest) (api.Changes, error) {
	var response api.Changes
	err := c.client.Get().
		URL(c.baseURL+"/apis/data.nrc.no/v1/changes").
		WithHeader("Accept", "application/json").
		WithQueryParam("since", strconv.FormatInt(request.Since, 10)).
		Do(ctx).Into(&response)
	return response, err
}

func (c *httpClient) PutRecord(ctx context.Context, request api.PutRecordRequest) (api.Record, error) {
	var response api.Record
	err := c.client.Put().
		URL(c.baseURL+"/apis/data.nrc.no/v1/tables/"+request.Record.Table+"/records/"+request.Record.ID).
		WithQueryParam("replication", strconv.FormatBool(request.IsReplication)).
		WithBody(request).
		Do(ctx).Into(&response)
	return response, err
}

func (c *httpClient) CreateTable(ctx context.Context, request api.Table) (api.Table, error) {
	var response api.Table
	err := c.client.Put().
		URL(c.baseURL+"/apis/data.nrc.no/v1/tables/"+request.Name).
		WithHeader("Accept", "application/json").
		WithHeader("Content-Type", "application/json").
		WithBody(request).
		Do(ctx).Into(&response)
	return response, err
}
