package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/sessionmanager"
	"github.com/nrc-no/core/pkg/validation"
	"github.com/satori/go.uuid"
	"github.com/xeonx/timeago"
	"golang.org/x/sync/errgroup"
	"net/http"
	"net/url"
)

func (s *Server) Cases(w http.ResponseWriter, req *http.Request) {

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

	options := &CasesListOptions{}
	if err := options.UnmarshalQueryParams(req.URL.Query()); err != nil {
		s.Error(w, err)
		return
	}

	var kases *cms.CaseList
	var caseTypes *cms.CaseTypeList
	var partyList *iam.PartyList
	var teams *iam.TeamList

	caseListOptions := &cms.CaseListOptions{}
	caseListOptions.Done = options.Closed
	caseListOptions.CaseTypeIDs = options.CaseTypeIDs
	caseListOptions.TeamIDs = options.TeamIDs

	wg, wgCtx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		var err error
		kases, err = cmsClient.Cases().List(wgCtx, *caseListOptions)
		return err
	})
	wg.Go(func() error {
		var err error
		caseTypes, err = cmsClient.CaseTypes().List(wgCtx, cms.CaseTypeListOptions{})
		return err
	})
	wg.Go(func() error {
		var err error
		partyList, err = iamClient.Parties().List(wgCtx, iam.PartyListOptions{})
		return err
	})
	wg.Go(func() error {
		var err error
		teams, err = iamClient.Teams().List(wgCtx, iam.TeamListOptions{})
		return err
	})

	if err := wg.Wait(); err != nil {
		s.Error(w, err)
		return
	}

	if req.Method == "POST" {
		s.PostCase(req, w, ctx, &cms.Case{})
		return
	}

	if err := s.renderFactory.New(req).ExecuteTemplate(w, "cases", map[string]interface{}{
		"Cases":         kases,
		"CaseTypes":     caseTypes,
		"Parties":       partyList,
		"Teams":         teams,
		"FilterOptions": options,
	}); err != nil {
		s.Error(w, err)
		return
	}

}

func (s *Server) Case(w http.ResponseWriter, req *http.Request) {

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

	caseID, ok := mux.Vars(req)["id"]
	if !ok || len(caseID) == 0 {
		err := fmt.Errorf("no id in path")
		s.Error(w, err)
		return
	}

	var recipientParty *iam.Party
	var team *iam.Team
	var kase *cms.Case
	var parent *cms.Case
	var kaseTypes *cms.CaseTypeList
	var referrals *cms.CaseList
	var comments *cms.CommentList
	var creator *iam.Party

	g, waitCtx := errgroup.WithContext(ctx)

	// Get Case
	g.Go(func() error {
		if caseID == "new" {
			kase = &cms.Case{}
			return nil
		}
		var err error
		kase, err = cmsClient.Cases().Get(waitCtx, caseID)
		return err
	})

	// Get comments for case
	g.Go(func() error {
		var err error
		comments, err = cmsClient.Comments().List(waitCtx, cms.CommentListOptions{
			CaseID: caseID,
		})
		return err
	})

	if err := g.Wait(); err != nil {
		s.Error(w, err)
		return
	}

	if caseID != "new" && len(kase.CreatorID) > 0 {
		var err error
		creator, err = iamClient.Parties().Get(ctx, kase.CreatorID)
		if err != nil {
			s.Error(w, err)
			return
		}
	}

	// Get all comment author IDs
	var commentAuthorIDMap = map[string]bool{}
	for _, comment := range comments.Items {
		commentAuthorIDMap[comment.AuthorID] = true
	}
	var commentAuthorIDs []string
	for authorID := range commentAuthorIDMap {
		commentAuthorIDs = append(commentAuthorIDs, authorID)
	}

	// Get all authors
	commentAuthors, err := iamClient.Parties().Search(ctx, iam.PartySearchOptions{
		PartyTypeIDs: []string{iam.IndividualPartyType.ID},
		PartyIDs:     commentAuthorIDs,
	})
	if err != nil {
		s.Error(w, err)
		return
	}
	var commentAuthorMap = map[string]*iam.Party{}
	for _, author := range commentAuthors.Items {
		commentAuthorMap[author.ID] = author
	}

	// Get case type team
	if caseID != "new" {
		teamRes, err := iamClient.Teams().Get(ctx, kase.TeamID)
		if err != nil {
			s.Error(w, err)
			return
		}
		team = teamRes
	}

	// Get parent case
	if len(kase.ParentID) > 0 {
		parent, err = cmsClient.Cases().Get(ctx, kase.ParentID)
		if err != nil {
			s.Error(w, err)
			return
		}
	}

	// Get case recipient
	recipientParty, err = iamClient.Parties().Get(ctx, kase.PartyID)
	if err != nil {
		s.Error(w, err)
		return
	}

	// Get case types
	kaseTypes, err = cmsClient.CaseTypes().List(ctx, cms.CaseTypeListOptions{
		PartyTypeIDs: recipientParty.PartyTypeIDs,
	})
	if err != nil {
		s.Error(w, err)
		return
	}

	// Get case referrals
	referrals, err = cmsClient.Cases().List(ctx, cms.CaseListOptions{
		ParentID: kase.ID,
	})
	if err != nil {
		s.Error(w, err)
		return
	}
	qry := req.URL.Query()

	var referralCaseType *cms.CaseType
	if referralCaseTypeID := qry.Get("referralCaseTypeId"); len(referralCaseTypeID) > 0 {
		referralCaseType, err = cmsClient.CaseTypes().Get(ctx, referralCaseTypeID)
		if err != nil {
			s.Error(w, err)
			return
		}
	}
	if req.Method == "POST" {
		s.PostCase(req, w, ctx, kase)
		return
	}
	if err := s.renderFactory.New(req).ExecuteTemplate(w, "case", map[string]interface{}{
		"Case":             kase,
		"Parent":           parent,
		"CaseTypes":        kaseTypes,
		"Recipient":        recipientParty,
		"ReferralCaseType": referralCaseType,
		"Referrals":        referrals,
		"Team":             team,
		"CreatedBy":        creator,
		"Comments":         displayComments(comments, commentAuthorMap),
	}); err != nil {
		s.Error(w, err)
		return
	}

}

