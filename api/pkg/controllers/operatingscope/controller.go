package operatingscope

import (
	"fmt"
	v1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/api/pkg/auth/keycloak"
	informers "github.com/nrc-no/core/api/pkg/client/informers/core/v1"
	listers "github.com/nrc-no/core/api/pkg/client/listers/core/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/labels"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type OperatingScopeController struct {
	operatingScopeLister listers.OperatingScopeLister
	operatingScopeSynced cache.InformerSynced
	keycloakClient       *keycloak.KeycloakClient
	realmName            string
	syncFn               func(name string) error
	queue                workqueue.RateLimitingInterface
}

func NewOperatingScopeController(
	informer informers.OperatingScopeInformer,
	keycloakClient *keycloak.KeycloakClient,
	realmName string) *OperatingScopeController {
	c := &OperatingScopeController{
		operatingScopeLister: informer.Lister(),
		operatingScopeSynced: informer.Informer().HasSynced,
		queue:                workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "OperatingScopeController"),
		realmName:            realmName,
		keycloakClient:       keycloakClient,
	}

	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			castObj := obj.(*v1.OrganizationScope)
			logrus.Tracef("creating OrganizationScope %s", castObj.Name)
			c.queue.Add(castObj.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			castObj := newObj.(*v1.OrganizationScope)
			logrus.Tracef("updating OrganizationScope %s", castObj.Name)
			c.queue.Add(castObj.Name)
		},
		DeleteFunc: func(obj interface{}) {
			castObj, ok := obj.(*v1.OrganizationScope)
			if !ok {
				tombStone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					logrus.Errorf("couldn't get object from tombstone: %#v", obj)
					return
				}
				castObj, ok = tombStone.Obj.(*v1.OrganizationScope)
				if !ok {
					logrus.Errorf("tombstone contained object thtat is not expected: %#v", tombStone.Obj)
					return
				}
			}
			logrus.Tracef("deleting OrganizationScope %s", castObj.Name)
			c.queue.Add(castObj.Name)
		},
	})

	c.syncFn = c.sync
	return c
}

