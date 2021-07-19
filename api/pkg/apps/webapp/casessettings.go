package webapp

import (
	"net/http"
)

func (s *Server) CasesSettings(w http.ResponseWriter, req *http.Request) {
	if err := s.renderFactory.New(req).ExecuteTemplate(w, "casessettings", nil); err != nil {
		s.Error(w, err)
		return
	}
}
