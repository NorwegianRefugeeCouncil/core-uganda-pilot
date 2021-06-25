package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/cases/cases"
	"github.com/nrc-no/core-kafka/pkg/cases/casetypes"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"github.com/nrc-no/core-kafka/pkg/teams"
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

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "cases", map[string]interface{}{
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

	var referrals *cases.CaseList

	var team *teams.Team

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

	if id != "new" {
		teamRes, err := h.teamClient.Get(ctx, kase.TeamID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		team = teamRes
	}

	party, err := h.partyClient.Get(ctx, kase.PartyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	kaseTypes, err = h.caseTypeClient.List(ctx, casetypes.ListOptions{
		PartyTypeIDs: party.PartyTypeIDs,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	referrals, err = h.caseClient.List(ctx, cases.ListOptions{
		ParentID: kase.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	qry := req.URL.Query()

	var referralCaseType *casetypes.CaseType
	if referralCaseTypeID := qry.Get("referralCaseTypeId"); len(referralCaseTypeID) > 0 {
		referralCaseType, err = h.caseTypeClient.Get(ctx, referralCaseTypeID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "case", map[string]interface{}{
		"Case":             kase,
		"CaseTypes":        kaseTypes,
		"Party":            party,
		"Parties":          partyList,
		"ReferralCaseType": referralCaseType,
		"Referrals":        referrals,
		"Team":             team,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) NewCase(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var caseTypes *casetypes.CaseTypeList
	var p *parties.PartyList

	g, waitCtx := errgroup.WithContext(ctx)


	g.Go(func() error {
		var err error
		caseTypes, err = h.caseTypeClient.List(waitCtx, casetypes.ListOptions{})
		return err
	})

	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	qry := req.URL.Query()
	caseTypeID := qry.Get("caseTypeId")
	partyTypeID := ""
	teamID := ""
	for _, caseType := range caseTypes.Items {
		if caseType.ID == caseTypeID {
			partyTypeID = caseType.PartyTypeID
			teamID = caseType.TeamID
			break
		}
	}

	listOptions := parties.ListOptions{
		PartyTypeID: partyTypeID,
	}

	p, err := h.partyClient.List(ctx, listOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var team *teams.Team
	if len(teamID) > 0 {
		team, err = h.teamClient.Get(ctx, teamID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "casenew", map[string]interface{}{
		"PartyID":    qry.Get("partyId"),
		"CaseTypeID": qry.Get("caseTypeId"),
		"Team":       team,
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
	parentId := req.Form.Get("parentId")
	teamId := req.Form.Get("teamId")

	var kase *cases.Case
	if id == "" {
		var err error
		kase, err = h.caseClient.Create(ctx, &cases.Case{
			CaseTypeID:  caseTypeId,
			PartyID:     partyId,
			Description: description,
			Done:        false,
			ParentID:    parentId,
			TeamID:      teamId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		var err error
		kase, err = h.caseClient.Update(ctx, &cases.Case{
			ID:          id,
			CaseTypeID:  caseTypeId,
			PartyID:     partyId,
			Description: description,
			Done:        done == "on",
			ParentID:    parentId,
			TeamID:      teamId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if len(parentId) > 0 {
		w.Header().Set("Location", "/cases/"+parentId)
	} else {
		w.Header().Set("Location", "/cases/"+kase.ID)
	}
	w.WriteHeader(http.StatusSeeOther)

}
