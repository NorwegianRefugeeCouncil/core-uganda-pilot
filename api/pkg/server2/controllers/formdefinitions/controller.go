package formdefinitions

import (
	"context"
	"github.com/nrc-no/core/api/pkg/apis/core/helpers"
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	v1 "github.com/nrc-no/core/api/pkg/client/core/v1"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"reflect"
	"time"
)

// FormDefinitionController will
// 1. Create a CustomResourceDefinition for a given FormDefinition
// 2. Delete a CustomResourceDefinition if the FormDefinition was deleted
// 3. Update a CustomResourceDefinition if the FormDefinition was updated
// It will poll every 15 seconds to reconcile the above.
type FormDefinitionController struct {
	cli v1.CoreV1Interface
}

func NewFormDefinitionController(ctx context.Context, cli v1.CoreV1Interface) *FormDefinitionController {
	controller := &FormDefinitionController{
		cli: cli,
	}

	go func() {
		if err := wait.PollInfiniteWithContext(ctx, time.Minute, controller.syncFormDefinitions); err != nil {
			logrus.Errorf("failed to sync form definitions: %v", err)
			return
		}
	}()

	go func() {
		// TODO: find a way to run this on post startup hooks
		// this would fail otherwise because the server is not up yet
		// at this point
		time.Sleep(5 * time.Second)
		_, _ = controller.syncFormDefinitions(ctx)
	}()

	return controller
}

func (c *FormDefinitionController) syncFormDefinitions(ctx context.Context) (done bool, err error) {

	g, ctx1 := errgroup.WithContext(ctx)

	formDefNames := sets.String{}
	formDefMap := map[string]*corev1.FormDefinition{}
	g.Go(func() error {
		// Retrieve form definitions
		formDefs, err := c.cli.FormDefinitions().List(ctx1, metav1.ListOptions{})
		if err != nil {
			return err
		}

		// Build map of name -> formDefinition
		for _, formDef := range formDefs.Items {
			formDefMap[formDef.Name] = &formDef
			formDefNames.Insert(formDef.Name)
		}
		return nil
	})

	crdNames := sets.String{}
	crdMap := map[string]*corev1.CustomResourceDefinition{}
	g.Go(func() error {
		// Retrieve CRDs
		crds, err := c.cli.CustomResourceDefinitions().List(ctx1, metav1.ListOptions{})
		if err != nil {
			return err
		}

		// Build map of name -> crd
		for _, crd := range crds.Items {
			crdMap[crd.Name] = &crd
			crdNames.Insert(crd.Name)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		logrus.Warnf("errot while retrieving form definitions: %v", err)
		return false, nil
	}

	crdsToRemove := crdNames.Difference(formDefNames)
	crdsToAdd := formDefNames.Difference(crdNames)
	alreadyPresentCrds := formDefNames.Intersection(crdNames)

	g, ctx2 := errgroup.WithContext(ctx)
	for crdName, _ := range crdsToRemove {
		crdNameToRemove := crdName
		g.Go(func() error {
			logrus.Infof("deleting custom resource: %s", crdNameToRemove)
			err := c.cli.CustomResourceDefinitions().Delete(ctx2, crdNameToRemove, metav1.DeleteOptions{})
			return err
		})
	}
	for crdName, _ := range crdsToAdd {
		crdNameToAdd := crdName
		formDef := formDefMap[crdNameToAdd]
		g.Go(func() error {
			logrus.Infof("creating custom resource: %s", crdNameToAdd)
			crd := helpers.ConvertToCustomResourceDefinition(formDef)
			_, err := c.cli.CustomResourceDefinitions().Create(ctx2, crd, metav1.CreateOptions{})
			return err
		})
	}
	for crdName, _ := range alreadyPresentCrds {
		crdNameToAdd := crdName
		formDef := formDefMap[crdNameToAdd]
		actualCrd := crdMap[crdNameToAdd]
		desiredCrd := helpers.ConvertToCustomResourceDefinition(formDef)
		if reflect.DeepEqual(actualCrd.Spec, desiredCrd.Spec) {
			continue
		}
		g.Go(func() error {
			logrus.Infof("updating custom resource: %s", crdNameToAdd)
			_, err := c.cli.CustomResourceDefinitions().Update(ctx2, desiredCrd, metav1.UpdateOptions{})
			return err
		})
	}

	if err := g.Wait(); err != nil {
		logrus.Warnf("errot while removing/adding custom resource definitions: %v", err)
		return false, nil
	}

	return false, nil
}
