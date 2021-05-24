package server

import (
	"fmt"
	restclient "github.com/nrc-no/core/api/pkg/client/rest"
	"github.com/sirupsen/logrus"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"runtime/debug"
)

// PostStartHookContext provides information about this API server to a PostStartHookFunc
type PostStartHookContext struct {
	// LoopbackClientConfig is a config for a privileged loopback connection to the API server
	LoopbackClientConfig *restclient.Config
	// StopCh is the channel that will be closed when the server stops
	StopCh <-chan struct{}
}

type PostStartHookFunc func(context PostStartHookContext) error

type postStartHookEntry struct {
	hook PostStartHookFunc
	// originatingStack holds the stack that registered postStartHooks. This allows us to show a more helpful message
	// for duplicate registration.
	originatingStack string

	// done will be closed when the postHook is finished
	done chan struct{}
}

func (s *Server) AddPostStartHook(name string, hook PostStartHookFunc) error {
	if len(name) == 0 {
		return fmt.Errorf("missing name")
	}
	if hook == nil {
		return fmt.Errorf("hook cannot be nil")
	}
	s.postStartHookLock.Lock()
	defer s.postStartHookLock.Unlock()

	if s.postStartHookCalled {
		return fmt.Errorf("unable to add %s because PostStartHooks have already been called", name)
	}
	if postStartHook, exists := s.postStartHooks[name]; exists {
		return fmt.Errorf("unable to add PostStartHook %s because it was already registered by: %s", postStartHook.originatingStack)
	}

	done := make(chan struct{})
	s.postStartHooks[name] = postStartHookEntry{
		hook:             hook,
		originatingStack: string(debug.Stack()),
		done:             done,
	}

	return nil

}

func (s *Server) RunPostStartHooks(stopCh <-chan struct{}) {
	s.postStartHookLock.Lock()
	defer s.postStartHookLock.Unlock()
	s.postStartHookCalled = true

	context := PostStartHookContext{
		LoopbackClientConfig: s.LoopbackClientConfig,
		StopCh:               stopCh,
	}

	for hookName, hookEntry := range s.postStartHooks {
		go runPostStartHook(hookName, hookEntry, context)
	}

}

func (s *Server) AddPostStartHookOrDie(name string, hook PostStartHookFunc) {
	if err := s.AddPostStartHook(name, hook); err != nil {
		logrus.Fatalf("error registering PostStartHook %s: %v", name, err)
	}
}

func runPostStartHook(name string, entry postStartHookEntry, context PostStartHookContext) {
	var err error
	func() {
		defer utilruntime.HandleCrash()
		err = entry.hook(context)
	}()
	if err != nil {
		logrus.Fatalf("PostStartHook %v failed: %v", name, err)
	}
	close(entry.done)
}
