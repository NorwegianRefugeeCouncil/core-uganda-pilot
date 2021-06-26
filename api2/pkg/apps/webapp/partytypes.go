package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/apps/iam"
	"net/http"
)

func (h *Server) PartyTypes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	iamClient := h.IAMClient(ctx)

	if req.Method == "POST" {
		h.PostPartyType(ctx, &iam.PartyType{}, w, req)
		return
	}

	partyTypes, err := iamClient.PartyTypes().List(ctx, iam.PartyTypeListOptions{})
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

func (h *Server) PartyType(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	iamClient := h.IAMClient(ctx)

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id found in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var partyType = &iam.PartyType{}
	if id != "new" {
		var err error
		partyType, err = iamClient.PartyTypes().Get(ctx, id)
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

func (h *Server) PostPartyType(
	ctx context.Context,
	partyType *iam.PartyType,
	w http.ResponseWriter,
	req *http.Request,
) {

	iamClient := h.IAMClient(ctx)

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
		created, err := iamClient.PartyTypes().Create(ctx, partyType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/settings/partytypes/"+created.ID)
		w.WriteHeader(http.StatusSeeOther)
		return
	} else {
		updated, err := iamClient.PartyTypes().Update(ctx, partyType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/settings/partytypes/"+updated.ID)
		w.WriteHeader(http.StatusSeeOther)
		return
	}
}
