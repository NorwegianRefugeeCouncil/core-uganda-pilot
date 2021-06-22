package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
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

	var kase *cases.Case
	var kaseTypes *casetypes.CaseTypeList

	var party *parties.Party
	var partyList *parties.PartyList

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if id == "new" {
			kase = &cases.Case{}
			return nil
		}
		var err error
		kase, err = h.caseClient.Get(waitCtx, id)
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

	party, err := h.partyClient.Get(ctx, kase.PartyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	kaseTypes, err = h.caseTypeClient.List(waitCtx, casetypes.ListOptions{
		PartyTypes: party.PartyTypes,
	})

	qry := req.URL.Query()
	refCaseTypeID := qry.Get("refCaseTypeId")
	var refCaseType *casetypes.CaseType
	if len(refCaseTypeID) > 0 {
		refCaseType, err = h.caseTypeClient.Get(ctx, refCaseTypeID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := h.template.ExecuteTemplate(w, "case", map[string]interface{}{
		"Case":        kase,
		"CaseTypes":   kaseTypes,
		"Party":       party,
		"Parties":     partyList,
		"RefCaseType": refCaseType,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) NewCase(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	caseTypesClient := casetypes.NewClient("http://localhost:9000")
	partiesClient := parties.NewClient("http://localhost:9000")

	var caseTypes *casetypes.CaseTypeList
	var p *parties.PartyList

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		caseTypes, err = caseTypesClient.List(waitCtx, casetypes.ListOptions{})
		return err
	})

	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	qry := req.URL.Query()
	caseTypeID := qry.Get("caseTypeId")
	partyTypeID := ""
	for _, caseType := range caseTypes.Items {
		if caseType.ID == caseTypeID {
			partyTypeID = caseType.PartyTypeID
		}
	}

	listOptions := parties.ListOptions{
		PartyTypeID: partyTypeID,
	}

	p, err := partiesClient.List(waitCtx, listOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

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
	done := req.Form.Get("done")

	if id == "" {
		_, err := h.caseClient.Create(ctx, &cases.Case{
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
		_, err := h.caseClient.Update(ctx, &cases.Case{
			ID:          id,
			CaseTypeID:  caseTypeId,
			PartyID:     partyId,
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
