package client

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/rest"
	types2 "github.com/nrc-no/core/pkg/types"
)

type Client interface {
	ListDatabases(ctx context.Context, into *types2.DatabaseList) error
	CreateDatabase(ctx context.Context, obj *types2.Database, into *types2.Database) error
	DeleteDatabase(ctx context.Context, id string) error
	CreateFolder(ctx context.Context, obj *types2.Folder, into *types2.Folder) error
	DeleteFolder(ctx context.Context, id string) error
	CreateForm(ctx context.Context, obj *types2.FormDefinition, into *types2.FormDefinition) error
	ListForms(ctx context.Context, into *types2.FormDefinitionList) error
	CreateRecord(ctx context.Context, obj *types2.Record, into *types2.Record) error
	ListRecords(ctx context.Context, options types2.RecordListOptions, into *types2.RecordList) error
}

func NewClientFromConfig(c rest.Config) Client {
	return &client{
		c: rest.NewClient(&c),
	}
}

type client struct {
	c *rest.Client
}

func (c *client) ListDatabases(ctx context.Context, into *types2.DatabaseList) error {
	return c.c.Get().Path("/databases").Do(ctx).Into(into)
}

func (c *client) CreateDatabase(ctx context.Context, obj *types2.Database, into *types2.Database) error {
	return c.c.Post().Body(obj).Path("/databases").Do(ctx).Into(into)
}

func (c *client) DeleteDatabase(ctx context.Context, id string) error {
	return c.c.Delete().Path(fmt.Sprintf("/databases/%s", id)).Do(ctx).Error()
}

func (c *client) CreateFolder(ctx context.Context, obj *types2.Folder, into *types2.Folder) error {
	return c.c.Post().Body(obj).Path("/folders").Do(ctx).Into(into)
}

func (c *client) DeleteFolder(ctx context.Context, id string) error {
	return c.c.Delete().Path(fmt.Sprintf("/folders/%s", id)).Do(ctx).Error()
}

func (c *client) CreateForm(ctx context.Context, obj *types2.FormDefinition, into *types2.FormDefinition) error {
	return c.c.Post().Body(obj).Path("/forms").Do(ctx).Into(into)
}

func (c *client) ListForms(ctx context.Context, into *types2.FormDefinitionList) error {
	return c.c.Get().Path("/forms").Do(ctx).Into(into)
}

func (c *client) CreateRecord(ctx context.Context, obj *types2.Record, into *types2.Record) error {
	return c.c.Post().Body(obj.Values).Path(fmt.Sprintf("/databases/%s/forms/%s/records", obj.DatabaseID, obj.FormID)).Do(ctx).Into(into)
}

func (c *client) ListRecords(ctx context.Context, options types2.RecordListOptions, into *types2.RecordList) error {
	return c.c.Get().Path(fmt.Sprintf("/databases/%s/forms/%s/records", options.DatabaseID, options.FormID)).Do(ctx).Into(into)
}
