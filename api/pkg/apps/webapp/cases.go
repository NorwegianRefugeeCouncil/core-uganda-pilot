package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/satori/go.uuid"
	"github.com/xeonx/timeago"
	"golang.org/x/sync/errgroup"
	"net/http"
	"net/url"
)

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

	hasStatusArg := false
	for key, _ := range values {
		if key == "status" {
			hasStatusArg = true
		}
	}

	statusStr := values.Get("status")
	if !hasStatusArg && len(statusStr) == 0 {
		statusStr = "false"
	}
	if len(statusStr) > 0 {
		isClosed := statusStr == "closed"
		c.Closed = &isClosed
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

func (h *Server) Cases(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	cmsClient := h.CMSClient(ctx)
	iamClient := h.IAMClient(ctx)

	options := &CasesListOptions{}
	if err := options.UnmarshalQueryParams(req.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		var err error
		kases, err = cmsClient.Cases().List(ctx, *caseListOptions)
		return err
	})
	wg.Go(func() error {
		var err error
		caseTypes, err = cmsClient.CaseTypes().List(ctx, cms.CaseTypeListOptions{})
		return err
	})
	wg.Go(func() error {
		var err error
		partyList, err = iamClient.Parties().List(ctx, iam.PartyListOptions{})
		return err
	})
	wg.Go(func() error {
		var err error
		teams, err = iamClient.Teams().List(ctx, iam.TeamListOptions{})
		return err
	})

	if err := wg.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostCase(ctx, "", w, req)
		return
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "cases", map[string]interface{}{
		"Cases":         kases,
		"CaseTypes":     caseTypes,
		"Parties":       partyList,
		"Teams":         teams,
		"FilterOptions": options,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Server) Case(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	cmsClient := h.CMSClient(ctx)
	iamClient := h.IAMClient(ctx)

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostCase(ctx, id, w, req)
		return
	}

	var party *iam.Party
	var partyList *iam.PartyList
	var team *iam.Team

	var kase *cms.Case
	var kaseTypes *cms.CaseTypeList
	var referrals *cms.CaseList
	var comments *cms.CommentList

	g, waitCtx := errgroup.WithContext(ctx)

	// Get Case
	g.Go(func() error {
		if id == "new" {
			kase = &cms.Case{}
			return nil
		}
		var err error
		kase, err = cmsClient.Cases().Get(waitCtx, id)
		return err
	})

	// TODO : Remove this? We cannot list all parties
	g.Go(func() error {
		var err error
		partyList, err = iamClient.Parties().List(waitCtx, iam.PartyListOptions{})
		return err
	})

	// Get comments for case
	g.Go(func() error {
		var err error
		comments, err = cmsClient.Comments().List(waitCtx, cms.CommentListOptions{
			CaseID: id,
		})
		return err
	})

	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get case type team
	if id != "new" {
		teamRes, err := iamClient.Teams().Get(ctx, kase.TeamID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		team = teamRes
	}

	// Get case recipient
	party, err := iamClient.Parties().Get(ctx, kase.PartyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get case types
	kaseTypes, err = cmsClient.CaseTypes().List(ctx, cms.CaseTypeListOptions{
		PartyTypeIDs: party.PartyTypeIDs,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get case referrals
	referrals, err = cmsClient.Cases().List(ctx, cms.CaseListOptions{
		ParentID: kase.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	qry := req.URL.Query()

	var referralCaseType *cms.CaseType
	if referralCaseTypeID := qry.Get("referralCaseTypeId"); len(referralCaseTypeID) > 0 {
		referralCaseType, err = cmsClient.CaseTypes().Get(ctx, referralCaseTypeID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "case", map[string]interface{}{
		"Case":             kase,
		"CaseTypes":        kaseTypes,
		"Party":            party,
		"Parties":          partyList,
		"ReferralCaseType": referralCaseType,
		"Referrals":        referrals,
		"Team":             team,
		"Comments":         displayComments(comments),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

type displayComment struct {
	*cms.Comment
	TimeAgo string
}

func displayComments(comments *cms.CommentList) []*displayComment {
	var displayComments []*displayComment
	for _, item := range comments.Items {
		c := &displayComment{
			Comment: item,
			TimeAgo: timeago.English.Format(item.CreatedAt),
		}
		displayComments = append(displayComments, c)
	}
	return displayComments
}

func (h *Server) NewCase(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	cmsClient := h.CMSClient(ctx)
	iamClient := h.IAMClient(ctx)

	var caseTypes *cms.CaseTypeList
	var p *iam.PartyList

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		caseTypes, err = cmsClient.CaseTypes().List(waitCtx, cms.CaseTypeListOptions{})
		return err
	})

	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	p, err := iamClient.Parties().List(ctx, listOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var team *iam.Team
	if len(teamID) > 0 {
		team, err = iamClient.Teams().Get(ctx, teamID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "casenew", map[string]interface{}{
		"PartyID":    qry.Get("partyId"),
		"CaseTypeID": qry.Get("caseTypeId"),
		"Team":       team,
		"CaseTypes":  caseTypes,
		"Parties":    p,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Server) PostCase(ctx context.Context, id string, w http.ResponseWriter, req *http.Request) {

	cmsClient := h.CMSClient(ctx)

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	caseTypeId := req.Form.Get("caseTypeId")
	partyId := req.Form.Get("partyId")
	description := req.Form.Get("description")
	done := req.Form.Get("done")
	parentId := req.Form.Get("parentId")
	teamId := req.Form.Get("teamId")

	var kase *cms.Case
	if id == "" {
		var err error
		kase, err = cmsClient.Cases().Create(ctx, &cms.Case{
			CaseTypeID:  caseTypeId,
			PartyID:     partyId,
			Description: description,
			Done:        false,
			ParentID:    parentId,
			TeamID:      teamId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		var err error
		kase, err = cmsClient.Cases().Update(ctx, &cms.Case{
			ID:          id,
			CaseTypeID:  caseTypeId,
			PartyID:     partyId,
			Description: description,
			Done:        done == "on",
			ParentID:    parentId,
			TeamID:      teamId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if len(parentId) > 0 {
		w.Header().Set("Location", "/cases/"+parentId)
	} else {
		w.Header().Set("Location", "/cases/"+kase.ID)
	}
	w.WriteHeader(http.StatusSeeOther)

}
