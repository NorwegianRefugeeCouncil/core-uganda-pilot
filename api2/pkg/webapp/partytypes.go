package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"net/http"
)

func (h *Handler) PartyTypes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	if req.Method == "POST" {
		h.PostPartyType(ctx, &partytypes.PartyType{}, w, req)
		return
	}

	partyTypes, err := h.partyTypeClient.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "partytypes", map[string]interface{}{
		"PartyTypes": partyTypes,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) PartyType(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id found in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	partyTypesCli := partytypes.NewClient("http://localhost:9000")

	var partyType = &partytypes.PartyType{}
	if id != "new" {
		var err error
		partyType, err = partyTypesCli.Get(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if req.Method == "POST" {
		h.PostPartyType(ctx, partyType, w, req)
		return
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "partytype", map[string]interface{}{
		"PartyTypeID": partyType,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) PostPartyType(
	ctx context.Context,
	partyType *partytypes.PartyType,
	w http.ResponseWriter,
	req *http.Request,
) {

	isNew := false
	if len(partyType.ID) == 0 {
		isNew = true
	}

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := req.Form.Get("name")
	partyType.Name = name
	partyType.IsBuiltIn = req.Form.Get("isBuiltIn") == "true"

	if isNew {
		out, err := h.partyTypeClient.Create(ctx, partyType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/settings/partytypes/"+out.ID)
		w.WriteHeader(http.StatusSeeOther)
		return
	} else {
		_, err := h.partyTypeClient.Update(ctx, partyType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/settings/partytypes/"+partyType.ID)
		w.WriteHeader(http.StatusSeeOther)
		return
	}
}
