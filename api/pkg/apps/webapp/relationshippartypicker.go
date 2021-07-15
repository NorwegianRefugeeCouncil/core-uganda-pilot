package webapp

import (
	"encoding/json"
	"github.com/nrc-no/core/pkg/apps/iam"
	"net/http"
)

type PickedParty struct {
	Party       iam.Party `json:"party"`
	DisplayName string    `json:"displayName"`
}

func (h *Server) PickRelationshipParty(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	iamClient := h.IAMClient(ctx)

	var listOptions iam.PartyListOptions
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := iamClient.Parties().List(ctx, listOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responseList []PickedParty
	for _, party := range response.Items {
		responseList = append(responseList, PickedParty{
			*party,
			party.String(),
		})
	}

	responseBytes, err := json.Marshal(responseList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}
