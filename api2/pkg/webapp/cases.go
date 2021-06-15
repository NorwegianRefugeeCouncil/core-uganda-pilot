package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	casesapi "github.com/nrc-no/core-kafka/pkg/cases/api"
	"github.com/nrc-no/core-kafka/pkg/cases/cases"
	"github.com/nrc-no/core-kafka/pkg/cases/casetypes"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func (h *Handler) Cases(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	caseClient := cases.NewClient("http://localhost:9000")

	list, err := caseClient.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostCase(ctx, caseClient, "", w, req)
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
	caseTypesClient := casetypes.NewClient("http://localhost:9000")
	partyClient := parties.NewClient("http://localhost:9000")

	var kase *casesapi.Case
	var kaseTypes *casesapi.CaseTypeList

	var party *parties.Party
	var partyList *parties.PartyList

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if id == "new" {
			kase = &casesapi.Case{}
			return nil
		}
		var err error
		kase, err = caseClient.Get(waitCtx, id)
		return err
	})

	g.Go(func() error {
		var err error
		kaseTypes, err = caseTypesClient.List(waitCtx)
		return err
	})

	g.Go(func() error {
		var err error
		party, err = partyClient.Get(waitCtx, id)
		return err
	})

	listOptions := &parties.ListOptions{} //TODO

	g.Go(func() error {
		var err error
		partyList, err = partyClient.List(waitCtx, *listOptions)
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
		"Case":      kase,
		"CaseTypes": kaseTypes,
		"Party":     party,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) NewCase(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	caseTypesClient := casetypes.NewClient("http://localhost:9000")
	partiesClient := parties.NewClient("http://localhost:9000")

	var caseTypes *casesapi.CaseTypeList
	var p *parties.PartyList

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		caseTypes, err = caseTypesClient.List(waitCtx)
		return err
	})

	listOptions := &parties.ListOptions{} //TODO
	g.Go(func() error {
		var err error
		p, err = partiesClient.List(waitCtx, *listOptions)
		return err
	})

	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.template.ExecuteTemplate(w, "casenew", map[string]interface{}{
		"CaseTypes": caseTypes,
		"Parties":   p,
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

	var kase *casesapi.Case

	f := req.Form
	for key, value := range f {
		switch key {
		case "caseTypeId":
			kase.CaseTypeID = value[0]
		case "partyId":
			kase.PartyID = value[0]
		case "description":
			kase.Description = value[0]
		}
	}

	if id == "" {
		_, err := cli.Create(ctx, &casesapi.Case{
			CaseTypeID:  kase.CaseTypeID,
			PartyID:     kase.PartyID,
			Description: kase.Description,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		_, err := cli.Update(ctx, &casesapi.Case{
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
	w.Header().Set("Location", "/cases")
	w.WriteHeader(http.StatusSeeOther)

}
