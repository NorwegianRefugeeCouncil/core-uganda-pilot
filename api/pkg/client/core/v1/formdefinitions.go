package v1

import (
	"context"
	v1 "github.com/nrc-no/coreapi/pkg/apis/core/v1"
	metav1 "github.com/nrc-no/coreapi/pkg/apis/meta/v1"
	"github.com/nrc-no/coreapi/pkg/client/rest"
)

type FormDefinitionsGetter interface {
	FormDefinitions() FormDefinitionInterface
}

type FormDefinitionInterface interface {
	Create(ctx context.Context, formDefinition *v1.FormDefinition, opts metav1.CreateOptions) (*v1.FormDefinition, error)
	Update(ctx context.Context, formDefinition *v1.FormDefinition, opts metav1.UpdateOptions) (*v1.FormDefinition, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.FormDefinition, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.FormDefinitionList, error)
}

type formDefinitions struct {
	client rest.Interface
}

var _ FormDefinitionInterface = &formDefinitions{}

func newFormDefinitions(c *CoreV1Client) *formDefinitions {
	return &formDefinitions{
		client: c.restClient,
	}
}

func (c *formDefinitions) Create(ctx context.Context, formDefinition *v1.FormDefinition, opts metav1.CreateOptions) (*v1.FormDefinition, error) {
	result := &v1.FormDefinition{}
	err := c.client.Post().
		Resource("formdefinitions").
		Body(formDefinition).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *formDefinitions) Update(ctx context.Context, formDefinition *v1.FormDefinition, opts metav1.UpdateOptions) (*v1.FormDefinition, error) {
	result := &v1.FormDefinition{}
	err := c.client.Put().
		Resource("formdefinitions").
		Name(formDefinition.Name).
		Body(formDefinition).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *formDefinitions) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("formdefinitions").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *formDefinitions) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.FormDefinition, error) {
	result := &v1.FormDefinition{}
	err := c.client.Get().
		Resource("formdefinitions").
		Name(name).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *formDefinitions) List(ctx context.Context, opts metav1.ListOptions) (*v1.FormDefinitionList, error) {
	result := &v1.FormDefinitionList{}
	err := c.client.Get().
		Resource("formdefinitions").
		Do(ctx).
		Into(result)
	return result, err
}
