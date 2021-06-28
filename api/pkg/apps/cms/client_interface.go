package cms

import (
	"context"
	"net/url"
)

type Interface interface {
	Cases() CaseClient
	CaseTypes() CaseTypeClient
}

type CaseListOptions struct {
	PartyID    string
	CaseTypeID string
	ParentID   string
}

func (a *CaseListOptions) MarshalQueryParameters() (url.Values, error) {
	ret := url.Values{}
	if len(a.PartyID) > 0 {
		ret.Set("partyId", a.PartyID)
	}
	if len(a.CaseTypeID) > 0 {
		ret.Set("caseTypeId", a.CaseTypeID)
	}
	if len(a.ParentID) > 0 {
		ret.Set("parentId", a.ParentID)
	}
	return ret, nil
}

func (a *CaseListOptions) UnmarshalQueryParameters(values url.Values) error {
	a.PartyID = values.Get("partyId")
	a.CaseTypeID = values.Get("caseTypeId")
	a.ParentID = values.Get("parentId")
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
