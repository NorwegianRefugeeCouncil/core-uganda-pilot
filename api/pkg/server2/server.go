package server2

import (
	"context"
	"errors"
	"github.com/emicklei/go-restful"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	listenAddress      string
	handler            *APIServerHandler
	goRestfulContainer *restful.Container
}

func (s *Server) Run(stopCh <-chan struct{}) error {

	httpServer := http.Server{
		Addr:    s.listenAddress,
		Handler: s.handler.FullHandlerChain,
	}

	go func() {
		err := httpServer.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			logrus.Info("server shutting down")
			err = nil
		}

		select {
		case <-stopCh:
			logrus.Info("stopped listening")
		default:
			logrus.Errorf("stopped listening because of: %v", err)
		}
	}()

	<-stopCh

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logrus.Errorf("server shutdown failed: %v", err)
		return err
	}

	return nil

}
