package server2

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
)

type Server struct {
	listener net.Listener
	handler  *APIServerHandler
	ctx      context.Context
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

	select {
	case <-s.ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			logrus.Errorf("server shutdown failed: %v", err)
			return err
		}
	case err := <-errChan:
		return err
	}

	return nil

}
