package webapp

import (
	"net/http"
)

func (h *Server) CasesSettings(w http.ResponseWriter, req *http.Request) {
	if err := h.renderFactory.New(req).ExecuteTemplate(w, "casessettings", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
