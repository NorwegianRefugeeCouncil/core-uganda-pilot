package registration

import (
	"context"
	discoveryv1 "github.com/nrc-no/core/api/pkg/apis/discovery/v1"
	informers "github.com/nrc-no/core/api/pkg/client/informers/discovery/v1"
	listers "github.com/nrc-no/core/api/pkg/client/listers/discovery/v1"
	discoveryclient "github.com/nrc-no/core/api/pkg/client/typed/discovery/v1"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"reflect"
	"sync"
	"time"
)

const (
	manageOnStart            = "onstart"
	manageContinuously       = "true"
	AutoRegisterManagedLabel = "aggregator.nrc.no/automanaged"
)

type AutoAPIServiceRegistration interface {
	AddAPIServiceToSyncOnStart(in *discoveryv1.APIService)
	AddAPIServiceToSync(in *discoveryv1.APIService)
	RemoveAPIServiceToSync(name string)
}

type AutoRegisterController struct {
	apiServiceLister listers.APIServiceLister
	apiServiceSynced cache.InformerSynced
	apiServiceClient discoveryclient.APIServicesGetter

	apiServicesToSyncLock sync.RWMutex
	apiServicesToSync     map[string]*discoveryv1.APIService

	syncHandler func(apiServiceName string) error

	syncedSuccessfullyLock *sync.RWMutex
	syncedSuccessfully     map[string]bool

	apiServicesAtStart map[string]bool

	queue workqueue.RateLimitingInterface
}

func NewAutoRegisterController(
	apiServiceInformer informers.APIServiceInformer,
	apiServiceClient discoveryclient.APIServicesGetter,
) *AutoRegisterController {

	c := &AutoRegisterController{
		apiServiceLister:       apiServiceInformer.Lister(),
		apiServiceSynced:       apiServiceInformer.Informer().HasSynced,
		apiServiceClient:       apiServiceClient,
		apiServicesToSync:      map[string]*discoveryv1.APIService{},
		apiServicesAtStart:     map[string]bool{},
		syncedSuccessfullyLock: &sync.RWMutex{},
		syncedSuccessfully:     map[string]bool{},
		queue:                  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "AutoRegisterController"),
	}
	c.syncHandler = c.checkAPIService

	apiServiceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			castObj := obj.(*discoveryv1.APIService)
			c.queue.Add(castObj.Name)
		},
		UpdateFunc: func(_, obj interface{}) {
			castObj := obj.(*discoveryv1.APIService)
			c.queue.Add(castObj.Name)
		},
		DeleteFunc: func(obj interface{}) {
			castObj, ok := obj.(*discoveryv1.APIService)
			if !ok {
				tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					logrus.Errorf("couldn't get object from tombstone %#v", obj)
					return
				}
				castObj, ok = tombstone.Obj.(*discoveryv1.APIService)
				if !ok {
					logrus.Errorf("tombstone contained unexpected object: %#v", obj)
					return
				}
			}
			c.queue.Add(castObj.Name)
		},
	})

	return c
}

func (c *AutoRegisterController) Run(stopCh <-chan struct{}) {

	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	logrus.Infof("starting AutoRegisterController")
	defer logrus.Infof("shutting down AutoRegisterController")

	if !cache.WaitForCacheSync(stopCh, c.apiServiceSynced) {
		logrus.Errorf("unable to sync caches for AutoRegisterController")
		return
	}

	if services, err := c.apiServiceLister.List(labels.Everything()); err == nil {
		for _, service := range services {
			c.apiServicesAtStart[service.Name] = true
		}
	}

	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh

}

func (c *AutoRegisterController) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *AutoRegisterController) processNextWorkItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}

	defer c.queue.Done(key)

	err := c.syncHandler(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	logrus.Errorf("%s failed with %v", key, err)
	c.queue.AddRateLimited(key)
	return true

}

