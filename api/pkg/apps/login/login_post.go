package login

import (
	"fmt"
	"net/http"
)

func (s *Server) PostLogin(w http.ResponseWriter, req *http.Request) {

	var payload SetCredentialPayload
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, fmt.Errorf("failed to bind request payload: %v", err))
		return
	}

	if err := s.SetPassword(req.Context(), payload.PartyID, payload.PlaintextPassword); err != nil {
		s.Error(w, fmt.Errorf("failed to set password: %v", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
