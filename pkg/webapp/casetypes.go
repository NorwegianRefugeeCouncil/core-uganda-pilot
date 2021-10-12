package webapp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/cms"
	"github.com/nrc-no/core/pkg/form"
	"github.com/nrc-no/core/pkg/iam"
	"github.com/nrc-no/core/pkg/sessionmanager"
	"github.com/nrc-no/core/pkg/validation"
	"golang.org/x/sync/errgroup"
	"net/http"
	"net/url"
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

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "casetypes", map[string]interface{}{
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

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "casetype", map[string]interface{}{
		"CaseType":   caseType,
		"PartyTypes": partyTypes,
		"Teams":      teams,
	}); err != nil {
		s.Error(w, err)
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

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "casetype", map[string]interface{}{
		"PartyTypes": p,
		"Teams":      teamsData,
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
	var err error
	cmsClient, err := s.CMSClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	if err := req.ParseForm(); err != nil {
		s.Error(w, err)
		return
	}

	if err := UnmarshalCaseTypeFormData(caseType, req.Form); err != nil {
		s.Error(w, err)
		return
	}

	var action string
	isNewCaseType := len(caseType.ID) == 0
	if isNewCaseType {
		_, err = cmsClient.CaseTypes().Create(ctx, caseType)
		action = "created"
	} else {
		_, err = cmsClient.CaseTypes().Update(ctx, caseType)
		action = "updated"
	}
	s.processCaseTypeValidation(req, w, caseType, partyTypes, teams, err, action)
}

func (s *Server) processCaseTypeValidation(req *http.Request, w http.ResponseWriter, caseType *cms.CaseType, partyTypes *iam.PartyTypeList, teams *iam.TeamList, err error, action string) {
	// Examine the error argument
	if err != nil {
		if status, ok := err.(*validation.Status); ok {
			// If the error is a validation.Status proceed with rendering the validated form.
			s.renderCaseTypeValidation(req, w, caseType, partyTypes, teams, status)
		} else {
			// Otherwise write the error
			s.Error(w, err)
		}
	} else {
		// If there was no error was passed in, proceed with success
		s.redirectAfterSuccessfulCaseTypePost(req, w, caseType, action)
	}
}

func (s *Server) renderCaseTypeValidation(req *http.Request, w http.ResponseWriter, caseType *cms.CaseType, partyTypes *iam.PartyTypeList, teams *iam.TeamList, status *validation.Status) {
	// Set notification
	if err := s.sessionManager.AddNotification(req, w, &sessionmanager.Notification{
		Message: "There seems to be an problem with the data you have submitted. See below for errors.",
		Theme:   "danger",
	}); err != nil {
		s.Error(w, err)
		return
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "casetype", map[string]interface{}{
		"CaseType":   caseType,
		"PartyTypes": partyTypes,
		"Teams":      teams,
		"Errors":     status.Errors,
	}); err != nil {
		s.Error(w, err)
		return
	}
	return
}

func (s *Server) redirectAfterSuccessfulCaseTypePost(req *http.Request, w http.ResponseWriter, caseType *cms.CaseType, action string) bool {
	if err := s.sessionManager.AddNotification(req, w, &sessionmanager.Notification{
		Message: fmt.Sprintf("Case type \"%s\" successfully %s", caseType.Name, action),
		Theme:   "success",
	}); err != nil {
		s.Error(w, err)
		return true
	}
	w.Header().Set("Location", "/settings/casetypes")
	w.WriteHeader(http.StatusSeeOther)
	return false
}

// UnmarshalCaseTypeFormData retrieves entries from a url.Values and applies them to a cms.CaseType object.
func UnmarshalCaseTypeFormData(c *cms.CaseType, values url.Values) error {
	c.Name = values.Get("name")
	c.PartyTypeID = values.Get("partyTypeId")
	c.TeamID = values.Get("teamId")
	templateString := values.Get("template")
	if templateString == "" {
		c.Form = form.Form{}
	} else {
		if err := json.Unmarshal([]byte(templateString), &c.Form); err != nil {
			return err
		}
	}
	return nil
}
