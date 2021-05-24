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

type CustomResourceDefinitionsGetter interface {
	CustomResourceDefinitions() CustomResourceDefinitionInterface
}

type CustomResourceDefinitionInterface interface {
	Create(ctx context.Context, customResourceDefinition *v1.CustomResourceDefinition, opts metav1.CreateOptions) (*v1.CustomResourceDefinition, error)
	Update(ctx context.Context, customResourceDefinition *v1.CustomResourceDefinition, opts metav1.UpdateOptions) (*v1.CustomResourceDefinition, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.CustomResourceDefinition, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.CustomResourceDefinitionList, error)
	Watch(ctx context.Context, opts coremetav1.ListResourcesOptions) (watch.Interface, error)
}

// customResourceDefinitions implements CustomResourceDefinitionInterface
type customResourceDefinitions struct {
	client rest.Interface
}

// newCustomResourceDefinitions returns a customResourceDefinitions
func newCustomResourceDefinitions(c *CoreV1Client) *customResourceDefinitions {
	return &customResourceDefinitions{
		client: c.RESTClient(),
	}
}

// Get takes name of the customResourceDefinition, and returns the corresponding customResourceDefinition object, and an error if there is any.
func (c *customResourceDefinitions) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.CustomResourceDefinition, err error) {
	result = &v1.CustomResourceDefinition{}
	err = c.client.Get().
		Resource("customresourcedefinitions").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CustomResourceDefinitions that match those selectors.
func (c *customResourceDefinitions) List(ctx context.Context, opts metav1.ListOptions) (result *v1.CustomResourceDefinitionList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.CustomResourceDefinitionList{}
	err = c.client.Get().
		Resource("customresourcedefinitions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Create takes the representation of a customResourceDefinition and creates it.  Returns the server's representation of the customResourceDefinition, and an error, if there is any.
func (c *customResourceDefinitions) Create(ctx context.Context, customResourceDefinition *v1.CustomResourceDefinition, opts metav1.CreateOptions) (result *v1.CustomResourceDefinition, err error) {
	result = &v1.CustomResourceDefinition{}
	err = c.client.Post().
		Resource("customresourcedefinitions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(customResourceDefinition).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a customResourceDefinition and updates it. Returns the server's representation of the customResourceDefinition, and an error, if there is any.
func (c *customResourceDefinitions) Update(ctx context.Context, customResourceDefinition *v1.CustomResourceDefinition, opts metav1.UpdateOptions) (result *v1.CustomResourceDefinition, err error) {
	result = &v1.CustomResourceDefinition{}
	err = c.client.Put().
		Resource("customresourcedefinitions").
		Name(customResourceDefinition.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(customResourceDefinition).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the customResourceDefinition and deletes it. Returns an error if one occurs.
func (c *customResourceDefinitions) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("customresourcedefinitions").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// Watch returns a watch.Interface that watches the requested customResourceDefinitions.
func (c *customResourceDefinitions) Watch(ctx context.Context, opts coremetav1.ListResourcesOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("customresourcedefinitions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}
