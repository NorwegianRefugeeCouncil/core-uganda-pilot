package customresource

import (
	"fmt"
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	discoveryv1 "github.com/nrc-no/core/api/pkg/apis/discovery/v1"
	informers "github.com/nrc-no/core/api/pkg/client/informers/core/v1"
	listers "github.com/nrc-no/core/api/pkg/client/listers/core/v1"
	"github.com/nrc-no/core/api/pkg/endpoints/discovery"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"sort"
	"time"
)

type CRDDiscoveryController struct {
	versionHandler *CRDVersionDiscoveryHandler
	groupHandler   *CRDGroupDiscoveryHandler

	crdLister  listers.CustomResourceDefinitionLister
	crdsSynced cache.InformerSynced
	syncFn     func(version schema.GroupVersion) error

	codecFactory serializer.CodecFactory
	queue        workqueue.RateLimitingInterface
}

func NewCRDDiscoveryController(
	versionHandler *CRDVersionDiscoveryHandler,
	groupHandler *CRDGroupDiscoveryHandler,
	codecFactory serializer.CodecFactory,
	crdInformer informers.CustomResourceDefinitionInformer,
) *CRDDiscoveryController {
	controller := &CRDDiscoveryController{
		versionHandler: versionHandler,
		groupHandler:   groupHandler,

		codecFactory: codecFactory,

		crdLister:  crdInformer.Lister(),
		crdsSynced: crdInformer.Informer().HasSynced,
		queue:      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "CRDDiscovery"),
	}

	crdInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.addCustomResourceDefinition,
		UpdateFunc: controller.updateCustomResourceDefinition,
		DeleteFunc: controller.deleteCustomResourceDefinition,
	})

	controller.syncFn = controller.sync

	return controller

}

func (c *CRDDiscoveryController) sync(version schema.GroupVersion) error {

	var apiVersionsForDiscovery []discoveryv1.GroupVersionForDiscovery
	var apiResourcesForDiscovery []discoveryv1.APIResource
	versionsForDiscoveryMap := map[metav1.GroupVersion]bool{}

	crds, err := c.crdLister.List(labels.Everything())
	if err != nil {
		return err
	}

	foundVersion := false
	foundGroup := false

	for _, crd := range crds {

		if crd.Spec.Group != version.Group {
			continue
		}

		foundThisVersion := false
		for _, v := range crd.Spec.Versions {

			if !v.Served {
				continue
			}
			foundGroup = true

			gv := metav1.GroupVersion{Group: crd.Spec.Group, Version: v.Name}
			if !versionsForDiscoveryMap[gv] {
				versionsForDiscoveryMap[gv] = true
				apiVersionsForDiscovery = append(apiVersionsForDiscovery, discoveryv1.GroupVersionForDiscovery{
					GroupVersion: crd.Spec.Group + "/" + v.Name,
					Version:      v.Name,
				})
			}
			if v.Name == version.Version {
				foundThisVersion = true
			}
			if v.Storage {
				//
			}

		}

		if !foundThisVersion {
			continue
		}
		foundVersion = true

		verbs := discoveryv1.Verbs([]string{"delete", "get", "list", "create", "update", "delete"})

		apiResourcesForDiscovery = append(apiResourcesForDiscovery, discoveryv1.APIResource{
			Name:         crd.Spec.Names.Plural,
			SingularName: crd.Spec.Names.Singular,
			Namespaced:   false,
			Kind:         crd.Spec.Names.Kind,
			Verbs:        verbs,
		})

	}

	if !foundGroup {
		c.groupHandler.unsetDiscovery(version.Group)
		c.versionHandler.unsetDiscovery(version)
		return nil
	}

	sortGroupDiscoveryByKubeAwareVersion(apiVersionsForDiscovery)

	if len(apiVersionsForDiscovery) == 0 {
		apiVersionsForDiscovery = []discoveryv1.GroupVersionForDiscovery{}
	}

	apiGroup := discoveryv1.APIGroup{
		Name:             version.Group,
		Versions:         apiVersionsForDiscovery,
		PreferredVersion: apiVersionsForDiscovery[0],
	}
	c.groupHandler.setDiscovery(version.Group, discovery.NewAPIGroupHandler(c.codecFactory, apiGroup))

	if !foundVersion {
		c.versionHandler.unsetDiscovery(version)
		return nil
	}

	c.versionHandler.setDiscovery(version, discovery.NewAPIVersionHandler(c.codecFactory, version, discovery.APIResourceListerFunc(func() []discoveryv1.APIResource {
		return apiResourcesForDiscovery
	})))

	return nil

}

