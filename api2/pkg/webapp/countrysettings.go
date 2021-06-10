package webapp

import "net/http"

func (h *Handler) CountrySettings(w http.ResponseWriter, req *http.Request) {
	if err := h.template.ExecuteTemplate(w, "countrysettings", map[string]interface{}{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
