package nrc

import (
	"context"
	"github.com/nrc-no/core/apps/api/pkg/apis"
	"github.com/nrc-no/core/apps/api/pkg/client/rest"
)

type NrcCoreClient struct {
	restClient rest.Interface
}

func New(rest rest.Interface) *NrcCoreClient {
	return &NrcCoreClient{restClient: rest}
}

func NewForConfig(c *rest.Config) (*NrcCoreClient, error) {
	config := *c
	config.Group = "core.nrc.no"
	config.Version = "v1"
	config.APIPath = "apis"
	config.ContentType = "application/json"
	config.AcceptContentType = "application/json"
	restClient, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return New(restClient), nil
}

func (c *NrcCoreClient) FormDefinitions() FormDefinitionsInterface {
	return &formDefinitionsClient{client: c.restClient}
}

type FormDefinitionsInterface interface {
	Create(ctx context.Context, formDefinition *apis.FormDefinition) (*apis.FormDefinition, error)
	Get(ctx context.Context, id string) (*apis.FormDefinition, error)
	List(ctx context.Context) (*apis.FormDefinitionList, error)
	Update(ctx context.Context, formDefinition *apis.FormDefinition) (result *apis.FormDefinition, err error)
}

type formDefinitionsClient struct {
	client rest.Interface
}

func (c *formDefinitionsClient) Get(ctx context.Context, name string) (result *apis.FormDefinition, err error) {
	result = &apis.FormDefinition{}
	err = c.client.Get().
		Resource("formdefinitions").
		Name(name).
		Do(ctx).
		Into(result)
	return
}

func (c *formDefinitionsClient) List(ctx context.Context) (result *apis.FormDefinitionList, err error) {
	result = &apis.FormDefinitionList{}
	err = c.client.Get().
		Resource("formdefinitions").
		Do(ctx).
		Into(result)
	return
}

func (c *formDefinitionsClient) Create(ctx context.Context, formDefinition *apis.FormDefinition) (result *apis.FormDefinition, err error) {
	result = &apis.FormDefinition{}
	err = c.client.Post().
		Resource("formdefinitions").
		Body(formDefinition).
		Do(ctx).
		Into(result)
	return
}

func (c *formDefinitionsClient) Update(ctx context.Context, formDefinition *apis.FormDefinition) (result *apis.FormDefinition, err error) {
	result = &apis.FormDefinition{}
	err = c.client.Put().
		Resource("formdefinitions").
		Name(formDefinition.UID).
		Body(formDefinition).
		Do(ctx).
		Into(result)
	return
}
