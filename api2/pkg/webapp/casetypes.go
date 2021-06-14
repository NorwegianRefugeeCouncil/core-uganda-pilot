package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	casesapi "github.com/nrc-no/core-kafka/pkg/cases/api"
	"github.com/nrc-no/core-kafka/pkg/cases/casetypes"
	partiesapi "github.com/nrc-no/core-kafka/pkg/parties/api"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func (h *Handler) CaseType(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	caseTypesClient := casetypes.NewClient("http://localhost:9000")
	partiesClient := parties.NewClient("http://localhost:9000")

	var caseTypes *casesapi.CaseTypeList
	var p *partiesapi.PartyList

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		caseTypes, err = caseTypesClient.List(waitCtx)
		return err
	})

	if req.Method == "POST" {
		h.PostCaseType(ctx, caseTypesClient, &casesapi.CaseType{}, w, req)
		return
	}

	g.Go(func() error {
		var err error
		p, err = partiesClient.List(waitCtx, parties.ListOptions{})
		return err
	})

	if err := h.template.ExecuteTemplate(w, "newcasetype", map[string]interface{}{
		"CaseTypes": caseTypes,
		"Parties":   p,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) NewCaseType(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	caseTypesClient := casetypes.NewClient("http://localhost:9000")
	partiesClient := parties.NewClient("http://localhost:9000")

	var caseTypes *casesapi.CaseTypeList
	var p *partiesapi.PartyList

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		caseTypes, err = caseTypesClient.List(waitCtx)
		return err
	})

	g.Go(func() error {
		var err error
		p, err = partiesClient.List(waitCtx, parties.ListOptions{})
		return err
	})

	if err := h.template.ExecuteTemplate(w, "newcasetype", map[string]interface{}{
		"CaseTypes": caseTypes,
		"Parties":   p,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PostCaseType(
	ctx context.Context,
	caseTypesCli *casetypes.Client,
	caseType *casesapi.CaseType,
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

	name := req.Form.Get("name")
	partypeID := req.Form.Get("partyTypeId")
	caseType.Name = name
	caseType.PartyTypeID = partypeID

	if isNew {
		out, err := caseTypesCli.Create(ctx, caseType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/cases/settings/casetypes/"+out.ID)
		w.WriteHeader(http.StatusSeeOther)
		return
	} else {
		_, err := caseTypesCli.Update(ctx, caseType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/cases")
		w.WriteHeader(http.StatusSeeOther)
		return
	}
}
