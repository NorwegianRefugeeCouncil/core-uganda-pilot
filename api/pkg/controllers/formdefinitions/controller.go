package formdefinitions

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/api/pkg/apis/core/helpers"
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/api/pkg/client/core"
	informers "github.com/nrc-no/core/api/pkg/client/informers/core/v1"
	listers "github.com/nrc-no/core/api/pkg/client/listers/core/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"reflect"
	"time"
)

// FormDefinitionController will
// 1. Create a CustomResourceDefinition for a given FormDefinition
// 2. Delete a CustomResourceDefinition if the FormDefinition was deleted
// 3. Update a CustomResourceDefinition if the FormDefinition was updated
// It will poll every 15 seconds to reconcile the above.
type FormDefinitionController struct {
	formDefinitionsLister listers.FormDefinitionLister
	formDefinitionsSynced cache.InformerSynced

	crdLister listers.CustomResourceDefinitionLister
	crdSynced cache.InformerSynced

	syncFn func(name string) error

	queue workqueue.RateLimitingInterface

	cli core.Interface
}

func NewFormDefinitionController(
	cli core.Interface,
	formDefinitionsInformer informers.FormDefinitionInformer,
	crdsInformer informers.CustomResourceDefinitionInformer,
) *FormDefinitionController {
	controller := &FormDefinitionController{
		cli:                   cli,
		formDefinitionsLister: formDefinitionsInformer.Lister(),
		formDefinitionsSynced: formDefinitionsInformer.Informer().HasSynced,
		crdLister:             crdsInformer.Lister(),
		crdSynced:             crdsInformer.Informer().HasSynced,
		queue:                 workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "FormDefinitionController"),
	}

	formDefinitionsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.addFormDefinition,
		UpdateFunc: controller.updateFormDefinition,
		DeleteFunc: controller.deleteFormDefinition,
	})
	crdsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.addCustomResourceDefinition,
		UpdateFunc: controller.updateCustomResourceDefinition,
		DeleteFunc: controller.deleteCustomResourceDefinition,
	})

	controller.syncFn = controller.sync

	return controller
}

func (c *FormDefinitionController) Run(stopCh <-chan struct{}, syncedCh chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	logrus.Infof("starting FormDefinitionController")
	defer logrus.Infof("shutting down FormDefinitionController")

	if !cache.WaitForCacheSync(stopCh, c.formDefinitionsSynced) {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

	if !cache.WaitForCacheSync(stopCh, c.crdSynced) {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

	if err := wait.PollImmediateUntil(time.Second*10, func() (done bool, err error) {
		formDefinitions, err := c.formDefinitionsLister.List(labels.Everything())
		if err != nil {
			utilruntime.HandleError(fmt.Errorf("failed to initially list formdefinitions: %v", err))
			return false, nil
		}
		for _, formDefinition := range formDefinitions {
			if err := c.syncFn(formDefinition.Name); err != nil {
				utilruntime.HandleError(fmt.Errorf("failed to initially sync formdefinition %s: %v", formDefinition.Name, err))
				return false, nil
			}
		}
		return true, nil

	}, stopCh); err == wait.ErrWaitTimeout {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for formdefinitions to sync: %v", err))
		return
	} else if err != nil {
		panic(fmt.Errorf("unexpected error: %v", err))
	}
	close(syncedCh)

	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh
}

func (c *FormDefinitionController) addFormDefinition(obj interface{}) {
	castObj := obj.(*corev1.FormDefinition)
	logrus.Infof("adding formdefinition %s", castObj.Name)
	c.queue.Add(castObj.Name)
}

func (c *FormDefinitionController) updateFormDefinition(oldObj interface{}, newObj interface{}) {
	castObj := oldObj.(*corev1.FormDefinition)
	logrus.Infof("updating formdefinition %s", castObj.Name)
	c.queue.Add(castObj.Name)
}

func (c *FormDefinitionController) deleteFormDefinition(obj interface{}) {
	castObj := obj.(*corev1.FormDefinition)
	logrus.Infof("deleting formdefinition %s", castObj.Name)
	c.queue.Add(castObj.Name)
}

func (c *FormDefinitionController) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *FormDefinitionController) processNextWorkItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.sync(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	utilruntime.HandleError(fmt.Errorf("failed to sync formdefinition: %v", err))
	c.queue.AddRateLimited(key)

	return true
}

func (c *FormDefinitionController) sync(s string) error {

	formDef, err := c.formDefinitionsLister.Get(s)
	if errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	desiredCrd := helpers.ConvertToCustomResourceDefinition(formDef)

	actualCrd, err := c.crdLister.Get(s)
	if errors.IsNotFound(err) {
		_, err := c.cli.CoreV1().CustomResourceDefinitions().Create(context.TODO(), desiredCrd, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	if !reflect.DeepEqual(actualCrd.Spec, desiredCrd.Spec) {
		_, err := c.cli.CoreV1().CustomResourceDefinitions().Update(context.TODO(), desiredCrd, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	return nil

}

func (c *FormDefinitionController) addCustomResourceDefinition(obj interface{}) {
	castObj := obj.(*corev1.CustomResourceDefinition)
	c.queue.Add(castObj.Name)
}

func (c *FormDefinitionController) updateCustomResourceDefinition(oldObj interface{}, newObj interface{}) {
	castObj := oldObj.(*corev1.CustomResourceDefinition)
	c.queue.Add(castObj.Name)
}

func (c *FormDefinitionController) deleteCustomResourceDefinition(obj interface{}) {
	castObj := obj.(*corev1.CustomResourceDefinition)
	c.queue.Add(castObj.Name)
}
