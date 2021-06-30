package cms

import (
	"context"
	"net/url"
)

type Interface interface {
	Cases() CaseClient
	CaseTypes() CaseTypeClient
	Comments() CommentClient
}

type CaseListOptions struct {
	PartyIDs    []string
	TeamIDs     []string
	CaseTypeIDs []string
	ParentID    string
	Done        *bool
}

func (a *CaseListOptions) MarshalQueryParameters() (url.Values, error) {
	ret := url.Values{}
	for _, partyID := range a.PartyIDs {
		ret.Add("partyId", partyID)
	}
	for _, caseTypeID := range a.CaseTypeIDs {
		ret.Add("caseTypeId", caseTypeID)
	}
	for _, teamID := range a.TeamIDs {
		ret.Add("teamId", teamID)
	}
	if len(a.ParentID) > 0 {
		ret.Set("parentId", a.ParentID)
	}
	if a.Done != nil {
		if *a.Done {
			ret.Set("done", "true")
		} else {
			ret.Set("done", "false")
		}
	}
	return ret, nil
}

func (a *CaseListOptions) UnmarshalQueryParameters(values url.Values) error {
	a.PartyIDs = values["partyId"]
	a.CaseTypeIDs = values["caseTypeId"]
	a.TeamIDs = values["teamId"]
	a.ParentID = values.Get("parentId")
	doneStr := values.Get("done")
	if len(doneStr) > 0 {
		isDone := doneStr == "true"
		a.Done = &isDone
	}
	return nil
}

type CaseClient interface {
	Get(ctx context.Context, id string) (*Case, error)
	Create(ctx context.Context, create *Case) (*Case, error)
	Update(ctx context.Context, update *Case) (*Case, error)
	List(ctx context.Context, listOptions CaseListOptions) (*CaseList, error)
}

type CaseTypeListOptions struct {
	PartyTypeIDs []string
}

func (a *CaseTypeListOptions) MarshalQueryParameters() (url.Values, error) {
	ret := url.Values{}
	for _, partyTypeID := range a.PartyTypeIDs {
		ret.Add("partyTypeId", partyTypeID)
	}
	return ret, nil
}

func (a *CaseTypeListOptions) UnmarshalQueryParameters(values url.Values) error {
	a.PartyTypeIDs = values["partyTypeId"]
	return nil
}

type CaseTypeClient interface {
	Get(ctx context.Context, id string) (*CaseType, error)
	Create(ctx context.Context, create *CaseType) (*CaseType, error)
	Update(ctx context.Context, update *CaseType) (*CaseType, error)
	List(ctx context.Context, listOptions CaseTypeListOptions) (*CaseTypeList, error)
}

type CommentListOptions struct {
	CaseID string
}

func (a *CommentListOptions) MarshalQueryParameters() (url.Values, error) {
	ret := url.Values{}
	ret.Set("caseId", a.CaseID)
	return ret, nil
}

func (a *CommentListOptions) UnmarshalQueryParameters(values url.Values) error {
	a.CaseID = values.Get("caseId")
	return nil
}

type CommentClient interface {
	Get(ctx context.Context, id string) (*Comment, error)
	Create(ctx context.Context, create *Comment) (*Comment, error)
	Update(ctx context.Context, update *Comment) (*Comment, error)
	List(ctx context.Context, listOptions CommentListOptions) (*CommentList, error)
	Delete(ctx context.Context, id string) error
}
