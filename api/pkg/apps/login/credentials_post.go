package login

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func (s *Server) PostCredentials(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload SetCredentialPayload
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	res, err := s.Collection.Find(ctx, bson.M{"partyId": payload.PartyID})
	if err != nil {
		s.Error(w, err)
		return
	}

	var credentials []*Credential
	for {
		if !res.Next(ctx) {
			break
		}
		var c Credential
		if err := res.Decode(&c); err != nil {
			s.Error(w, err)
			return
		}
		credentials = append(credentials, &c)
	}

	saltedHash, err := HashAndSalt(s.BCryptCost, []byte(payload.PlaintextPassword))
	if err != nil {
		s.Error(w, err)
		return
	}

	var returnValue Credential

	// no credentials found, need to create new
	if len(credentials) == 0 {
		var newCredential = Credential{
			PartyID: payload.PartyID,
			Hash:    saltedHash,
		}

		_, err = s.Collection.InsertOne(ctx, newCredential)
		if err != nil {
			s.Error(w, err)
			return
		}

		returnValue = newCredential
	}

	// existing credential found, updating
	if len(credentials) == 1 {
		err = s.SetPassword(ctx, payload.PartyID, payload.PlaintextPassword)

		if err != nil {
			s.Error(w, err)
			return
		}
	}

	// more than one credential found, error
	if len(credentials) > 1 {
		s.Error(w, errors.New("more than one credential found for user"))
		return
	}

	s.json(w, http.StatusOK, returnValue)
}
