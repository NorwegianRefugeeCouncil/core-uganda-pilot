package webapp

import (
	"context"
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/apps/iam"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func (h *Server) RelationshipTypes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	r := &iam.RelationshipType{}

	if req.Method == "POST" {
		h.PostRelationshipType(ctx, r, w, req)
		return
	}

	relationshipTypes, err := h.iam.RelationshipTypes().List(ctx, iam.RelationshipTypeListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "relationshiptypes", map[string]interface{}{
		"RelationshipTypes": relationshipTypes,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Server) NewRelationshipType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	p, err := h.iam.PartyTypes().List(ctx, iam.PartyTypeListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "relationshiptype", map[string]interface{}{
		"PartyTypes": p,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Server) RelationshipType(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id found in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r, err := h.iam.RelationshipTypes().Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p, err := h.iam.PartyTypes().List(ctx, iam.PartyTypeListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostRelationshipType(ctx, r, w, req)
		return
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "relationshiptype", map[string]interface{}{
		"RelationshipType": r,
		"PartyTypes":       p,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Server) PostRelationshipType(
	ctx context.Context,
	r *iam.RelationshipType,
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

	r.Rules = []iam.RelationshipTypeRule{
		{
			PartyTypeRule: &iam.PartyTypeRule{
				FirstPartyTypeID:  formValues.Get("rules[0].firstPartyTypeId"),
				SecondPartyTypeID: formValues.Get("rules[0].secondPartyTypeId"),
			},
		},
	}

	if isNew {
		out, err := h.iam.RelationshipTypes().Create(ctx, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/settings/relationshiptypes/"+out.ID)
		w.WriteHeader(http.StatusSeeOther)
	} else {
		out, err := h.iam.RelationshipTypes().Update(ctx, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/settings/relationshiptypes/"+out.ID)
		w.WriteHeader(http.StatusSeeOther)
	}

}
