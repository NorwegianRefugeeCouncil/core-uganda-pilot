package webapp

import (
	"fmt"
	"github.com/nrc-no/core/pkg/cms"
	"github.com/nrc-no/core/pkg/iam"
	"github.com/nrc-no/core/pkg/teamstatusctrl"
	"net/http"
)

func (s *Server) GetTeamStatusController(req *http.Request, individual *iam.Individual) (teamstatusctrl.TeamStatusCtrl, error) {
	tsc := teamstatusctrl.TeamStatusCtrl{}

	// Get Cases For Individual
	cmsClient, err := s.CMSClient(req)
	if err != nil {
		return tsc, err
	}
	cases, err := cmsClient.Cases().List(req.Context(), cms.CaseListOptions{
		PartyIDs: []string{individual.ID},
	})
	if err != nil {
		return tsc, err
	}

	// Get currently logged in user
	session, err := s.sessionManager.Get(req)
	if err != nil {
		return tsc, err
	}

	profile := session.Values["profile"]

	userSessionClaims, ok := profile.(*Claims)
	if !ok {
		return tsc, fmt.Errorf("Failed to cast profile to claims")
	}

	// Get memberships For currently logged in user
	iamClient, err := s.IAMClient(req)
	if err != nil {
		return tsc, err
	}
	memberships, err := iamClient.Memberships().List(req.Context(), iam.MembershipListOptions{
		IndividualID: userSessionClaims.Subject,
	})
	if err != nil {
		return tsc, err
	}

	// Infer list of team ids from memberships
	teamIdsForIndividual := []string{}
	for _, membership := range memberships.Items {
		teamIdAlreadyInList := false
		for _, tid := range teamIdsForIndividual {
			if tid == membership.TeamID {
				teamIdAlreadyInList = true
				break
			}
		}

		if !teamIdAlreadyInList {
			teamIdsForIndividual = append(teamIdsForIndividual, membership.TeamID)
		}
	}

	// Get team intake case types for currently logged in user
	teamIntakeCaseTypes := []*cms.CaseType{}
	caseTypes, err := cmsClient.CaseTypes().List(req.Context(), cms.CaseTypeListOptions{
		TeamIDs: teamIdsForIndividual,
	})
	if err != nil {
		return tsc, nil
	}
	for _, ct := range caseTypes.Items {
		if ct.IntakeCaseType {
			ctAlreadyInList := false

			for _, tict := range teamIntakeCaseTypes {
				if tict.ID == ct.ID {
					ctAlreadyInList = true
				}
			}
			if !ctAlreadyInList {
				teamIntakeCaseTypes = append(teamIntakeCaseTypes, ct)
			}
		}
	}

	tsc.TeamIntakeCaseTypes = teamIntakeCaseTypes
	tsc.CasesForIndividual = cases.Items

	return tsc, nil
}