func sortGroupDiscoveryByKubeAwareVersion(gd []discoveryv1.GroupVersionForDiscovery) {
	sort.Slice(gd, func(i, j int) bool {
		return version.CompareKubeAwareVersionStrings(gd[i].Version, gd[j].Version) > 0
	})
}

func (c *CRDDiscoveryController) Run(stopCh <-chan struct{}, syncedCh chan<- struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShuttingDown()
	defer logrus.Infof("shutting down CRDDiscoveryController")

	logrus.Infof("starting CRDDiscoveryController")

	if !cache.WaitForCacheSync(stopCh, c.crdsSynced) {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

	if err := wait.PollImmediateUntil(time.Second*20, func() (done bool, err error) {
		crds, err := c.crdLister.List(labels.Everything())
		if err != nil {
			utilruntime.HandleError(fmt.Errorf("failed to initially list CRDs: %v", err))
			return false, nil
		}
		for _, crd := range crds {
			for _, v := range crd.Spec.Versions {
				gv := schema.GroupVersion{Group: crd.Spec.Group, Version: v.Name}
				if err := c.sync(gv); err != nil {
					utilruntime.HandleError(fmt.Errorf("failed to initially sync crd version %#v", gv))
					return false, nil
				}
			}
		}
		return true, nil
	}, stopCh); err == wait.ErrWaitTimeout {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for discovery endpoint to initialize"))
		return
	} else if err != nil {
		panic(fmt.Errorf("unexpected error: %v", err))
	}
	close(syncedCh)

	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh
}

func (c *CRDDiscoveryController) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *CRDDiscoveryController) processNextWorkItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.syncFn(key.(schema.GroupVersion))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	utilruntime.HandleError(fmt.Errorf("key %v failed with %v", key, err))
	c.queue.AddRateLimited(key)

	return true

}

func (c *CRDDiscoveryController) addCustomResourceDefinition(obj interface{}) {
	castObj := obj.(*corev1.CustomResourceDefinition)
	logrus.Infof("adding customresourcedefinition %s", castObj.Name)
	c.enqueue(castObj)
}

func (c *CRDDiscoveryController) updateCustomResourceDefinition(oldObj, newObj interface{}) {
	castNewObj := newObj.(*corev1.CustomResourceDefinition)
	castOldObj := oldObj.(*corev1.CustomResourceDefinition)
	logrus.Tracef("updating customresourcedefinition %s", castOldObj.Name)
	c.enqueue(castNewObj)
	c.enqueue(castOldObj)
}

func (c *CRDDiscoveryController) deleteCustomResourceDefinition(obj interface{}) {
	castObj, ok := obj.(*corev1.CustomResourceDefinition)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			logrus.Errorf("couldn't get object from tombstone: %#v", obj)
			return
		}
		castObj, ok = tombstone.Obj.(*corev1.CustomResourceDefinition)
		if !ok {
			logrus.Errorf("tombstone contained object that is not expected: %#v", castObj)
			return
		}
	}
	logrus.Infof("deleting customresourcedefinition: %s", castObj.Name)
	c.enqueue(castObj)
}

func (c *CRDDiscoveryController) enqueue(obj *corev1.CustomResourceDefinition) {
	for _, v := range obj.Spec.Versions {
		c.queue.Add(schema.GroupVersion{Group: obj.Spec.Group, Version: v.Name})
	}
}
