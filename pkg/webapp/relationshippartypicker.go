package webapp

import (
	"encoding/json"
	iam2 "github.com/nrc-no/core/pkg/iam"
	"net/http"
)

type PickedParty struct {
	Party       iam2.Party `json:"party"`
	DisplayName string     `json:"displayName"`
}

func (s *Server) PickRelationshipParty(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	var listOptions iam2.PartyListOptions
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.Error(w, err)
		return
	}

	response, err := iamClient.Parties().List(ctx, listOptions)
	if err != nil {
		s.Error(w, err)
		return
	}

	var responseList []PickedParty
	for _, party := range response.Items {
		responseList = append(responseList, PickedParty{
			Party:       *party,
			DisplayName: party.String(),
		})
	}

	responseBytes, err := json.Marshal(responseList)
	if err != nil {
		s.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}
