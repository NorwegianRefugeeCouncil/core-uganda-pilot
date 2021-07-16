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
	if errResponse(w, err) {
		return
	}
	partyTypes, err := iamClient.PartyTypes().List(ctx, iam.PartyTypeListOptions{})
	if errResponse(w, err) {
		return
	}
	teams, err := iamClient.Teams().List(ctx, iam.TeamListOptions{})
	if errResponse(w, err) {
		return
	}

	if req.Method == "POST" {
		h.PostCaseType(ctx, &cms.CaseType{}, w, req, partyTypes, teams)
		return
	}

	if err := s.renderFactory.New(req).ExecuteTemplate(w, "casetypes", map[string]interface{}{
		"CaseTypes":  caseTypes,
		"PartyTypes": partyTypes,
		"Teams":      teams,
	}); errResponse(w, err) {
		return
	}
}

func (s *Server) CaseType(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	cmsClient := s.CMSClient(ctx)
	iamClient := s.IAMClient(ctx)

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		errResponse(w, fmt.Errorf("no id in path"))
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
	if err := g.Wait(); errResponse(w, err) {
		return
	}

	if req.Method == "POST" {
		h.PostCaseType(ctx, caseType, w, req, partyTypes, teamsData)
		return
	}

	if err := s.renderFactory.New(req).ExecuteTemplate(w, "casetype", map[string]interface{}{
		"CaseType":   caseType,
		"PartyTypes": partyTypes,
		"Teams":      teamsData,
		//"ErrList": validation,
	}); errResponse(w, err) {
		return
	}
}

func (h *Server) PostCaseType(
	ctx context.Context,
	caseType *cms.CaseType,
	w http.ResponseWriter,
	req *http.Request,
	partyTypes *iam.PartyTypeList,
	teamsData *iam.TeamList,
) {
	cmsClient := h.CMSClient(ctx)

	isNew := false
	if len(caseType.ID) == 0 {
		isNew = true
	}

	if err := req.ParseForm(); errResponse(w, err) {
		return
	}

	values := req.Form

	err := caseType.UnmarshalFormData(values)
	if errResponse(w, err) {
		return
	}

	if isNew {
		h.CreateCaseType(ctx, caseType, w, cmsClient)
	} else {
		h.UpdateCaseType(ctx, caseType, w, cmsClient)
	}
}

func (h *Server) CreateCaseType(ctx context.Context, caseType *cms.CaseType, w http.ResponseWriter, cmsClient cms.Interface) {
	_, err := cmsClient.CaseTypes().Create(ctx, caseType)
	if errResponse(w, err) {
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
	if errResponse(w, err) {
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
	if errResponse(w, err) {
		return
	}

	teamsData, err := iamClient.Teams().List(ctx, iam.TeamListOptions{})
	if errResponse(w, err) {
		return
	}

	if err := s.renderFactory.New(req).ExecuteTemplate(w, "casetype", map[string]interface{}{
		"PartyTypes": p,
		"Teams":      teamsData,
	}); errResponse(w, err) {
		return
	}
}

func errResponse(w http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}
