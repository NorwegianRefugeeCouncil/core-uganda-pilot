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

	kases, err := h.caseClient.List(ctx, cases.ListOptions{})
	caseTypes, err := h.caseTypeClient.List(ctx, casetypes.ListOptions{})
	partyList, err := h.partyClient.List(ctx, parties.ListOptions{})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostCase(ctx, "", w, req)
		return
	}

	if err := h.template.ExecuteTemplate(w, "cases", map[string]interface{}{
		"Cases":     kases,
		"CaseTypes": caseTypes,
		"Parties":   partyList,
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

	if req.Method == "POST" {
		h.PostCase(ctx, id, w, req)
		return
	}

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
		kase, err = h.caseClient.Get(waitCtx, id)
		return err
	})

	g.Go(func() error {
		var err error
		kaseTypes, err = h.caseTypeClient.List(waitCtx, casetypes.ListOptions{})
		return err
	})

	g.Go(func() error {
		var err error
		partyList, err = h.partyClient.List(waitCtx, parties.ListOptions{})
		return err
	})

	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	party, err := h.partyClient.Get(waitCtx, kase.PartyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.template.ExecuteTemplate(w, "case", map[string]interface{}{
		"Case":      kase,
		"CaseTypes": kaseTypes,
		"Party":     party,
		"Parties":   partyList,
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
		caseTypes, err = caseTypesClient.List(waitCtx, casetypes.ListOptions{})
		return err
	})

	listOptions := &parties.ListOptions{} // TODO
	g.Go(func() error {
		var err error
		p, err = partiesClient.List(waitCtx, *listOptions)
		return err
	})

	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	qry := req.URL.Query()

	if err := h.template.ExecuteTemplate(w, "casenew", map[string]interface{}{
		"PartyID":    qry.Get("partyId"),
		"CaseTypeID": qry.Get("caseTypeId"),
		"CaseTypes":  caseTypes,
		"Parties":    p,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PostCase(ctx context.Context, id string, w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	caseTypeId := req.Form.Get("caseTypeId")
	partyId := req.Form.Get("partyId")
	description := req.Form.Get("description")
	done := req.Form.Get("done-check")

	if id == "" {
		_, err := h.caseClient.Create(ctx, &casesapi.Case{
			CaseTypeID:  caseTypeId,
			PartyID:     partyId,
			Description: description,
			Done:        false,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		_, err := h.caseClient.Update(ctx, &casesapi.Case{
			ID:          id,
			Description: description,
			Done:        done == "on",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Location", "/cases")
	w.WriteHeader(http.StatusSeeOther)

}