func (c *OperatingScopeController) Run(stopCh <-chan struct{}, syncedCh chan struct{}) {

	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	logrus.Infof("starting OperatingScopeController")
	defer logrus.Infof("shutting down OperatingScopeController")

	if !cache.WaitForCacheSync(stopCh, c.operatingScopeSynced) {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

	if err := wait.PollImmediateUntil(time.Second*10, func() (done bool, err error) {
		operatingScopes, err := c.operatingScopeLister.List(labels.Everything())
		if err != nil {
			utilruntime.HandleError(fmt.Errorf("failed to initially list operatingScopes: %v", err))
			return false, nil
		}
		for _, operatingScope := range operatingScopes {
			if err := c.syncFn(operatingScope.Name); err != nil {
				utilruntime.HandleError(fmt.Errorf("failed to initially sync formdefinition %s: %v", operatingScope.Name, err))
				return false, nil
			}
		}
		return true, nil

	}, stopCh); err == wait.ErrWaitTimeout {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for formdefinition to sync: %v", err))
		return
	} else if err != nil {
		panic(fmt.Errorf("unexpected error: %v", err))
	}

	close(syncedCh)

	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh

}

func (c *OperatingScopeController) processNextWorkItem() bool {
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

	utilruntime.HandleError(fmt.Errorf("failed to sync operatingScope: %v", err))
	c.queue.AddRateLimited(key)

	return true
}

func (c *OperatingScopeController) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *OperatingScopeController) sync(name string) error {

	operatingScope, err := c.operatingScopeLister.Get(name)
	if err != nil {
		return err
	}

	token, err := c.keycloakClient.GetToken(c.realmName)
	if err != nil {
		return err
	}

	clientScopes, err := c.keycloakClient.ListClientScopes(token, c.realmName)
	if err != nil {
		return err
	}

	var found *keycloak.ClientScope
	for _, scope := range clientScopes {
		if scope.Name == operatingScope.Name {
			found = scope
			break
		}
	}

	if found == nil {

		// must create keycloak client scope
		logrus.Infof("creating keycloak scope: %s", operatingScope.Name)
		clientScope := mapOperatingScopeToKeycloakClientScope(operatingScope)
		_, err := c.keycloakClient.CreateClientScope(token, c.realmName, clientScope)
		if err != nil {
			return err
		}

	} else {

		// must reconcile client scope
		// Since we cannot simply update the Client protocol mappers using a PUT
		// on the Client Scope (we must reconcile the individual protocol mappers)
		//
		// We will construct a map of the protocol mappers we have to
		// 1. Delete
		// 2. Update
		// 3. Create
		//
		// Then will process each individually

		// This is the desired Client Scope with the mappers contained in it
		desiredClientScope := mapOperatingScopeToKeycloakClientScope(operatingScope)

		if !clientScopeEqual(desiredClientScope, found) {
			logrus.Infof("updating keycloak ClientScope: %s", operatingScope.Name)
			desiredClientScope.ID = found.ID
			_, err := c.keycloakClient.UpdateClientScope(token, c.realmName, desiredClientScope)
			if err != nil {
				return err
			}
		}

		desiredMapperMap := map[string]*keycloak.ProtocolMapper{}
		foundMapperMap := map[string]*keycloak.ProtocolMapper{}
		toDeleteMap := map[string]bool{}
		toCreateMap := map[string]bool{}
		toUpdateMap := map[string]bool{}

		for _, foundMapper := range found.ProtocolMappers {
			foundMapperMap[foundMapper.Name] = foundMapper
			toDeleteMap[foundMapper.Name] = true
		}

		for _, desiredMapper := range desiredClientScope.ProtocolMappers {
			desiredMapperMap[desiredMapper.Name] = desiredMapper
			toCreateMap[desiredMapper.Name] = true
			if _, ok := toDeleteMap[desiredMapper.Name]; ok {
				toDeleteMap[desiredMapper.Name] = false
				toUpdateMap[desiredMapper.Name] = true
			}
		}

		for _, foundMapper := range found.ProtocolMappers {
			if _, ok := toDeleteMap[foundMapper.Name]; ok {
				toCreateMap[foundMapper.Name] = false
			}
		}

		for toDeleteName, toDelete := range toDeleteMap {
			if toDelete {

				logrus.Tracef("deleting ProtocolMapper with name '%s' for ClientScope %s",
					toDeleteName, found.Name)

				if err := c.keycloakClient.DeleteProtocolMapper(token, c.realmName, found.ID, foundMapperMap[toDeleteName].ID); err != nil {
					return err
				}
			}
		}

		for toCreateName, shouldCreate := range toCreateMap {
			if shouldCreate {

				logrus.Tracef("creating ProtocolMapper with name '%s' for ClientScope %s",
					toCreateName, found.Name)

				if _, err := c.keycloakClient.CreateProtocolMapper(token, c.realmName, found.ID, desiredMapperMap[toCreateName]); err != nil {
					return err
				}
			}
		}

		for toUpdateName, shouldUpdate := range toUpdateMap {
			if shouldUpdate {
				foundMapper := foundMapperMap[toUpdateName]
				desired := desiredMapperMap[toUpdateName]
				desired.ID = foundMapper.ID

				if !protocolMapperEqual(foundMapper, desired) {

					logrus.Tracef("updating ProtocolMapper with name '%s' for ClientScope %s",
						desired.Name,
						found.Name)

					if _, err := c.keycloakClient.UpdateProtocolMapper(token, c.realmName, found.ID, desired); err != nil {
						return err
					}
				}

			}
		}

	}
	return nil
}

func clientScopeEqual(a, b *keycloak.ClientScope) bool {
	if a.Name != b.Name {
		return false
	}
	if a.Description != b.Description {
		return false
	}
	if a.Protocol != b.Protocol {
		return false
	}
	return true
}

func protocolMapperEqual(a, b *keycloak.ProtocolMapper) bool {
	if a.ProtocolMapper != b.ProtocolMapper {
		return false
	}
	if a.Name != b.Name {
		return false
	}
	if a.Protocol != b.Protocol {
		return false
	}
	if len(a.Configuration) != len(b.Configuration) {
		return false
	}
	for key, configA := range a.Configuration {
		configB, ok := b.Configuration[key]
		if !ok {
			return false
		}
		if configA != configB {
			return false
		}
	}
	return true
}

func mapOperatingScopeToKeycloakClientScope(operatingScope *v1.OrganizationScope) *keycloak.ClientScope {
	scope := &keycloak.ClientScope{
		Name:        operatingScope.Name,
		Protocol:    "openid-connect",
		Description: fmt.Sprintf("OpenID scope for %s", operatingScope.Name),
	}

	for _, additionalInfo := range operatingScope.Spec.AdditionalBeneficiaryInformation {
		protocolMapper := &keycloak.ProtocolMapper{
			Name:           additionalInfo.Key,
			Protocol:       "openid-connect",
			ProtocolMapper: "oidc-usermodel-attribute-mapper",
			Configuration: map[string]string{
				"access.token.claim":   "true",
				"claim.name":           "core.nrc.no/" + additionalInfo.Key,
				"id.token.claim":       "true",
				"jsonType.label":       "String",
				"user.attribute":       "core.nrc.no/" + additionalInfo.Key,
				"userinfo.token.claim": "true",
			},
		}
		scope.ProtocolMappers = append(scope.ProtocolMappers, protocolMapper)
	}

	return scope
}
