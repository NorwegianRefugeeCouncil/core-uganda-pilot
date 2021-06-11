package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	casesapi "github.com/nrc-no/core-kafka/pkg/cases/api"
	"github.com/nrc-no/core-kafka/pkg/cases/cases"
	"github.com/nrc-no/core-kafka/pkg/cases/casetypes"
	partiesapi "github.com/nrc-no/core-kafka/pkg/parties/api"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
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
	caseTypesClient := casetypes.NewClient("http://localhost:9000")
	partyTypesClient := partytypes.NewClient("http://localhost:9000")

	var c *casesapi.Case
	var cList *casesapi.CaseList
	var cTypes *casesapi.CaseTypeList
	var pTypes *partiesapi.PartyTypeList

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

	g.Go(func() error {
		var err error
		cTypes, err = caseTypesClient.List(waitCtx)
		return err
	})

	g.Go(func() error {
		var err error
		pTypes, err = partyTypesClient.List(waitCtx)
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
		"Case":       c,
		"Cases":      cList,
		"CaseTypes":  cTypes,
		"PartyTypes": pTypes,
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
	var p *partiesapi.PartyList

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		caseTypes, err = caseTypesClient.List(waitCtx)
		return err
	})

	var listoptions parties.ListOptions
	listoptions.Party = ""

	g.Go(func() error {
		var err error
		p, err = partiesClient.List(waitCtx, listoptions)
		return err
	})

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
		case "caseType":
			kase.CaseTypeID = value[0]
		case "partyId":
			kase.PartyID = value[0]
		case "description":
			kase.Description = value[0]
		}
	}

	if id == "" {
		_, err := cli.Create(ctx, &casesapi.Case{
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
