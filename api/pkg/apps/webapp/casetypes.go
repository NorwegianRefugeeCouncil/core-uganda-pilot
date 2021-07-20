package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/sessionmanager"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func (s *Server) CaseTypes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	cmsClient := s.CMSClient(ctx)
	iamClient := s.IAMClient(ctx)

	caseTypes, err := cmsClient.CaseTypes().List(ctx, cms.CaseTypeListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	partyTypes, err := iamClient.PartyTypes().List(ctx, iam.PartyTypeListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	teams, err := iamClient.Teams().List(ctx, iam.TeamListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		s.PostCaseType(ctx, &cms.CaseType{}, w, req)
		return
	}

	if err := s.renderFactory.New(req).ExecuteTemplate(w, "casetypes", map[string]interface{}{
		"CaseTypes":  caseTypes,
		"PartyTypes": partyTypes,
		"Teams":      teams,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) CaseType(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	cmsClient := s.CMSClient(ctx)
	iamClient := s.IAMClient(ctx)

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		http.Error(w, fmt.Errorf("no id in path").Error(), http.StatusInternalServerError)
		return
	}

	var caseType *cms.CaseType
	var partyTypes *iam.PartyTypeList
	var teamsData *iam.TeamList

	g, waitCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if id == "new" {
			caseType = &cms.CaseType{}
			return nil
		}
		var err error
		caseType, err = cmsClient.CaseTypes().Get(waitCtx, id)
		return err
	})
	g.Go(func() error {
		var err error
		partyTypes, err = iamClient.PartyTypes().List(waitCtx, iam.PartyTypeListOptions{})
		return err
	})
	g.Go(func() error {
		var err error
		teamsData, err = iamClient.Teams().List(ctx, iam.TeamListOptions{})
		return err
	})
	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if req.Method == "POST" {
		s.PostCaseType(ctx, caseType, w, req)
		return
	}

	if err := s.renderFactory.New(req).ExecuteTemplate(w, "casetype", map[string]interface{}{
		"CaseType":   caseType,
		"PartyTypes": partyTypes,
		"Teams":      teamsData,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Server) PostCaseType(
	ctx context.Context,
	caseType *cms.CaseType,
	w http.ResponseWriter,
	req *http.Request,
) {
	cmsClient := h.CMSClient(ctx)

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	values := req.Form

	err := caseType.UnmarshalFormData(values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isNew := len(caseType.ID) == 0
	if isNew {
		h.CreateCaseType(ctx, caseType, w, cmsClient)
	} else {
		h.UpdateCaseType(ctx, caseType, w, cmsClient)
	}
}

func (h *Server) CreateCaseType(ctx context.Context, caseType *cms.CaseType, w http.ResponseWriter, cmsClient cms.Interface) {
	_, err := cmsClient.CaseTypes().Create(ctx, caseType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.sessionManager.AddNotification(ctx, &sessionmanager.Notification{
		Message: fmt.Sprintf("Case type \"%s\" successfully created", caseType.Name),
		Theme:   "success",
	})
	w.Header().Set("Location", "/settings/casetypes")
	w.WriteHeader(http.StatusSeeOther)
}

func (h *Server) UpdateCaseType(ctx context.Context, caseType *cms.CaseType, w http.ResponseWriter, cmsClient cms.Interface) {
	_, err := cmsClient.CaseTypes().Update(ctx, caseType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.sessionManager.AddNotification(ctx, &sessionmanager.Notification{
		Message: fmt.Sprintf("Case type \"%s\" successfully updated", caseType.Name),
		Theme:   "success",
	})
	w.Header().Set("Location", "/settings/casetypes")
	w.WriteHeader(http.StatusSeeOther)
}

func (s *Server) NewCaseType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	iamClient := s.IAMClient(ctx)

	p, err := iamClient.PartyTypes().List(ctx, iam.PartyTypeListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	teamsData, err := iamClient.Teams().List(ctx, iam.TeamListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.renderFactory.New(req).ExecuteTemplate(w, "casetype", map[string]interface{}{
		"PartyTypes": p,
		"Teams":      teamsData,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
