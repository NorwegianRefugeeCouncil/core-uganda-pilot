package webapp

import (
	"net/http"
)

func (h *Handler) CasesSettings(w http.ResponseWriter, req *http.Request) {
	if err := h.template.ExecuteTemplate(w, "casessettings", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
