package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/parties/api"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"github.com/nrc-no/core-kafka/pkg/parties/relationshiptypes"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (h *Handler) RelationshipTypes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	cli := relationshiptypes.NewClient("http://localhost:9000")
	r := &api.RelationshipType{}

	if req.Method == "POST" {
		h.PostRelationshipType(ctx, cli, r, w, req)
		return
	}

	relationshipTypes, err := cli.List(ctx)
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
	relationshipTypesCli := relationshiptypes.NewClient("http://localhost:9000")
	partyTypesCli := partytypes.NewClient("http://localhost:9000")

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id found in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r, err := relationshipTypesCli.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p, err := partyTypesCli.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostRelationshipType(ctx, relationshipTypesCli, r, w, req)
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
	cli *relationshiptypes.Client,
	r *api.RelationshipType,
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
	r.FirstPartyRole = formValues.Get("firstPartyRole")
	r.SecondPartyRole = formValues.Get("secondPartyRole")

	r.Rules = []api.RelationshipTypeRule{
		api.RelationshipTypeRule{
			api.PartyTypeRule{
				FirstPartyType:  formValues.Get("rules[0].firstPartyType"),
				SecondPartyType: formValues.Get("rules[0].secondPartyType"),
			},
		},
	}

	if isNew {
		out, err := cli.Create(ctx, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/settings/relationshiptypes/"+out.ID)
		w.WriteHeader(http.StatusSeeOther)
	} else {
		out, err := cli.Update(ctx, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/settings/relationshiptypes/"+out.ID)
		w.WriteHeader(http.StatusSeeOther)
	}

}
