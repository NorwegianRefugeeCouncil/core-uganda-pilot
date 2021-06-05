package v1

import (
	"context"
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	coremetav1 "github.com/nrc-no/core/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/api/pkg/client/rest"
	scheme2 "github.com/nrc-no/core/api/pkg/client/typed/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"time"
)

// OperatingScopesGetter has a method to return a OperatingScopeInterface.
// A group's client should implement this interface.
type OperatingScopesGetter interface {
	OperatingScopes() OperatingScopeInterface
}

type OperatingScopeInterface interface {
	Create(ctx context.Context, OperatingScope *corev1.OrganizationScope, opts metav1.CreateOptions) (*corev1.OrganizationScope, error)
	Update(ctx context.Context, OperatingScope *corev1.OrganizationScope, opts metav1.UpdateOptions) (*corev1.OrganizationScope, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*corev1.OrganizationScope, error)
	List(ctx context.Context, opts metav1.ListOptions) (*corev1.OperatingScopeList, error)
	Watch(ctx context.Context, opts coremetav1.ListResourcesOptions) (watch.Interface, error)
}

// operatingScopes implements OperatingScopeInterface
type operatingScopes struct {
	client rest.Interface
}

// newOperatingScopes returns a operatingScopes
func newOperatingScopes(c *CoreV1Client) *operatingScopes {
	return &operatingScopes{
		client: c.RESTClient(),
	}
}

// Get takes name of the operatingScope, and returns the corresponding operatingScope object, and an error if there is any.
func (c *operatingScopes) Get(ctx context.Context, name string, options metav1.GetOptions) (result *corev1.OrganizationScope, err error) {
	result = &corev1.OrganizationScope{}
	err = c.client.Get().
		Resource("operatingscopes").
		Name(name).
		VersionedParams(&options, scheme2.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of OperatingScopes that match those selectors.
func (c *operatingScopes) List(ctx context.Context, opts metav1.ListOptions) (result *corev1.OperatingScopeList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &corev1.OperatingScopeList{}
	err = c.client.Get().
		Resource("operatingscopes").
		VersionedParams(&opts, scheme2.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Create takes the representation of a operatingScope and creates it.  Returns the server's representation of the operatingScope, and an error, if there is any.
func (c *operatingScopes) Create(ctx context.Context, operatingScope *corev1.OrganizationScope, opts metav1.CreateOptions) (result *corev1.OrganizationScope, err error) {
	result = &corev1.OrganizationScope{}
	err = c.client.Post().
		Resource("operatingscopes").
		VersionedParams(&opts, scheme2.ParameterCodec).
		Body(operatingScope).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a operatingScope and updates it. Returns the server's representation of the operatingScope, and an error, if there is any.
func (c *operatingScopes) Update(ctx context.Context, operatingScope *corev1.OrganizationScope, opts metav1.UpdateOptions) (result *corev1.OrganizationScope, err error) {
	result = &corev1.OrganizationScope{}
	err = c.client.Put().
		Resource("operatingscopes").
		Name(operatingScope.Name).
		VersionedParams(&opts, scheme2.ParameterCodec).
		Body(operatingScope).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the operatingScope and deletes it. Returns an error if one occurs.
func (c *operatingScopes) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("operatingscopes").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// Watch returns a watch.Interface that watches the requested customResourceDefinitions.
func (c *operatingScopes) Watch(ctx context.Context, opts coremetav1.ListResourcesOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("operatingscopes").
		VersionedParams(&opts, scheme2.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}
