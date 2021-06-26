package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/apps/iam"
	"github.com/nrc-no/core-kafka/pkg/cases/casetypes"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func (h *Handler) CaseTypes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	if req.Method == "POST" {
		h.PostCaseType(ctx, &casetypes.CaseType{}, w, req)
		return
	}

	caseTypes, err := h.caseTypeClient.List(ctx, casetypes.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	partyTypes, err := h.iam.PartyTypes().List(ctx, iam.PartyTypeListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	teams, err := h.iam.Teams().List(ctx, iam.TeamListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "casetypes", map[string]interface{}{
		"CaseTypes":  caseTypes,
		"PartyTypes": partyTypes,
		"Teams":      teams,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PostCaseType(
	ctx context.Context,
	caseType *casetypes.CaseType,
	w http.ResponseWriter,
	req *http.Request,
) {

	isNew := false
	if len(caseType.ID) == 0 {
		isNew = true
	}

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	values := req.Form
	caseType.Name = values.Get("name")
	caseType.PartyTypeID = values.Get("partyTypeId")
	caseType.TeamID = values.Get("teamId")

	if isNew {
		_, err := h.caseTypeClient.Create(ctx, caseType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/settings/casetypes")
		w.WriteHeader(http.StatusSeeOther)
		return
	} else {
		_, err := h.caseTypeClient.Update(ctx, caseType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/settings/casetypes")
		w.WriteHeader(http.StatusSeeOther)
		return
	}
}

func (h *Handler) CaseType(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var caseType *casetypes.CaseType
	var partyTypes *iam.PartyTypeList
	var teamsData *iam.TeamList

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if id == "new" {
			caseType = &casetypes.CaseType{}
			return nil
		}
		var err error
		caseType, err = h.caseTypeClient.Get(waitCtx, id)
		return err
	})

	g.Go(func() error {
		var err error
		partyTypes, err = h.iam.PartyTypes().List(waitCtx, iam.PartyTypeListOptions{})
		return err
	})

	g.Go(func() error {
		var err error
		teamsData, err = h.iam.Teams().List(ctx, iam.TeamListOptions{})
		return err
	})

	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostCaseType(ctx, caseType, w, req)
		return
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "casetype", map[string]interface{}{
		"CaseType":   caseType,
		"PartyTypes": partyTypes,
		"Teams":      teamsData,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) NewCaseType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	p, err := h.iam.PartyTypes().List(ctx, iam.PartyTypeListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	teamsData, err := h.iam.Teams().List(ctx, iam.TeamListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "casetype", map[string]interface{}{
		"PartyTypes": p,
		"Teams":      teamsData,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
