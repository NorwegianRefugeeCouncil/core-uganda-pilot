package v1

import (
	"context"
	v1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	coremetav1 "github.com/nrc-no/core/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/api/pkg/client/core/scheme"
	"github.com/nrc-no/core/api/pkg/client/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"time"
)

// FormDefinitionsGetter has a method to return a FormDefinitionInterface.
// A group's client should implement this interface.
type FormDefinitionsGetter interface {
	FormDefinitions() FormDefinitionInterface
}

type FormDefinitionInterface interface {
	Create(ctx context.Context, FormDefinition *v1.FormDefinition, opts metav1.CreateOptions) (*v1.FormDefinition, error)
	Update(ctx context.Context, FormDefinition *v1.FormDefinition, opts metav1.UpdateOptions) (*v1.FormDefinition, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.FormDefinition, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.FormDefinitionList, error)
	Watch(ctx context.Context, opts coremetav1.ListResourcesOptions) (watch.Interface, error)
}

// formDefinitions implements FormDefinitionInterface
type formDefinitions struct {
	client rest.Interface
}

// newFormDefinitions returns a formDefinitions
func newFormDefinitions(c *CoreV1Client) *formDefinitions {
	return &formDefinitions{
		client: c.RESTClient(),
	}
}

// Get takes name of the formDefinition, and returns the corresponding formDefinition object, and an error if there is any.
func (c *formDefinitions) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.FormDefinition, err error) {
	result = &v1.FormDefinition{}
	err = c.client.Get().
		Resource("formdefinitions").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of FormDefinitions that match those selectors.
func (c *formDefinitions) List(ctx context.Context, opts metav1.ListOptions) (result *v1.FormDefinitionList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.FormDefinitionList{}
	err = c.client.Get().
		Resource("formdefinitions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Create takes the representation of a formDefinition and creates it.  Returns the server's representation of the formDefinition, and an error, if there is any.
func (c *formDefinitions) Create(ctx context.Context, formDefinition *v1.FormDefinition, opts metav1.CreateOptions) (result *v1.FormDefinition, err error) {
	result = &v1.FormDefinition{}
	err = c.client.Post().
		Resource("formdefinitions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(formDefinition).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a formDefinition and updates it. Returns the server's representation of the formDefinition, and an error, if there is any.
func (c *formDefinitions) Update(ctx context.Context, formDefinition *v1.FormDefinition, opts metav1.UpdateOptions) (result *v1.FormDefinition, err error) {
	result = &v1.FormDefinition{}
	err = c.client.Put().
		Resource("formdefinitions").
		Name(formDefinition.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(formDefinition).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the formDefinition and deletes it. Returns an error if one occurs.
func (c *formDefinitions) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("formdefinitions").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// Watch returns a watch.Interface that watches the requested customResourceDefinitions.
func (c *formDefinitions) Watch(ctx context.Context, opts coremetav1.ListResourcesOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("formdefinitions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}
