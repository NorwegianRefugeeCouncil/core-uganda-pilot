package server

import "net/http"

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