func (s *Server) NewCase(w http.ResponseWriter, req *http.Request) {

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

	var caseTypes *cms.CaseTypeList
	var p *iam.PartyList

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		caseTypes, err = cmsClient.CaseTypes().List(waitCtx, cms.CaseTypeListOptions{})
		return err
	})

	if err := g.Wait(); err != nil {
		s.Error(w, err)
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

	p, err = iamClient.Parties().List(ctx, listOptions)
	if err != nil {
		s.Error(w, err)
	}

	var caseType *cms.CaseType
	if caseTypeID != "" {
		caseType, err = cmsClient.CaseTypes().Get(ctx, caseTypeID)
		if err != nil {
			s.Error(w, err)
			return
		}
	}

	var team *iam.Team
	if len(teamID) > 0 {
		team, err = iamClient.Teams().Get(ctx, teamID)
		if err != nil {
			s.Error(w, err)
			return
		}
	}

	if err := s.renderFactory.New(req).ExecuteTemplate(w, "casenew", map[string]interface{}{
		"PartyID":   qry.Get("partyId"),
		"CaseType":  caseType,
		"Team":      team,
		"CaseTypes": caseTypes,
		"Parties":   p,
	}); err != nil {
		s.Error(w, err)
		return
	}
}

func (s *Server) PostCase(req *http.Request, w http.ResponseWriter, ctx context.Context, kase *cms.Case) {
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

	caseTypeId := req.Form.Get("caseTypeId")

	caseType, err := cmsClient.CaseTypes().Get(ctx, caseTypeId)
	if err != nil {
		s.Error(w, err)
		return
	}

	err = UnmarshalCaseFormData(kase, caseType.Template, req.Form)
	if err != nil {
		s.Error(w, err)
		return
	}

	var isNewCase = kase.ID == ""
	var action string
	var postedCase *cms.Case
	if isNewCase {
		subject := ctx.Value("Subject")
		if subject == nil {
			kase.CreatorID = ""
		} else {
			kase.CreatorID = subject.(string)
		}
		kase.IntakeCase = caseType.IntakeCaseType
		postedCase, err = cmsClient.Cases().Create(ctx, kase)
		action = "created"
	} else {
		postedCase, err = cmsClient.Cases().Update(ctx, kase)
		action = "updated"
	}
	if err != nil {
		if status, ok := err.(*validation.Status); ok {
			s.renderCaseValidation(req, w, kase, status)
			return
		} else {
			s.Error(w, err)
			return
		}
	}
	s.redirectAfterSuccessfulCasePost(req, w, postedCase, action)
}

type displayComment struct {
	*cms.Comment
	Author  *iam.Party
	TimeAgo string
}

