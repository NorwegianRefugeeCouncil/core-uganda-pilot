package client

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/rest"
)

type Client interface {
	GetDatabase(ctx context.Context, id string, into *types.Database) error
	ListDatabases(ctx context.Context, into *types.DatabaseList) error
	CreateDatabase(ctx context.Context, obj *types.Database, into *types.Database) error
	DeleteDatabase(ctx context.Context, id string) error
	CreateFolder(ctx context.Context, obj *types.Folder, into *types.Folder) error
	ListFolders(ctx context.Context, into *types.FolderList) error
	DeleteFolder(ctx context.Context, id string) error
	CreateForm(ctx context.Context, obj *types.FormDefinition, into *types.FormDefinition) error
	ListForms(ctx context.Context, into *types.FormDefinitionList) error
	DeleteForm(ctx context.Context, id string) error
	CreateRecord(ctx context.Context, obj *types.Record, into *types.Record) error
	ListRecords(ctx context.Context, options types.RecordListOptions, into *types.RecordList) error
}

func NewClientFromConfig(c rest.Config) Client {
	return &client{
		c: rest.NewClient(&c),
	}
}

type client struct {
	c *rest.Client
}

func (c *client) GetDatabase(ctx context.Context, id string, into *types.Database) error {
	return c.c.Get().Path(fmt.Sprintf("/databases/%s", id)).Do(ctx).Into(into)
}

func (c *client) ListDatabases(ctx context.Context, into *types.DatabaseList) error {
	return c.c.Get().Path("/databases").Do(ctx).Into(into)
}

func (c *client) CreateDatabase(ctx context.Context, obj *types.Database, into *types.Database) error {
	return c.c.Post().Body(obj).Path("/databases").Do(ctx).Into(into)
}

func (c *client) DeleteDatabase(ctx context.Context, id string) error {
	return c.c.Delete().Path(fmt.Sprintf("/databases/%s", id)).Do(ctx).Error()
}

func (c *client) CreateFolder(ctx context.Context, obj *types.Folder, into *types.Folder) error {
	return c.c.Post().Body(obj).Path("/folders").Do(ctx).Into(into)
}

func (c *client) ListFolders(ctx context.Context, into *types.FolderList) error {
	return c.c.Get().Path(fmt.Sprintf("/folders")).Do(ctx).Into(into)
}

func (c *client) DeleteFolder(ctx context.Context, id string) error {
	return c.c.Delete().Path(fmt.Sprintf("/folders/%s", id)).Do(ctx).Error()
}

func (c *client) CreateForm(ctx context.Context, obj *types.FormDefinition, into *types.FormDefinition) error {
	return c.c.Post().Body(obj).Path("/forms").Do(ctx).Into(into)
}

func (c *client) ListForms(ctx context.Context, into *types.FormDefinitionList) error {
	return c.c.Get().Path("/forms").Do(ctx).Into(into)
}

func (c *client) DeleteForm(ctx context.Context, id string) error {
	return c.c.Delete().Path(fmt.Sprintf("/forms/%s", id)).Do(ctx).Error()
}

func (c *client) CreateRecord(ctx context.Context, obj *types.Record, into *types.Record) error {
	return c.c.Post().Body(obj.Values).Path(fmt.Sprintf("/databases/%s/forms/%s/records", obj.DatabaseID, obj.FormID)).Do(ctx).Into(into)
}

func (c *client) ListRecords(ctx context.Context, options types.RecordListOptions, into *types.RecordList) error {
	return c.c.Get().Path(fmt.Sprintf("/databases/%s/forms/%s/records", options.DatabaseID, options.FormID)).Do(ctx).Into(into)
}
