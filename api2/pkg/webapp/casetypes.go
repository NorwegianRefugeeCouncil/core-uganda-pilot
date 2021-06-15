package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	casesapi "github.com/nrc-no/core-kafka/pkg/cases/api"
	"github.com/nrc-no/core-kafka/pkg/cases/casetypes"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func (h *Handler) CaseTypes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var caseTypes *casesapi.CaseTypeList

	caseTypes, err := h.caseTypeClient.List(ctx, casetypes.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostCaseType(ctx, &casesapi.CaseType{}, w, req)
		return
	}

	if err := h.template.ExecuteTemplate(w, "casetypes", map[string]interface{}{
		"CaseTypes": caseTypes,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	var caseType *casesapi.CaseType
	var partyTypes *partytypes.PartyTypeList

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if id == "new" {
			caseType = &casesapi.CaseType{}
			return nil
		}
		var err error
		caseType, err = h.caseTypeClient.Get(waitCtx, id)
		return err
	})

	g.Go(func() error {
		var err error
		partyTypes, err = h.partyTypeClient.List(waitCtx)
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

	if err := h.template.ExecuteTemplate(w, "casetype", map[string]interface{}{
		"CaseType":   caseType,
		"PartyTypes": partyTypes,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) NewCaseType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	partyTypesClient := partytypes.NewClient("http://localhost:9000")

	p, err := partyTypesClient.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.template.ExecuteTemplate(w, "casetype", map[string]interface{}{
		"PartyTypes": p,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PostCaseType(
	ctx context.Context,
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
	caseType.Name = name
	partyTypeID := req.Form.Get("partyTypeId")
	caseType.PartyTypeID = partyTypeID

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
