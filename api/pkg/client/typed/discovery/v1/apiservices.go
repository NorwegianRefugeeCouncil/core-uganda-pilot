package v1

import (
	"context"
	discoveryv1 "github.com/nrc-no/core/api/pkg/apis/discovery/v1"
	coremetav1 "github.com/nrc-no/core/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/api/pkg/client/rest"
	"github.com/nrc-no/core/api/pkg/client/typed/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"time"
)

// APIServicesGetter has a method to return a APIServiceInterface.
// A group's client should implement this interface.
type APIServicesGetter interface {
	APIServices() APIServiceInterface
}

type APIServiceInterface interface {
	Create(ctx context.Context, APIService *discoveryv1.APIService, opts metav1.CreateOptions) (*discoveryv1.APIService, error)
	Update(ctx context.Context, APIService *discoveryv1.APIService, opts metav1.UpdateOptions) (*discoveryv1.APIService, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*discoveryv1.APIService, error)
	List(ctx context.Context, opts metav1.ListOptions) (*discoveryv1.APIServiceList, error)
	Watch(ctx context.Context, opts coremetav1.ListResourcesOptions) (watch.Interface, error)
}

// apiServices implements APIServiceInterface
type apiServices struct {
	client rest.Interface
}

// newAPIServices returns a apiServices
func newAPIServices(c *DiscoveryV1Client) *apiServices {
	return &apiServices{
		client: c.RESTClient(),
	}
}

// Get takes name of the apiService, and returns the corresponding apiService object, and an error if there is any.
func (c *apiServices) Get(ctx context.Context, name string, options metav1.GetOptions) (result *discoveryv1.APIService, err error) {
	result = &discoveryv1.APIService{}
	err = c.client.Get().
		Resource("apiservices").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of APIServices that match those selectors.
func (c *apiServices) List(ctx context.Context, opts metav1.ListOptions) (result *discoveryv1.APIServiceList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &discoveryv1.APIServiceList{}
	err = c.client.Get().
		Resource("apiservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Create takes the representation of a apiService and creates it.  Returns the server's representation of the apiService, and an error, if there is any.
func (c *apiServices) Create(ctx context.Context, apiService *discoveryv1.APIService, opts metav1.CreateOptions) (result *discoveryv1.APIService, err error) {
	result = &discoveryv1.APIService{}
	err = c.client.Post().
		Resource("apiservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(apiService).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a apiService and updates it. Returns the server's representation of the apiService, and an error, if there is any.
func (c *apiServices) Update(ctx context.Context, apiService *discoveryv1.APIService, opts metav1.UpdateOptions) (result *discoveryv1.APIService, err error) {
	result = &discoveryv1.APIService{}
	err = c.client.Put().
		Resource("apiservices").
		Name(apiService.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(apiService).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the apiService and deletes it. Returns an error if one occurs.
func (c *apiServices) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("apiservices").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// Watch returns a watch.Interface that watches the requested customResourceDefinitions.
func (c *apiServices) Watch(ctx context.Context, opts coremetav1.ListResourcesOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("apiservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}
