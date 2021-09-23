package login

import (
	"errors"
	"fmt"
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

	credentialsCollection, err := s.credentialsCollectionFn()
	if err != nil {
		err = fmt.Errorf("failed to get credentials collection: %v", err)
		s.Error(w, err)
	}

	res, err := credentialsCollection.Find(ctx, bson.M{"partyId": payload.PartyID})
	if err != nil {
		s.Error(w, fmt.Errorf("failed to get credentials: %v", err))
		return
	}

	var credentials []*Credential
	for {
		if !res.Next(ctx) {
			break
		}
		var c Credential
		if err := res.Decode(&c); err != nil {
			s.Error(w, fmt.Errorf("failed to decode credentials: %v", err))
			return
		}
		credentials = append(credentials, &c)
	}

	// no credentials found, need to create new
	if len(credentials) == 0 {
		err = s.CreatePassword(ctx, payload.PartyID, payload.PlaintextPassword)
		if err != nil {
			s.Error(w, fmt.Errorf("failed to create password: %v", err))
			return
		}
	}

	// existing credential found, updating
	if len(credentials) == 1 {
		err = s.SetPassword(ctx, payload.PartyID, payload.PlaintextPassword)
		if err != nil {
			s.Error(w, fmt.Errorf("failed to set password: %v", err))
			return
		}
	}

	// more than one credential found, error
	if len(credentials) > 1 {
		s.Error(w, errors.New("more than one credential found for user"))
		return
	}

	var i interface{}
	s.json(w, http.StatusOK, i)
}
