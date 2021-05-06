package nrc

import (
	"context"
	"github.com/nrc-no/core/apps/api/pkgs2/client/rest"
	"github.com/nrc-no/core/apps/api/pkgs2/models"
)

type NrcCoreClient struct {
	restClient rest.Interface
}

func New(rest rest.Interface) *NrcCoreClient {
	return &NrcCoreClient{restClient: rest}
}

func (c *NrcCoreClient) FormDefinitions() FormDefinitionsInterface {
	return &formDefinitionsClient{client: c.restClient}
}

type FormDefinitionsInterface interface {
	Get(ctx context.Context, name string) (*models.FormDefinition, error)
	List(ctx context.Context) (*models.FormDefinitionList, error)
}

type formDefinitionsClient struct {
	client rest.Interface
}

func (c *formDefinitionsClient) Get(ctx context.Context, name string) (result *models.FormDefinition, err error) {
	result = &models.FormDefinition{}
	err = c.client.Get().
		Resource("formdefinitions").
		Name(name).
		Do(ctx).
		Into(result)
	return
}

func (c *formDefinitionsClient) List(ctx context.Context) (result *models.FormDefinitionList, err error) {
	result = &models.FormDefinitionList{}
	err = c.client.Get().
		Resource("formdefinitions").
		Do(ctx).
		Into(result)
	return
}
