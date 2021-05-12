package testing

import (
	"github.com/nrc-no/core/apps/api/pkg/tools/cache"
	"github.com/stretchr/testify/assert"
	"time"
)

func (s *MainTestSuite) TestInformer() {

	informer := s.informer.Core().V1().FormDefinitions().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			s.T().Logf("Added")
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			s.T().Logf("Updated")
		},
		DeleteFunc: func(obj interface{}) {
			s.T().Logf("Deleted")
		},
	})

	stopCh := make(chan struct{})
	go informer.Run(stopCh)

	_, err := s.nrcClient.CoreV1().FormDefinitions().Create(s.ctx, aValidFormDefinition("testgroup", "TestKind"))
	if !assert.NoError(s.T(), err) {
		stopCh <- struct{}{}
		return
	}

	time.Sleep(10 * time.Second)
	stopCh <- struct{}{}

}
