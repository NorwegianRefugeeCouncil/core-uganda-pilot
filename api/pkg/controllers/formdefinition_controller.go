package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nrc-no/core/api/pkg/apis/core/helpers"
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
	"strings"
	"time"
)

// FormDefinitionController is supposed to create CustomResourceDefinitions from
// FormDefinitions. It will create the CustomResourceDefinition if it doesn't exist,
// or will reconcile it if there is a discrepancy between the FormDefinition
// and the CustomResourceDefinition.
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
		err := c.syncFormDefinition(formDef)
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

func (c *FormDefinitionController) syncFormDefinition(formDef *v1.FormDefinition) error {

	startTime := time.Now()
	defer func() {
		klog.V(4).Infof("finished syncing form definition: %s (%v)", formDef.Name, time.Since(startTime))
	}()

	crdName := formDef.Name
	crd, err := c.crdLister.Get(crdName)
	if errors.IsNotFound(err) {
		crd := helpers.ConvertToCustomResourceDefinition(formDef)
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

// reconcileFormDefinition Detects and reconciles changes between FormDefinition and CustomResourceDefinition
func (c *FormDefinitionController) reconcileFormDefinition(formDef *v1.FormDefinition, actualCrd *apiextensionsv1.CustomResourceDefinition) error {

	desiredCrd := helpers.ConvertToCustomResourceDefinition(formDef)

	desiredBytes, err := json.Marshal(desiredCrd.Spec)
	if err != nil {
		return err
	}
	actualBytes, err := json.Marshal(actualCrd.Spec)
	if err != nil {
		return err
	}

	if strings.Compare(string(actualBytes), string(desiredBytes)) == 0 {
		return nil
	}

	klog.Infof("detected difference between actual and desired CustomResourceDefinition %v: \nactual:\n%s\ndesired\n%s", formDef.Name, string(actualBytes), string(desiredBytes))

	desiredCrd.ObjectMeta.ResourceVersion = actualCrd.ResourceVersion
	_, err = c.crdClient.Update(context.TODO(), desiredCrd, metav1.UpdateOptions{})
	if err != nil {
		klog.Warningf("error while trying to update CustomResourceDefinition: %v", err)
		return err
	}

	return nil
}
