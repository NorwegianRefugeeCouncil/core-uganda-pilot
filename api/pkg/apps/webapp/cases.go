package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/form"
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
		s.Case(w, req)
		return
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "cases", map[string]interface{}{
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
	var (
		err                error
		recipientParty     *iam.Party
		team               *iam.Team
		kase               *cms.Case
		parent             *cms.Case
		kaseTypes          *cms.CaseTypeList
		referrals          *cms.CaseList
		referralCaseType   *cms.CaseType
		creator            *iam.Party
		comments           *cms.CommentList
		commentAuthors     *iam.PartyList
		commentAuthorIDMap = make(map[string]bool)
		commentAuthorMap   = make(map[string]*iam.Party)
		commentAuthorIDs   []string
	)

	ctx := req.Context()
	qry := req.URL.Query()
	referralCaseTypeID := qry.Get("referralCaseTypeId")

	caseID, ok := mux.Vars(req)["id"]
	isNewCase := !ok || len(caseID) == 0
	if isNewCase {
		if req.Method != "POST" {
			caseID = "new"
		} else {
			err := fmt.Errorf("no id in path")
			s.Error(w, err)
			return
		}
	}

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

	if isNewCase {
		kase = &cms.Case{}
		if len(referralCaseTypeID) > 0 {
			kase.CaseTypeID = referralCaseTypeID
		}
	} else {
		if kase, err = cmsClient.Cases().Get(ctx, caseID); err != nil {
			s.Error(w, err)
			return
		}
	}

	if req.Method == "POST" {
		posted, err := s.PostCase(req, w, ctx, kase)
		if err != nil {
			if status, ok := err.(*validation.Status); ok {
				validatedElements := zipTemplateAndErrors(status.Errors, kase.Template)
				s.json(w, status.Code, validatedElements)
			} else {
				s.Error(w, err)
			}
			return
		}
		var action string
		if isNewCase {
			action = "created"
		} else {
			action = "updated"
		}
		if err := s.sessionManager.AddNotification(req, w, &sessionmanager.Notification{
			Message: fmt.Sprintf("Case successfully %s", action),
			Theme:   "success",
		}); err != nil {
			s.Error(w, err)
			return
		}

		if len(posted.ParentID) > 0 {
			w.Header().Set("Location", "/cases/"+posted.ParentID)
		} else {
			w.Header().Set("Location", "/cases/"+posted.ID)
		}
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	if len(kase.ID) > 0 {
		g, waitCtx := errgroup.WithContext(ctx)
		// Comments
		g.Go(func() error {
			if comments, err = cmsClient.Comments().List(waitCtx, cms.CommentListOptions{
				CaseID: caseID,
			}); err != nil {
				return err
			}
			// Get all comment author IDs
			for _, comment := range comments.Items {
				commentAuthorIDMap[comment.AuthorID] = true
			}
			for authorID := range commentAuthorIDMap {
				commentAuthorIDs = append(commentAuthorIDs, authorID)
			}
			if commentAuthors, err = iamClient.Parties().Search(waitCtx, iam.PartySearchOptions{
				PartyTypeIDs: []string{iam.IndividualPartyType.ID},
				PartyIDs:     commentAuthorIDs,
			}); err != nil {
				return err
			}
			for _, author := range commentAuthors.Items {
				commentAuthorMap[author.ID] = author
			}
			return nil
		})
		// Creator
		if len(kase.CreatorID) > 0 {
			g.Go(func() error {
				if creator, err = iamClient.Parties().Get(waitCtx, kase.CreatorID); err != nil {
					return err
				}
				return nil
			})
		}
		// Team
		g.Go(func() error {
			if team, err = iamClient.Teams().Get(waitCtx, kase.TeamID); err != nil {
				return err
			}
			return nil
		})
		// Parent Case
		if len(kase.ParentID) > 0 {
			g.Go(func() error {
				if parent, err = cmsClient.Cases().Get(waitCtx, kase.ParentID); err != nil {
					return err
				}
				return nil
			})
		}
		// Recipient & CaseTypes
		g.Go(func() error {
			if recipientParty, err = iamClient.Parties().Get(waitCtx, kase.PartyID); err != nil {
				return err
			}
			if kaseTypes, err = cmsClient.CaseTypes().List(waitCtx, cms.CaseTypeListOptions{
				PartyTypeIDs: recipientParty.PartyTypeIDs,
			}); err != nil {
				return err
			}
			return nil
		})
		// Referrals
		g.Go(func() error {
			if referrals, err = cmsClient.Cases().List(waitCtx, cms.CaseListOptions{ParentID: kase.ID}); err != nil {
				return err
			}
			return nil
		})
		if len(referralCaseTypeID) > 0 {
			g.Go(func() error {
				if referralCaseType, err = cmsClient.CaseTypes().Get(waitCtx, referralCaseTypeID); err != nil {
					return err
				}
				return nil
			})
		}
		if err = g.Wait(); err != nil {
			s.Error(w, err)
			return
		}
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "case", map[string]interface{}{
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
	}
}

// zipTemplateAndErrors returns a slice of form.FormElement populated with validated template form elements.
func zipTemplateAndErrors(errors validation.ErrorList, template *cms.CaseTemplate) []form.FormElement {
	var formElements []form.FormElement
	for _, element := range template.FormElements {
		if errs := errors.FindFamily(element.Attributes.Name); len(*errs) > 0 {
			element.Errors = errs
			formElements = append(formElements, element)
		}
	}
	return formElements
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

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "casenew", map[string]interface{}{
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

func (s *Server) PostCase(req *http.Request, w http.ResponseWriter, ctx context.Context, kase *cms.Case) (*cms.Case, error) {
	var err error

	cmsClient, err := s.CMSClient(req)
	if err != nil {
		return nil, err
	}

	if err := req.ParseForm(); err != nil {
		return nil, err
	}

	if kase.CaseTypeID == "" {
		// we need to get the caseTypeId from the form data
		kase.CaseTypeID = req.Form.Get("caseTypeId")
		if kase.CaseTypeID == "" {
			err = fmt.Errorf("unable to detect case type id for new case")
			return nil, err
		}
	}

	caseType, err := cmsClient.CaseTypes().Get(ctx, kase.CaseTypeID)
	if err != nil {
		return nil, err
	}

	err = UnmarshalCaseFormData(kase, caseType.Template, req.Form)
	if err != nil {
		return nil, err
	}

	var isNewCase = kase.ID == ""
	if isNewCase {
		subject := ctx.Value("Subject")
		if subject == nil {
			kase.CreatorID = ""
		} else {
			kase.CreatorID = subject.(string)
		}
		kase.IntakeCase = caseType.IntakeCaseType
		kase, err = cmsClient.Cases().Create(ctx, kase)
	} else {
		kase, err = cmsClient.Cases().Update(ctx, kase)
	}
	return kase, err
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

func (s *Server) processCaseValidation(req *http.Request, w http.ResponseWriter, kase *cms.Case, err error) {
	if err != nil {
		if status, ok := err.(*validation.Status); ok {
			s.renderWithValidation(req, w, kase, status)
			return
		} else {
			s.Error(w, err)
			return
		}
	} else {
		s.redirectAfterSuccessfulCasePost(req, w, kase)
	}
}

func (s *Server) redirectAfterSuccessfulCasePost(req *http.Request, w http.ResponseWriter, kase *cms.Case) {

}

func (s *Server) renderWithValidation(req *http.Request, w http.ResponseWriter, kase *cms.Case, status *validation.Status) {
	validatedForm := zipTemplateAndErrors(status.Errors, kase.Template)
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
	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "casenew", map[string]interface{}{
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
	var formElements []form.FormElement
	for _, formElement := range template.FormElements {
		formElement.Attributes.Value = values[formElement.Attributes.Name]
		formElements = append(formElements, formElement)
	}
	c.Template = &cms.CaseTemplate{FormElements: formElements}
	return nil
}
