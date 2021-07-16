package webapp

import (
	"net/http"
)

func (s *Server) Settings(w http.ResponseWriter, req *http.Request) {
	if err := s.renderFactory.New(req).ExecuteTemplate(w, "settings", nil); err != nil {
		s.Error(w, err)
		return
	}
}
