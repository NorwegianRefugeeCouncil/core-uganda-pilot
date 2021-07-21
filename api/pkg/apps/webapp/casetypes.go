package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/sessionmanager"
	"github.com/nrc-no/core/pkg/validation"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func (s *Server) CaseTypes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	cmsClient, err := s.CMSClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	caseTypes, err := cmsClient.CaseTypes().List(ctx, cms.CaseTypeListOptions{})
	if err != nil {
		s.Error(w, err)
		return
	}

	partyTypes, err := iamClient.PartyTypes().List(ctx, iam.PartyTypeListOptions{})
	if err != nil {
		s.Error(w, err)
		return
	}

	teams, err := iamClient.Teams().List(ctx, iam.TeamListOptions{})
	if err != nil {
		s.Error(w, err)
		return
	}

	if req.Method == "POST" {
		s.PostCaseType(ctx, &cms.CaseType{}, partyTypes, teams, w, req)
		return
	}

	if err := s.renderFactory.New(req).ExecuteTemplate(w, "casetypes", map[string]interface{}{
		"CaseTypes":  caseTypes,
		"PartyTypes": partyTypes,
		"Teams":      teams,
	}); err != nil {
		s.Error(w, err)
		return
	}
}

func (s *Server) CaseType(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	cmsClient, err := s.CMSClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		s.Error(w, err)
		return
	}

	var caseType *cms.CaseType
	var partyTypes *iam.PartyTypeList
	var teams *iam.TeamList

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
		teams, err = iamClient.Teams().List(ctx, iam.TeamListOptions{})
		return err
	})

	if err := g.Wait(); err != nil {
		s.Error(w, err)
		return
	}

	if req.Method == "POST" {
		s.PostCaseType(ctx, caseType, partyTypes, teams, w, req)
		return
	}

	if err := s.renderFactory.New(req).ExecuteTemplate(w, "casetype", map[string]interface{}{
		"CaseType":   caseType,
		"PartyTypes": partyTypes,
		"Teams":      teams,
	}); err != nil {
		s.Error(w, err)
		return
	}
}

func (s *Server) PostCaseType(
	ctx context.Context,
	caseType *cms.CaseType,
	partyTypes *iam.PartyTypeList,
	teams *iam.TeamList,
	w http.ResponseWriter,
	req *http.Request,
) {

	cmsClient, err := s.CMSClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	isNewCaseType := len(caseType.ID) == 0

	if err := req.ParseForm(); err != nil {
		s.Error(w, err)
		return
	}

	if err := caseType.UnmarshalFormData(req.Form); err != nil {
		s.Error(w, err)
		return
	}

	if isNewCaseType {
		_, err := cmsClient.CaseTypes().Create(ctx, caseType)
		if err != nil {
			if status, ok := err.(*validation.Status); ok {
				if err := s.renderFactory.New(req).ExecuteTemplate(w, "casetype", map[string]interface{}{
					"CaseType":   caseType,
					"PartyTypes": partyTypes,
					"Teams":      teams,
					"Errors":     status.Errors,
				}); err != nil {
					s.Error(w, err)
					return
				}
			} else {
				s.Error(w, err)
				return
			}
		}

		if err := s.sessionManager.AddNotification(req, w, &sessionmanager.Notification{
			Message: fmt.Sprintf("Case type \"%s\" successfully created", caseType.Name),
			Theme:   "success",
		}); err != nil {
			s.Error(w, err)
			return
		}

		w.Header().Set("Location", "/settings/casetypes")
		w.WriteHeader(http.StatusSeeOther)
		return
	} else {
		_, err := cmsClient.CaseTypes().Update(ctx, caseType)
		if err != nil {
			s.Error(w, err)
			return
		}

		if err := s.sessionManager.AddNotification(req, w, &sessionmanager.Notification{
			Message: fmt.Sprintf("Case type \"%s\" successfully updated", caseType.Name),
			Theme:   "success",
		}); err != nil {
			s.Error(w, err)
			return
		}

		w.Header().Set("Location", "/settings/casetypes")
		w.WriteHeader(http.StatusSeeOther)
		return
	}
}

func (s *Server) NewCaseType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	p, err := iamClient.PartyTypes().List(ctx, iam.PartyTypeListOptions{})
	if err != nil {
		s.Error(w, err)
		return
	}

	teamsData, err := iamClient.Teams().List(ctx, iam.TeamListOptions{})
	if err != nil {
		s.Error(w, err)
		return
	}

	if err := s.renderFactory.New(req).ExecuteTemplate(w, "casetype", map[string]interface{}{
		"PartyTypes": p,
		"Teams":      teamsData,
	}); err != nil {
		s.Error(w, err)
		return
	}
}
