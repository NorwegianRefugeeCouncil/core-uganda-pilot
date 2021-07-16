package webapp

import (
	"encoding/json"
	"github.com/nrc-no/core/pkg/apps/iam"
	"net/http"
)

func (s *Server) PickRelationshipParty(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	iamClient := s.IAMClient(ctx)

	var listOptions iam.PartyListOptions
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.Error(w, err)
		return
	}

	response, err := iamClient.Parties().List(ctx, listOptions)
	if err != nil {
		s.Error(w, err)
		return
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		s.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}
