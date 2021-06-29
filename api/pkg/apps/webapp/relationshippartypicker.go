package webapp

import (
	"encoding/json"
	"github.com/nrc-no/core/pkg/apps/iam"
	"net/http"
)

func (h *Server) PickRelationshipParty(w http.ResponseWriter, req *http.Request){
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

	responseBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}