func (c *AutoRegisterController) checkAPIService(name string) (err error) {

	desired := c.GetAPIServiceToSync(name)
	curr, err := c.apiServiceLister.Get(name)

	hasSynced := c.hasSyncedSuccessfully(name)
	if !hasSynced {
		defer func() {
			if err == nil {
				c.setSyncedSuccessfully(name)
			}
		}()
	}

	switch {
	case err != nil && !apierrors.IsNotFound(err):
		return err
	case apierrors.IsNotFound(err) && desired == nil:
		return nil
	case isAutoManagedOnStart(desired) && hasSynced:
		return nil
	case apierrors.IsNotFound(err) && desired != nil:
		_, err := c.apiServiceClient.APIServices().Create(context.TODO(), desired, metav1.CreateOptions{})
		if apierrors.IsAlreadyExists(err) {
			return nil
		}
		return err
	case !isAutoManaged(curr):
		return nil
	case isAutoManagedOnStart(curr) && !c.apiServicesAtStart[name]:
		return nil
	case isAutoManagedOnStart(curr) && hasSynced:
		return nil
	case desired == nil:
		opts := metav1.DeleteOptions{Preconditions: metav1.NewUIDPreconditions(string(curr.UID))}
		err := c.apiServiceClient.APIServices().Delete(context.TODO(), curr.Name, opts)
		if apierrors.IsNotFound(err) || apierrors.IsConflict(err) {
			return nil
		}
		return err
	case reflect.DeepEqual(curr.Spec, desired.Spec):
		return nil
	}

	apiService := curr.DeepCopy()
	apiService.Spec = desired.Spec
	_, err = c.apiServiceClient.APIServices().Update(context.TODO(), apiService, metav1.UpdateOptions{})
	if apierrors.IsNotFound(err) || apierrors.IsConflict(err) {
		return nil
	}
	return err
}

func (c *AutoRegisterController) GetAPIServiceToSync(name string) *discoveryv1.APIService {
	c.apiServicesToSyncLock.RLock()
	defer c.apiServicesToSyncLock.RUnlock()
	return c.apiServicesToSync[name]
}

func (c *AutoRegisterController) AddAPIServiceToSyncOnStart(in *discoveryv1.APIService) {
	c.addApiServiceToSync(in, manageOnStart)
}

func (c *AutoRegisterController) AddAPIServiceToSync(in *discoveryv1.APIService) {
	c.addApiServiceToSync(in, manageContinuously)
}
func (c *AutoRegisterController) addApiServiceToSync(in *discoveryv1.APIService, syncType string) {
	c.apiServicesToSyncLock.Lock()
	defer c.apiServicesToSyncLock.Unlock()

	apiService := in.DeepCopy()
	if apiService.Labels == nil {
		apiService.Labels = map[string]string{}
	}
	apiService.Labels[AutoRegisterManagedLabel] = syncType
	c.apiServicesToSync[apiService.Name] = apiService
	c.queue.Add(apiService.Name)
}

func (c *AutoRegisterController) RemoveAPIServiceToSync(name string) {
	c.apiServicesToSyncLock.Lock()
	defer c.apiServicesToSyncLock.Unlock()

	delete(c.apiServicesToSync, name)
	c.queue.Add(name)
}

func (c *AutoRegisterController) hasSyncedSuccessfully(name string) bool {
	c.syncedSuccessfullyLock.RLock()
	defer c.syncedSuccessfullyLock.RUnlock()
	return c.syncedSuccessfully[name]
}

func (c *AutoRegisterController) setSyncedSuccessfully(name string) {
	c.syncedSuccessfullyLock.Lock()
	defer c.syncedSuccessfullyLock.Unlock()
	c.syncedSuccessfully[name] = true
}

func autoManagedType(service *discoveryv1.APIService) string {
	if service == nil {
		return ""
	}
	return service.Labels[AutoRegisterManagedLabel]
}

func isAutoManagedOnStart(service *discoveryv1.APIService) bool {
	return autoManagedType(service) == manageOnStart
}
func isAutoManaged(service *discoveryv1.APIService) bool {
	managedType := autoManagedType(service)
	return managedType == manageContinuously || managedType == manageOnStart
}
