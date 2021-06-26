package middleware

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
	"net/http"
	"time"
)

func UseLogging() func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()

			stWriter := &statusWriter{w: w}
			handler.ServeHTTP(stWriter, req)

			end := time.Now()

			statusCode := stWriter.statusCode
			if stWriter.statusCode == 0 {
				statusCode = 200
			}

			fields := logrus.Fields{
				"method":     req.Method,
				"statusCode": statusCode,
				"path":       req.URL.Path,
				"responseMs": math.Round(float64(end.Sub(start).Nanoseconds())/1000000.0*100.0) / 100.0,
			}

			if stWriter.statusCode < 400 {
				logrus.WithFields(fields).Infof("")
			} else {
				logrus.WithFields(fields).
					WithError(fmt.Errorf("inbound request failed with status code: %d", statusCode)).
					Errorf("")
			}
		})
	}
}

type statusWriter struct {
	w          http.ResponseWriter
	statusCode int
}

var _ http.ResponseWriter = &statusWriter{}

func (s *statusWriter) Header() http.Header {
	return s.w.Header()
}

func (s *statusWriter) Write(bytes []byte) (int, error) {
	return s.w.Write(bytes)
}

func (s *statusWriter) WriteHeader(statusCode int) {
	s.statusCode = statusCode
	s.w.WriteHeader(statusCode)
}
