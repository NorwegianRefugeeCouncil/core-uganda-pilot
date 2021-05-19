package controllers

import (
	"context"
	"fmt"
	v1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	informersv1 "github.com/nrc-no/core/api/pkg/generated/informers/externalversions/core/v1"
	corev1 "github.com/nrc-no/core/api/pkg/generated/listers/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsclientv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	apiextensionsv1informers "k8s.io/apiextensions-apiserver/pkg/client/informers/externalversions/apiextensions/v1"
	listers "k8s.io/apiextensions-apiserver/pkg/client/listers/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"strconv"
	"time"
)

type FormDefinitionController struct {
	lister    corev1.FormDefinitionLister
	queue     workqueue.RateLimitingInterface
	synced    cache.InformerSynced
	crdLister listers.CustomResourceDefinitionLister
	crdClient apiextensionsclientv1.CustomResourceDefinitionInterface
}

func NewFormDefinitionController(
	formDefinitionInformer informersv1.FormDefinitionInformer,
	crdClient apiextensionsclientv1.CustomResourceDefinitionInterface,
	crdInformer apiextensionsv1informers.CustomResourceDefinitionInformer,
) *FormDefinitionController {
	c := &FormDefinitionController{
		queue:     workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "formdefinition"),
		lister:    formDefinitionInformer.Lister(),
		synced:    formDefinitionInformer.Informer().HasSynced,
		crdClient: crdClient,
		crdLister: crdInformer.Lister(),
	}

	formDefinitionInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			crd := obj.(*v1.FormDefinition)
			c.addFormDefinition(crd)
		},
	})

	return c
}

func (c *FormDefinitionController) addFormDefinition(crd *v1.FormDefinition) {
	c.queue.Add(crd)
}

func (c *FormDefinitionController) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()
	defer klog.Info("Shutting down FormDefinitionController")

	klog.Info("Starting FormDefinitionController")

	if !cache.WaitForNamedCacheSync("formdefinition", stopCh, c.synced) {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

	go wait.Until(c.worker, time.Second, stopCh)

	<-stopCh
}

func (c *FormDefinitionController) worker() {
	workFunc := func() bool {
		key, quit := c.queue.Get()
		if quit {
			return true
		}
		defer c.queue.Done(key)

		formDef := key.(*v1.FormDefinition)
		err := c.syncFormDefinitionFromKey(formDef)
		if err != nil {
			c.queue.Forget(key)
			return false
		}

		return false
	}
	for {
		if workFunc() {
			return
		}
	}
}

func (c *FormDefinitionController) syncFormDefinitionFromKey(formDef *v1.FormDefinition) error {

	startTime := time.Now()
	defer func() {
		klog.V(4).Infof("finished syncing form definition: %s (%v)", formDef.Name, time.Since(startTime))
	}()

	crdName := formDef.Name
	crd, err := c.crdLister.Get(crdName)
	if errors.IsNotFound(err) {
		crd := createCrdFromFormDefinition(formDef)
		_, err := c.crdClient.Create(context.TODO(), crd, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("error while retrieving form definition %s: %v", crdName, err))
		return err
	}
	return c.reconcileFormDefinition(formDef, crd)

}

func createCrdFromFormDefinition(formDefinition *v1.FormDefinition) *apiextensionsv1.CustomResourceDefinition {

	crd := &apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: formDefinition.Spec.Names.Plural + "." + formDefinition.Spec.Group,
		},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Group: formDefinition.Spec.Group,
			Scope: apiextensionsv1.ClusterScoped,
			Names: apiextensionsv1.CustomResourceDefinitionNames{
				Plural:   formDefinition.Spec.Names.Plural,
				Singular: formDefinition.Spec.Names.Singular,
				Kind:     formDefinition.Spec.Names.Kind,
			},
		},
	}

	for i, version := range formDefinition.Spec.Versions {
		crdVersion := apiextensionsv1.CustomResourceDefinitionVersion{
			Name:    version.Name,
			Storage: i == 0,
			Served:  true,
			Schema:  &apiextensionsv1.CustomResourceValidation{},
		}
		validation := createCrdValidationFromFormDefinitionVersion(formDefinition, version)
		crdVersion.Schema = validation
		crd.Spec.Versions = append(crd.Spec.Versions, crdVersion)
	}

	return crd
}

