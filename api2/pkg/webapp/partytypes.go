package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/apps/iam"
	"net/http"
)

func (h *Handler) PartyTypes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	if req.Method == "POST" {
		h.PostPartyType(ctx, &iam.PartyType{}, w, req)
		return
	}

	partyTypes, err := h.iam.PartyTypes().List(ctx, iam.PartyTypeListOptions{})
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

	var partyType = &iam.PartyType{}
	if id != "new" {
		var err error
		partyType, err = h.iam.PartyTypes().Get(ctx, id)
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
		"PartyType": partyType,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) PostPartyType(
	ctx context.Context,
	partyType *iam.PartyType,
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
		_, err := h.iam.PartyTypes().Create(ctx, partyType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/settings/partytypes/")
		w.WriteHeader(http.StatusSeeOther)
		return
	} else {
		_, err := h.iam.PartyTypes().Update(ctx, partyType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/settings/partytypes")
		w.WriteHeader(http.StatusSeeOther)
		return
	}
}
