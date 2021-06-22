package relationshipparties

import (
	"encoding/json"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
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

func (h *Handler) PickParty(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	qry := req.URL.Query()

	pickPartyOptions := parties.ListOptions{
		PartyTypeID: qry.Get("partyTypeId"),
		SearchParam: qry.Get("searchParam"),
	}

	ret, err := h.partiesStore.store.List(ctx, pickPartyOptions)
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
