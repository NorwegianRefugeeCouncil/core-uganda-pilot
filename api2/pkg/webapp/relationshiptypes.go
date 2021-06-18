package webapp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"github.com/nrc-no/core-kafka/pkg/parties/relationshiptypes"
	uuid "github.com/satori/go.uuid"
)

func (h *Handler) RelationshipTypes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	r := &relationshiptypes.RelationshipType{}

	if req.Method == "POST" {
		h.PostRelationshipType(ctx, r, w, req)
		return
	}

	relationshipTypes, err := h.relationshipTypeClient.List(ctx, relationshiptypes.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.template.ExecuteTemplate(w, "relationshiptypes", map[string]interface{}{
		"RelationshipTypes": relationshipTypes,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) NewRelationshipType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	partyTypesCli := partytypes.NewClient("http://localhost:9000")

	p, err := partyTypesCli.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.template.ExecuteTemplate(w, "relationshiptype", map[string]interface{}{
		"PartyTypes": p,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RelationshipType(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id found in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r, err := h.relationshipTypeClient.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p, err := h.partyTypeClient.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostRelationshipType(ctx, r, w, req)
		return
	}

	if err := h.template.ExecuteTemplate(w, "relationshiptype", map[string]interface{}{
		"RelationshipType": r,
		"PartyTypes":       p,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PostRelationshipType(
	ctx context.Context,
	r *relationshiptypes.RelationshipType,
	w http.ResponseWriter,
	req *http.Request,
) {

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	formValues := req.Form

	isNew := false
	if len(r.ID) == 0 {
		r.ID = uuid.NewV4().String()
		isNew = true
	}

	r.Name = formValues.Get("name")

	if formValues.Get("isDirectional") == "true" {
		r.IsDirectional = true
	} else {
		r.IsDirectional = false
	}
	r.FirstPartyRole = formValues.Get("firstPartyRole")
	r.SecondPartyRole = formValues.Get("secondPartyRole")

	r.Rules = []relationshiptypes.RelationshipTypeRule{
		{
			relationshiptypes.PartyTypeRule{
				FirstPartyType:  formValues.Get("rules[0].firstPartyType"),
				SecondPartyType: formValues.Get("rules[0].secondPartyType"),
			},
		},
	}

	if isNew {
		out, err := h.relationshipTypeClient.Create(ctx, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/settings/relationshiptypes/"+out.ID)
		w.WriteHeader(http.StatusSeeOther)
	} else {
		out, err := h.relationshipTypeClient.Update(ctx, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/settings/relationshiptypes/"+out.ID)
		w.WriteHeader(http.StatusSeeOther)
	}

}
