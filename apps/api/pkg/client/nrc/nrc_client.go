package nrc

import (
	"context"
	"github.com/nrc-no/core/apps/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/apps/api/pkg/client/rest"
	"github.com/nrc-no/core/apps/api/pkg/watch"
)

type NrcCoreClient struct {
	restClient rest.Interface
}

func New(rest rest.Interface) *NrcCoreClient {
	return &NrcCoreClient{restClient: rest}
}

func NewForConfig(c *rest.Config) (*NrcCoreClient, error) {
	config := *c
	config.GroupVersion = v1.SchemeGroupVersion
	config.APIPath = "apis"
	config.ContentType = "application/json"
	config.AcceptContentType = "application/json"
	config.ContentConfig.Serializer = v1.Codecs

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
	Create(ctx context.Context, formDefinition *v1.FormDefinition) (*v1.FormDefinition, error)
	Get(ctx context.Context, id string) (*v1.FormDefinition, error)
	List(ctx context.Context) (*v1.FormDefinitionList, error)
	Update(ctx context.Context, formDefinition *v1.FormDefinition) (result *v1.FormDefinition, err error)
	Watch(ctx context.Context) (p watch.Interface, err error)
}

type formDefinitionsClient struct {
	client rest.Interface
}

func (c *formDefinitionsClient) Get(ctx context.Context, name string) (result *v1.FormDefinition, err error) {
	result = &v1.FormDefinition{}
	err = c.client.Get().
		Resource("formdefinitions").
		Name(name).
		Do(ctx).
		Into(result)
	return
}

func (c *formDefinitionsClient) List(ctx context.Context) (result *v1.FormDefinitionList, err error) {
	result = &v1.FormDefinitionList{}
	err = c.client.Get().
		Resource("formdefinitions").
		Do(ctx).
		Into(result)
	return
}

func (c *formDefinitionsClient) Create(ctx context.Context, formDefinition *v1.FormDefinition) (result *v1.FormDefinition, err error) {
	result = &v1.FormDefinition{}
	err = c.client.Post().
		Resource("formdefinitions").
		Body(formDefinition).
		Do(ctx).
		Into(result)
	return
}

func (c *formDefinitionsClient) Update(ctx context.Context, formDefinition *v1.FormDefinition) (result *v1.FormDefinition, err error) {
	result = &v1.FormDefinition{}
	err = c.client.Put().
		Resource("formdefinitions").
		Name(formDefinition.UID).
		Body(formDefinition).
		Do(ctx).
		Into(result)
	return
}

func (c *formDefinitionsClient) Watch(ctx context.Context) (w watch.Interface, err error) {
	return c.client.Get().
		Resource("formdefinitions").
		Watch(ctx)
}
