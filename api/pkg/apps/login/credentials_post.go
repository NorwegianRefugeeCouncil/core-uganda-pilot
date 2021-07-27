package login

import (
	"net/http"
)

func (s *Server) PostCredentials(w http.ResponseWriter, req *http.Request) {
	saltedHash, err := HashAndSalt(s.BCryptCost, []byte(req.Form.Get("plaintextPassword")))
	if err != nil {
		s.Error(w, err)
		return
	}

	var newCredential = Credential{
		PartyID: req.Form.Get("partyId"),
		Hash:    saltedHash,
	}

	_, err = s.Collection.InsertOne(req.Context(), newCredential)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.json(w, http.StatusOK, newCredential)
}
