package login

import (
	"net/http"
)

func (s *Server) PostCredentials(w http.ResponseWriter, req *http.Request) {
	var payload SetCredentialPayload
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	saltedHash, err := HashAndSalt(s.BCryptCost, []byte(payload.PlaintextPassword))
	if err != nil {
		s.Error(w, err)
		return
	}

	var newCredential = Credential{
		PartyID: payload.PartyID,
		Hash:    saltedHash,
	}

	_, err = s.Collection.InsertOne(req.Context(), newCredential)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.json(w, http.StatusOK, newCredential)
}
