package server

import (
	"context"
	"errors"
	"github.com/nrc-no/core/api/pkg/client/informers"
	"github.com/nrc-no/core/api/pkg/client/rest"
	"github.com/nrc-no/core/api/pkg/server/routes"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	listener net.Listener
	handler  *APIServerHandler
	ctx      context.Context

	postStartHooks       map[string]postStartHookEntry
	postStartHookLock    sync.Mutex
	postStartHookCalled  bool
	informers            informers.SharedInformerFactory
	LoopbackClientConfig *rest.Config
	listedPathProvider   routes.ListedPathProvider
}

func (s *Server) Run() error {

	httpServer := http.Server{
		Handler: s.handler.FullHandlerChain,
	}

	var errChan = make(chan error, 1)
	go func() {

		err := httpServer.Serve(s.listener)

		if errors.Is(err, http.ErrServerClosed) {
			logrus.Info("server shutting down")
			err = nil
		}

		select {
		case <-s.ctx.Done():
			logrus.Info("stopped listening")
		default:
			logrus.Errorf("stopped listening because of: %v", err)
			errChan <- err
		}
	}()

	stopCh := make(chan struct{})
	go func() {
		select {
		case <-s.ctx.Done():
			close(stopCh)
		}
	}()
	s.RunPostStartHooks(stopCh)

	select {
	case <-s.ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			logrus.Errorf("server shutdown failed: %v", err)
			return err
		}
		time.Sleep(5 * time.Second)
	case err := <-errChan:
		return err
	}

	return nil

}