func createCrdValidationFromFormDefinitionVersion(fs *v1.FormDefinition, version v1.FormDefinitionVersion) *apiextensionsv1.CustomResourceValidation {

	specSchema := &apiextensionsv1.JSONSchemaProps{
		Description: `Defines the desired state fo ` + fs.Spec.Names.Kind,
		Type:        "object",
	}
	formSchema := version.Schema.FormSchema.Root
	walkFormSchema(formSchema, specSchema)

	return &apiextensionsv1.CustomResourceValidation{
		OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
			Description: "Schema for the " + fs.Spec.Names.Kind + " api",
			Type:        "object",
			Properties: map[string]apiextensionsv1.JSONSchemaProps{
				"apiVersion": {
					Description: `APIVersion defines the versioned schema of this representation
of an object. Servers should convert recognized schemas to the latest internal value, and may
reject unrecognized values.`,
					Type: "string",
				},
				"kind": {
					Description: `Kind is a string value representing the REST resource this 
object represents. Servers may infer this from the endpoint the client submits requests to.
Cannot be updated. In CamelCase.`,
					Type: "string",
				},
				"metadata": {
					Type: "object",
				},
				"spec": *specSchema,
			},
		},
	}
}

func walkFormSchema(element v1.FormElementDefinition, jsonProps *apiextensionsv1.JSONSchemaProps) {

	var intMultipleOf float64 = 1
	switch element.Type {
	case v1.IntegerType:
		jsonProps.Type = "number"
		jsonProps.MultipleOf = &intMultipleOf
	case v1.ShortTextType:
		jsonProps.Type = "string"
	case v1.LongTextType:
		jsonProps.Type = "string"
	case v1.SectionType:
		jsonProps.Type = "object"
	case v1.DateTimeType:
		jsonProps.Type = "string"
		jsonProps.Format = "datetime"
	case v1.DateType:
		jsonProps.Type = "string"
		jsonProps.Format = "date"
	case v1.SelectType:
	case v1.TimeType:
	}

	if element.MinLength != 0 {
		jsonProps.MinLength = &element.MinLength
	}
	if element.MaxLength != nil {
		jsonProps.MaxLength = element.MaxLength
	}
	if element.Max != "" {
		max, err := strconv.ParseFloat(element.Max, 64)
		if err == nil {
			jsonProps.Maximum = &max
		}
	}
	if element.Min != "" {
		min, err := strconv.ParseFloat(element.Min, 64)
		if err == nil {
			jsonProps.Minimum = &min
		}
	}
	if element.Pattern != "" {
		jsonProps.Pattern = element.Pattern
	}

	if jsonProps.Description == "" {
		jsonProps.Description = findDescription(element.Description)
	}

	for _, child := range element.Children {
		childJsonProps := &apiextensionsv1.JSONSchemaProps{}
		walkFormSchema(child, childJsonProps)

		if childJsonProps.Type == "object" {
			if childJsonProps.Properties != nil {
				if jsonProps.Properties == nil {
					jsonProps.Properties = map[string]apiextensionsv1.JSONSchemaProps{}
				}
				for key, props := range childJsonProps.Properties {
					jsonProps.Properties[key] = props
				}
				for _, propName := range childJsonProps.Required {
					jsonProps.Required = append(jsonProps.Required, propName)
				}
			}
		} else {
			if jsonProps.Properties == nil {
				jsonProps.Properties = map[string]apiextensionsv1.JSONSchemaProps{}
			}
			jsonProps.Properties[child.Key] = *childJsonProps
			if child.Required {
				jsonProps.Required = append(jsonProps.Required, child.Key)
			}
		}
	}
}

func findDescription(strs v1.TranslatedStrings) string {
	for _, str := range strs {
		if str.Locale == "en" {
			return str.Value
		}
	}
	if len(strs) > 0 {
		return strs[0].Value
	}
	return ""
}

func (c *FormDefinitionController) reconcileFormDefinition(formDefinition *v1.FormDefinition, crd *apiextensionsv1.CustomResourceDefinition) error {
	return nil
}
