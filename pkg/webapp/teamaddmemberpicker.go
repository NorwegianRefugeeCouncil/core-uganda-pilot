package webapp

import (
	"encoding/json"
	iam2 "github.com/nrc-no/core/pkg/iam"
	"net/http"
	"net/url"
)

type TeamPartyOptions struct {
	PartyTypeID string `json:"partyTypeId"`
	SearchParam string `json:"searchParam"`
	TeamId      string `json:"teamId"`
}

func (a *TeamPartyOptions) UnmarshalQueryParameters(values url.Values) error {
	a.SearchParam = values.Get("searchParam")
	a.TeamId = values.Get("teamId")
	return nil
}

func (s *Server) PickTeamParty(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	var listOptions TeamPartyOptions
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.Error(w, err)
		return
	}

	partyList, err := iamClient.Parties().List(ctx, iam2.PartyListOptions{
		PartyTypeID: iam2.IndividualPartyType.ID,
		SearchParam: listOptions.SearchParam,
	})
	if err != nil {
		s.Error(w, err)
		return
	}

	partiesInTeam, err := iamClient.Memberships().List(ctx, iam2.MembershipListOptions{
		TeamID: listOptions.TeamId,
	})
	if err != nil {
		s.Error(w, err)
		return
	}

	var returnList iam2.PartyList
	returnList.Items = []*iam2.Party{}
	for _, party := range partyList.Items {
		isMember := false
		for _, member := range partiesInTeam.Items {
			if party.ID == member.IndividualID {
				isMember = true
			}
		}
		if !isMember {
			returnList.Items = append(returnList.Items, party)
		}
	}

	responseBytes, err := json.Marshal(returnList)
	if err != nil {
		s.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}
