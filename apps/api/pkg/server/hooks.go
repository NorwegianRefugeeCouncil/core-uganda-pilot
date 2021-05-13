package server

import (
  "fmt"
  "k8s.io/klog/v2"
  "runtime/debug"
)

// AddPostStartHookOrDie allows you to add a PostStartHook, but dies on failure
func (s *Server) AddPostStartHookOrDie(name string, hook PostStartHookFunc) {
  if err := s.AddPostStartHook(name, hook); err != nil {
    klog.Fatalf("Error registering PostStartHook %q: %v", name, err)
  }
}
// AddPostStartHook allows you to add a PostStartHook.
func (s *Server) AddPostStartHook(name string, hook PostStartHookFunc) error {
  if len(name) == 0 {
    return fmt.Errorf("missing name")
  }
  if hook == nil {
    return fmt.Errorf("hook func may not be nil: %q", name)
  }
  if s.disabledPostStartHooks.Has(name) {
    klog.V(1).Infof("skipping %q because it was explicitly disabled", name)
    return nil
  }

  s.postStartHookLock.Lock()
  defer s.postStartHookLock.Unlock()

  if s.postStartHooksCalled {
    return fmt.Errorf("unable to add %q because PostStartHooks have already been called", name)
  }
  if postStartHook, exists := s.postStartHooks[name]; exists {
    // this is programmer error, but it can be hard to debug
    return fmt.Errorf("unable to add %q because it was already registered by: %s", name, postStartHook.originatingStack)
  }

  // done is closed when the poststarthook is finished.  This is used by the health check to be able to indicate
  // that the poststarthook is finished
  done := make(chan struct{})
  if err := s.AddBootSequenceHealthChecks(postStartHookHealthz{name: "poststarthook/" + name, done: done}); err != nil {
    return err
  }
  s.postStartHooks[name] = postStartHookEntry{hook: hook, originatingStack: string(debug.Stack()), done: done}

  return nil
}
