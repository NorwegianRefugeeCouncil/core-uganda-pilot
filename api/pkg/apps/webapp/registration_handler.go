package webapp

import (
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"net/http"
)

type IndividualRegistrationRequestHandler struct {
	s *Server
	i *iam.Individual
	r *http.Request
}

func NewIndividualRegistrationHandler(s *Server, i *iam.Individual, r *http.Request) *IndividualRegistrationRequestHandler {
	return &IndividualRegistrationRequestHandler{s, i, r}
}

func (irh *IndividualRegistrationRequestHandler) IndividualExists() bool {
	necessaryDataExistsForIndividual := true

	expectedAttributes := []*iam.PartyAttributeDefinition{
		&iam.FullNameAttribute,
		&iam.DisplayNameAttribute,
		&iam.GenderAttribute,
		&iam.DisplacementStatusAttribute,
		&iam.BirthDateAttribute,
	}

	for _, expectedAttribute := range expectedAttributes {
		isNotPresentOrIsNotSet := !irh.i.HasAttribute(expectedAttribute.ID) || (irh.i.HasAttribute(expectedAttribute.ID) && len(irh.i.GetAttribute(expectedAttribute.ID)) == 0)
		if isNotPresentOrIsNotSet {
			necessaryDataExistsForIndividual = false
			break
		}
	}
	return necessaryDataExistsForIndividual
}

func (irh *IndividualRegistrationRequestHandler) GetOpenCases() []*cms.Case {
	cmsClient, err := irh.s.CMSClient(irh.r)
	if err != nil {
		return nil
	}

	notDone := false

	cases, err := cmsClient.Cases().List(irh.r.Context(), cms.CaseListOptions{
		PartyIDs: []string{irh.i.ID},
		Done:     &notDone,
	})
	if err != nil {
		return nil
	}

	return cases.Items
}

func (irh *IndividualRegistrationRequestHandler) GetClosedCases() []*cms.Case {
	cmsClient, err := irh.s.CMSClient(irh.r)
	if err != nil {
		return nil
	}

	done := true

	cases, err := cmsClient.Cases().List(irh.r.Context(), cms.CaseListOptions{
		PartyIDs: []string{irh.i.ID},
		Done:     &done,
	})
	if err != nil {
		return nil
	}

	return cases.Items
}
