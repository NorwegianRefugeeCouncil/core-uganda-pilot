package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/apps/cms"
	"github.com/nrc-no/core-kafka/pkg/apps/iam"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func (h *Server) Cases(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	kases, err := h.cms.Cases().List(ctx, cms.CaseListOptions{})
	caseTypes, err := h.cms.CaseTypes().List(ctx, cms.CaseTypeListOptions{})
	partyList, err := h.iam.Parties().List(ctx, iam.PartyListOptions{})

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

func (h *Server) Case(w http.ResponseWriter, req *http.Request) {

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

	var party *iam.Party
	var partyList *iam.PartyList
	var team *iam.Team

	var kase *cms.Case
	var kaseTypes *cms.CaseTypeList
	var referrals *cms.CaseList

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if id == "new" {
			kase = &cms.Case{}
			return nil
		}
		var err error
		kase, err = h.cms.Cases().Get(waitCtx, id)
		return err
	})

	g.Go(func() error {
		var err error
		partyList, err = h.iam.Parties().List(waitCtx, iam.PartyListOptions{})
		return err
	})

	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if id != "new" {
		teamRes, err := h.iam.Teams().Get(ctx, kase.TeamID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		team = teamRes
	}

	party, err := h.iam.Parties().Get(ctx, kase.PartyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	kaseTypes, err = h.cms.CaseTypes().List(ctx, cms.CaseTypeListOptions{
		PartyTypeIDs: party.PartyTypeIDs,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	referrals, err = h.cms.Cases().List(ctx, cms.CaseListOptions{
		ParentID: kase.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	qry := req.URL.Query()

	var referralCaseType *cms.CaseType
	if referralCaseTypeID := qry.Get("referralCaseTypeId"); len(referralCaseTypeID) > 0 {
		referralCaseType, err = h.cms.CaseTypes().Get(ctx, referralCaseTypeID)
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

func (h *Server) NewCase(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var caseTypes *cms.CaseTypeList
	var p *iam.PartyList

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		caseTypes, err = h.cms.CaseTypes().List(waitCtx, cms.CaseTypeListOptions{})
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

	listOptions := iam.PartyListOptions{
		PartyTypeID: partyTypeID,
	}

	p, err := h.iam.Parties().List(ctx, listOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var team *iam.Team
	if len(teamID) > 0 {
		team, err = h.iam.Teams().Get(ctx, teamID)
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

func (h *Server) PostCase(ctx context.Context, id string, w http.ResponseWriter, req *http.Request) {

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

	var kase *cms.Case
	if id == "" {
		var err error
		kase, err = h.cms.Cases().Create(ctx, &cms.Case{
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
		kase, err = h.cms.Cases().Update(ctx, &cms.Case{
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
