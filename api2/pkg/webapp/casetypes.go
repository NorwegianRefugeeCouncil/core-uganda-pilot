package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/cases/api"
	"github.com/nrc-no/core-kafka/pkg/cases/casetypes"
	"net/http"
)

func (h *Handler) CaseTypes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	caseTypesCli := casetypes.NewClient("http://localhost:9000")

	if req.Method == "POST" {
		h.PostCaseType(ctx, caseTypesCli, &api.CaseType{}, w, req)
		return
	}

	caseTypes, err := caseTypesCli.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		err := fmt.Errorf("no id found in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	caseTypesCli := casetypes.NewClient("http://localhost:9000")

	var caseType = &api.CaseType{}
	if id != "new" {
		var err error
		caseType, err = caseTypesCli.Get(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if req.Method == "POST" {
		h.PostCaseType(ctx, caseTypesCli, caseType, w, req)
		return
	}

	if err := h.template.ExecuteTemplate(w, "casetypes", map[string]interface{}{
		"CaseType": caseType,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) NewCaseType(w http.ResponseWriter, req *http.Request) {
	if err := h.template.ExecuteTemplate(w, "casetype", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PostCaseType(
	ctx context.Context,
	caseTypesCli *casetypes.Client,
	caseType *api.CaseType,
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
		w.Header().Set("Location", "/settings/casetypes/"+out.ID)
		w.WriteHeader(http.StatusSeeOther)
		return
	} else {
		_, err := caseTypesCli.Update(ctx, caseType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/settings/casetypes/"+caseType.ID)
		w.WriteHeader(http.StatusSeeOther)
		return
	}
}
