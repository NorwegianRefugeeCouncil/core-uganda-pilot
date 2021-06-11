package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/cases/api"
	"github.com/nrc-no/core-kafka/pkg/cases/cases"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func (h *Handler) Cases(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	cli := cases.NewClient("http://localhost:9000")

	if req.Method == "POST" {
		h.PostCase(ctx, cli, "", w, req)
		return
	}

	list, err := cli.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.template.ExecuteTemplate(w, "cases", map[string]interface{}{
		"Cases": list,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) Case(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	caseClient := cases.NewClient("http://localhost:9000")

	var c *api.Case
	var cList *api.CaseList
	var cTypes *api.CaseTypeList
	var pTypes *api.PartyTypeList
	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if id == "new" {
			c = &api.Case{}
			return nil
		}
		var err error
		c, err = caseClient.Get(waitCtx, id)
		return err
	})

	g.Go(func() error {
		var err error
		cList, err = caseClient.List(waitCtx)
		return err
	})

	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostCase(ctx, caseClient, id, w, req)
		return
	}

	if err := h.template.ExecuteTemplate(w, "case", map[string]interface{}{
		"Case":      c,
		"CaseTypes": cTypes,
		"PartyType": pTypes,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) PostCase(ctx context.Context, cli *cases.Client, id string, w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var kase *api.Case

	f := req.Form
	for key, value := range f {
		switch key {
		case "caseType":
			kase.CaseTypeID = value[0]
		case "partyId":
			kase.PartyID = value[0]
		case "description":
			kase.Description = value[0]
		}
	}

	if id == "" {
		_, err := cli.Create(ctx, &api.Case{
			ID:          uuid.NewV4().String(),
			CaseTypeID:  kase.CaseTypeID,
			PartyID:     kase.PartyID,
			Description: kase.Description,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		_, err := cli.Update(ctx, &api.Case{
			ID:          id,
			CaseTypeID:  kase.CaseTypeID,
			PartyID:     kase.PartyID,
			Description: kase.Description,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Location", "/cases/"+id)
	w.WriteHeader(http.StatusSeeOther)

}
