package apiservice

import (
	"fmt"
	v1 "github.com/nrc-no/core/api/pkg/apis/discovery/v1"
	informers "github.com/nrc-no/core/api/pkg/client/informers/discovery/v1"
	listers "github.com/nrc-no/core/api/pkg/client/listers/discovery/v1"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type APIHandlerManager interface {
	AddAPIService(apiService *v1.APIService) error
	RemoveAPIService(apiServiceName string)
}

type APIServiceRegistrationController struct {
	apiHandlerManager APIHandlerManager

	apiServiceLister listers.APIServiceLister
	apiServiceSynced cache.InformerSynced

	syncFn func(key string) error

	queue workqueue.RateLimitingInterface
}

func NewAPIServiceRegistrationController(
	apiServiceInformer informers.APIServiceInformer,
	apiHandlerManager APIHandlerManager,
) *APIServiceRegistrationController {

	c := &APIServiceRegistrationController{
		apiHandlerManager: apiHandlerManager,
		apiServiceLister:  apiServiceInformer.Lister(),
		apiServiceSynced:  apiServiceInformer.Informer().HasSynced,
		queue:             workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "APIServiceRegistrationController"),
	}

	apiServiceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addAPIService,
		UpdateFunc: c.updateAPIService,
		DeleteFunc: c.deleteAPIService,
	})

	c.syncFn = c.sync

	return c
}

func (c *APIServiceRegistrationController) sync(key string) error {
	apiService, err := c.apiServiceLister.Get(key)
	if apierrors.IsNotFound(err) {
		c.apiHandlerManager.RemoveAPIService(key)
		return nil
	}
	if err != nil {
		return err
	}
	return c.apiHandlerManager.AddAPIService(apiService)
}

func (c *APIServiceRegistrationController) Run(stopCh <-chan struct{}, handlerSyncedCh chan<- struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	logrus.Infof("starting APIServiceRegistrationController")
	defer logrus.Infof("shutting down APIServiceRegistrationController")

	if !cache.WaitForCacheSync(stopCh) {
		logrus.Error("unable to sync caches for APIServiceRegistrationController", c.apiServiceSynced)
		return
	}

	if err := wait.PollImmediateUntil(time.Second, func() (bool, error) {
		services, err := c.apiServiceLister.List(labels.Everything())
		if err != nil {
			logrus.Errorf("failed to initially list APIServices: %v", err)
			return false, nil
		}
		for _, s := range services {
			if err := c.apiHandlerManager.AddAPIService(s); err != nil {
				logrus.Errorf("failed to initially sync APIService %s: %v", s.Name, err)
				return false, nil
			}
		}
		return true, nil
	}, stopCh); err == wait.ErrWaitTimeout {
		logrus.Errorf("timed out waiting for handler to initialize: %v", err)
		return
	} else {
		panic(fmt.Errorf("unexpected error: %v", err))
	}
	close(handlerSyncedCh)

	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh

}

func (c *APIServiceRegistrationController) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *APIServiceRegistrationController) processNextWorkItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}

	err := c.syncFn(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	logrus.Errorf("%v failed: %v", key, err)
	c.queue.AddRateLimited(key)

	return true
}

func (c *APIServiceRegistrationController) enqueueInternal(obj *v1.APIService) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		logrus.Errorf("couldn't get object key for object %#v: %v", obj, err)
		return
	}
	c.queue.Add(key)
}

func (c *APIServiceRegistrationController) addAPIService(obj interface{}) {
	castObj := obj.(*v1.APIService)
	logrus.Tracef("adding APIService %s", castObj.Name)
	c.enqueueInternal(castObj)
}

func (c *APIServiceRegistrationController) updateAPIService(obj, _ interface{}) {
	castObj := obj.(*v1.APIService)
	logrus.Tracef("updating APIService %s", castObj.Name)
	c.enqueueInternal(castObj)
}

func (c *APIServiceRegistrationController) deleteAPIService(obj interface{}) {
	castObj, ok := obj.(*v1.APIService)
	if !ok {
		tombStone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			logrus.Errorf("couldn't get object from tombstone: %#v", obj)
			return
		}
		castObj, ok = tombStone.Obj.(*v1.APIService)
		if !ok {
			logrus.Errorf("tombstone contained object thtat is not expected: %#v", tombStone.Obj)
			return
		}
	}
	logrus.Tracef("deleting APIService %s", castObj.Name)
	c.enqueueInternal(castObj)
}

func (c *APIServiceRegistrationController) Enqueue() {
	apiServices, err := c.apiServiceLister.List(labels.Everything())
	if err != nil {
		logrus.Errorf("failed to enqueue APIServices: %v", err)
		return
	}
	for _, service := range apiServices {
		c.addAPIService(service)
	}
}