func displayComments(comments *cms.CommentList, authorMap map[string]*iam.Party) []*displayComment {
	var displayComments []*displayComment
	for _, item := range comments.Items {
		c := &displayComment{
			Comment: item,
			TimeAgo: timeago.English.Format(item.CreatedAt),
			Author:  authorMap[item.AuthorID],
		}
		displayComments = append(displayComments, c)
	}
	return displayComments
}

func (s *Server) processCaseValidation(req *http.Request, w http.ResponseWriter, kase *cms.Case, err error, action string) {
	if err != nil {
		if status, ok := err.(*validation.Status); ok {
			s.renderCaseValidation(req, w, kase, status)
			return
		} else {
			s.Error(w, err)
			return
		}
	} else {
		s.redirectAfterSuccessfulCasePost(req, w, kase, action)
	}
}

func (s *Server) redirectAfterSuccessfulCasePost(req *http.Request, w http.ResponseWriter, kase *cms.Case, action string) {
	if err := s.sessionManager.AddNotification(req, w, &sessionmanager.Notification{
		Message: fmt.Sprintf("Case successfully %s", action),
		Theme:   "success",
	}); err != nil {
		s.Error(w, err)
		return
	}
	if len(kase.ParentID) > 0 {
		w.Header().Set("Location", "/cases/"+kase.ParentID)
	} else {
		w.Header().Set("Location", "/cases/"+kase.ID)
	}
	w.WriteHeader(http.StatusSeeOther)
}

func (s *Server) renderCaseValidation(req *http.Request, w http.ResponseWriter, kase *cms.Case, status *validation.Status) {
	validatedForm := NewValidatedTemplate(kase.Template, status.Errors)
	parties, err := s.retrieveParties(req)
	if err != nil {
		s.Error(w, err)
		return
	}
	caseTypes, err := s.retrieveCaseTypes(req)
	if err != nil {
		s.Error(w, err)
		return
	}
	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}
	cmsClient, err := s.CMSClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}
	team, err := iamClient.Teams().Get(req.Context(), kase.TeamID)
	if err != nil {
		s.Error(w, err)
	}
	caseType, err := cmsClient.CaseTypes().Get(req.Context(), kase.CaseTypeID)
	if err != nil {
		s.Error(w, err)
	}
	qry := req.URL.Query()

	partyID := qry.Get("partyId")
	if len(partyID) == 0 {
		partyID = kase.PartyID
	}
	// Set notification and render
	s.validationErrorNotification(req, w)
	if err := s.renderFactory.New(req).ExecuteTemplate(w, "casenew", map[string]interface{}{
		"Team":          team,
		"CaseTypes":     caseTypes,
		"CaseType":      caseType,
		"Parties":       parties,
		"PartyID":       partyID,
		"ValidatedForm": validatedForm,
	}); err != nil {
		s.Error(w, err)
		return
	}

}

type CasesListOptions struct {
	Closed      *bool
	TeamIDs     []string
	CaseTypeIDs []string
}

func (c *CasesListOptions) ClosedOnly() bool {
	return c.Closed != nil && *c.Closed == true
}

func (c *CasesListOptions) OpenOnly() bool {
	return c.Closed != nil && *c.Closed == false
}

func (c *CasesListOptions) UnmarshalQueryParams(values url.Values) error {

	if len(values["status"]) == 1 {
		closed := values["status"][0] == "closed"
		c.Closed = &closed
	}

	for _, teamId := range values["teamId"] {
		if _, err := uuid.FromString(teamId); err == nil {
			c.TeamIDs = append(c.TeamIDs, teamId)
		}
	}

	for _, caseTypeId := range values["caseTypeId"] {
		if _, err := uuid.FromString(caseTypeId); err == nil {
			c.CaseTypeIDs = append(c.CaseTypeIDs, caseTypeId)
		}
	}

	return nil

}

// UnmarshalCaseFormData retrieves entries from a url.Values and applies them to a cms.Case object via a cms.CaseTemplate.
func UnmarshalCaseFormData(c *cms.Case, template *cms.CaseTemplate, values url.Values) error {
	c.CaseTypeID = values.Get("caseTypeId")
	c.PartyID = values.Get("partyId")
	c.Done = values.Get("done") == "on"
	c.ParentID = values.Get("parentId")
	c.TeamID = values.Get("teamId")
	var formElements []cms.FormElement
	for _, formElement := range template.FormElements {
		formElement.Attributes.Value = values[formElement.Attributes.Name]
		formElements = append(formElements, formElement)
	}
	c.Template = &cms.CaseTemplate{FormElements: formElements}
	return nil
}
