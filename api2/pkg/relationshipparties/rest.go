package relationshipparties

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	partiesStore *PartiesStore
}

func NewHandler(partiesStore *PartiesStore) *Handler {
	return &Handler{
		partiesStore: partiesStore,
	}
}

type PickPartyOptions struct {
	PartyTypeID string
	SearchParam string
}

func (h *Handler) PickParty(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	qry := req.URL.Query()

	pickPartyOptions := PickPartyOptions{
		PartyTypeID: qry.Get("partyTypeId"),
		SearchParam: qry.Get("searchParam"),
	}

	ret, err := h.partiesStore.FilteredList(ctx, pickPartyOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)

}
