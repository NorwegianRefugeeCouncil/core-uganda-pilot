package form

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"net/url"
)

type Interface interface {
	Get(ctx context.Context, id string, options GetOptions) (*Definition, error)
	List(ctx context.Context, options ListOptions) (*Definition, error)
	Create(ctx context.Context, obj *Definition, options CreateOptions) (*Definition, error)
	Put(ctx context.Context, obj *Definition, options PutOptions) (*Definition, error)
	Validate(ctx context.Context, obj *Definition, options ValidateOptions) (*Definition, error)
}

type RESTFormDefinitionClient struct {
	c *rest.Client
}

func NewClientFromConfig(c *rest.Config) Interface {
	return &RESTFormDefinitionClient{
		c: rest.NewClient(c),
	}
}

type GetOptions struct {
	Version string
}

func (c *RESTFormDefinitionClient) Get(ctx context.Context, id string, options GetOptions) (*Definition, error) {
	var def Definition
	req := c.c.Get().Path(fmt.Sprintf("%s/%s", server.FormDefinitionsEndpoint, id))

	params := url.Values{}
	if len(options.Version) > 0 {
		params.Set("version", options.Version)
	}
	if len(params) > 0 {
		req = req.WithParams(params)
	}

	err := req.Do(ctx).Into(&def)
	return &def, err
}

type ListOptions struct {
}

func (c *RESTFormDefinitionClient) List(ctx context.Context, options ListOptions) (*Definition, error) {
	var def Definition
	err := c.c.Get().Path(fmt.Sprintf("%s", server.FormDefinitionsEndpoint)).Do(ctx).Into(&def)
	return &def, err
}

type CreateOptions struct {
}

func (c *RESTFormDefinitionClient) Create(ctx context.Context, obj *Definition, options CreateOptions) (*Definition, error) {
	var def Definition
	err := c.c.Post().Body(obj).Path(fmt.Sprintf("%s", server.FormDefinitionsEndpoint)).Do(ctx).Into(&def)
	return &def, err
}

type PutOptions struct {
}

func (c *RESTFormDefinitionClient) Put(ctx context.Context, obj *Definition, options PutOptions) (*Definition, error) {
	var def Definition
	err := c.c.Put().Body(obj).Path(fmt.Sprintf("%s/%s", server.FormDefinitionsEndpoint, obj.ID)).Do(ctx).Into(&def)
	return &def, err
}

type ValidateOptions struct {
}

func (c *RESTFormDefinitionClient) Validate(ctx context.Context, obj *Definition, options ValidateOptions) (*Definition, error) {
	var def Definition
	err := c.c.Post().Body(obj).Path(fmt.Sprintf("%s/%s/validate", server.FormDefinitionsEndpoint, obj.ID)).Do(ctx).Into(&def)
	return &def, err
}
