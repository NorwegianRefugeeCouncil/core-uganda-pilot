package client

import (
	"context"
	"net/http"
	"path"
	"strconv"

	"github.com/nrc-no/core/pkg/server/data/api"
)

type httpClient struct {
	baseURL string
	client  *client
}

const tablesUrl = "/apis/data.nrc.no/v1/tables"
const mimeJson = "application/json"
const headerAccept = "Accept"
const headerContentType = "Content-Type"

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
		URL(path.Join(c.baseURL, tablesUrl, request.TableName, "records", request.RecordID)).
		WithHeader(headerAccept, mimeJson).
		WithQueryParam("revision", request.Revision.String()).
		Do(ctx).Into(&response)
	return response, err
}

func (c *httpClient) GetTables(ctx context.Context, _ api.GetTablesRequest) (api.GetTablesResponse, error) {
	var response api.GetTablesResponse
	err := c.client.Get().
		URL(path.Join(c.baseURL, "/apis/data.nrc.no/v1/tables")).
		WithHeader(headerAccept, mimeJson).
		Do(ctx).Into(&response)
	return response, err
}

func (c *httpClient) GetTable(ctx context.Context, request api.GetTableRequest) (api.Table, error) {
	var response api.Table
	err := c.client.Get().
		URL(path.Join(c.baseURL+tablesUrl, request.TableName)).
		WithHeader(headerAccept, mimeJson).
		Do(ctx).Into(&response)
	return response, err
}

func (c *httpClient) GetRecords(ctx context.Context, request api.GetRecordsRequest) (api.RecordList, error) {
	var response api.RecordList
	err := c.client.Get().
		URL(path.Join(c.baseURL, tablesUrl, request.TableName, "/records")).
		WithHeader(headerAccept, mimeJson).
		WithQueryParam("revisions", strconv.FormatBool(request.Revisions)).
		Do(ctx).Into(&response)
	return response, err
}

func (c *httpClient) GetChanges(ctx context.Context, request api.GetChangesRequest) (api.Changes, error) {
	var response api.Changes
	err := c.client.Get().
		URL(path.Join(c.baseURL, "/apis/data.nrc.no/v1/changes")).
		WithHeader(headerAccept, mimeJson).
		WithQueryParam("since", strconv.FormatInt(request.Since, 10)).
		Do(ctx).Into(&response)
	return response, err
}

func (c *httpClient) PutRecord(ctx context.Context, request api.PutRecordRequest) (api.Record, error) {
	var response api.Record
	err := c.client.Put().
		URL(path.Join(c.baseURL, tablesUrl, request.Record.Table, "records", request.Record.ID)).
		WithQueryParam("replication", strconv.FormatBool(request.IsReplication)).
		WithBody(request.Record).
		Do(ctx).Into(&response)
	return response, err
}

func (c *httpClient) CreateTable(ctx context.Context, request api.Table) (api.Table, error) {
	var response api.Table
	err := c.client.Put().
		URL(path.Join(c.baseURL, tablesUrl, request.Name)).
		WithHeader(headerAccept, mimeJson).
		WithHeader(headerContentType, mimeJson).
		WithBody(request).
		Do(ctx).Into(&response)
	return response, err
}
