package webapp

import (
	"net/http"
)

func (h *Handler) Settings(w http.ResponseWriter, req *http.Request) {
	if err := h.renderFactory.New(req).ExecuteTemplate(w, "settings", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
