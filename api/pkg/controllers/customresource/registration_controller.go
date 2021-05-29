package customresource

import (
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	discoveryv1 "github.com/nrc-no/core/api/pkg/apis/discovery/v1"
	informers "github.com/nrc-no/core/api/pkg/client/informers/core/v1"
	listers "github.com/nrc-no/core/api/pkg/client/listers/core/v1"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type AutoAPIServiceRegistration interface {
	AddAPIServiceToSync(in *discoveryv1.APIService)
	RemoveAPIServiceToSync(name string)
}

type CRDRegistrationController struct {
	crdLister listers.CustomResourceDefinitionLister
	crdSynced cache.InformerSynced

	apiServiceRegistration AutoAPIServiceRegistration

	syncHandler      func(groupVersion schema.GroupVersion) error
	syncedInitialSet chan struct{}

	queue workqueue.RateLimitingInterface
}

func NewCRDRegistrationController(
	crdInformer informers.CustomResourceDefinitionInformer,
	apiServiceRegistration AutoAPIServiceRegistration,
) *CRDRegistrationController {
	c := &CRDRegistrationController{
		crdLister:              crdInformer.Lister(),
		crdSynced:              crdInformer.Informer().HasSynced,
		apiServiceRegistration: apiServiceRegistration,
		syncedInitialSet:       make(chan struct{}),
		queue:                  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "CRDRegistrationController"),
	}
	c.syncHandler = c.handleVersionUpdate
	crdInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			c.enqueueCRD(obj.(*corev1.CustomResourceDefinition))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			c.enqueueCRD(oldObj.(*corev1.CustomResourceDefinition))
			c.enqueueCRD(newObj.(*corev1.CustomResourceDefinition))
		},
		DeleteFunc: func(obj interface{}) {
			castObj, ok := obj.(*corev1.CustomResourceDefinition)
			if !ok {
				tombStone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					logrus.Errorf("couldn't get object from tombstone: %#v", obj)
					return
				}
				castObj, ok = tombStone.Obj.(*corev1.CustomResourceDefinition)
				if !ok {
					logrus.Errorf("tombstone container unexpected object: %#v", obj)
					return
				}
			}
			c.enqueueCRD(castObj)
		},
	})

	return c
}

func (c *CRDRegistrationController) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	logrus.Infof("starting CRDRegistrationController")
	defer logrus.Infof("shutting down CRDRegistrationController")

	if !cache.WaitForNamedCacheSync("CRDRegistrationController", stopCh, c.crdSynced) {
		return
	}

	if crds, err := c.crdLister.List(labels.Everything()); err != nil {
		logrus.Error(err)
	} else {
		for _, crd := range crds {
			for _, version := range crd.Spec.Versions {
				if err := c.syncHandler(schema.GroupVersion{Group: crd.Spec.Group, Version: version.Name}); err != nil {
					logrus.Errorf("failed to sync crd version %s.%s: %v", crd.Spec.Group, version.Name, err)
				}
			}
		}
	}
	close(c.syncedInitialSet)

	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh

}

func (c *CRDRegistrationController) WaitForInitialSync() {
	<-c.syncedInitialSet
}

func (c *CRDRegistrationController) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *CRDRegistrationController) processNextWorkItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}

	defer c.queue.Done(key)
	err := c.syncHandler(key.(schema.GroupVersion))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	logrus.Errorf("%s failed with: %v", key, err)
	c.queue.AddRateLimited(key)
	return true
}

func (c CRDRegistrationController) enqueueCRD(obj *corev1.CustomResourceDefinition) {
	for _, version := range obj.Spec.Versions {
		c.queue.Add(schema.GroupVersion{Group: obj.Spec.Group, Version: version.Name})
	}
}

func (c CRDRegistrationController) handleVersionUpdate(groupVersion schema.GroupVersion) error {
	apiServiceName := groupVersion.Version + "." + groupVersion.Group

	crds, err := c.crdLister.List(labels.Everything())
	if err != nil {
		return err
	}

	for _, crd := range crds {
		if crd.Spec.Group != groupVersion.Group {
			continue
		}
		for _, version := range crd.Spec.Versions {
			if version.Name != groupVersion.Version || !version.Served {
				continue
			}
		}
		c.apiServiceRegistration.AddAPIServiceToSync(&discoveryv1.APIService{
			ObjectMeta: metav1.ObjectMeta{
				Name: apiServiceName,
			},
			Spec: discoveryv1.APIServiceSpec{
				Group:   groupVersion.Group,
				Version: groupVersion.Version,
			},
		})
		return nil
	}

	c.apiServiceRegistration.RemoveAPIServiceToSync(apiServiceName)
	return nil
}
