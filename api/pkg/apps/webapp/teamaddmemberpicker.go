package webapp

import (
	"encoding/json"
	"github.com/nrc-no/core/pkg/apps/iam"
	"net/http"
	"net/url"
)

type TeamPartyOptions struct {
	PartyTypeID string `json:"partyTypeId"`
	SearchParam string `json:"searchParam"`
	TeamId string `json:"teamId"`
}

func (a *TeamPartyOptions) UnmarshalQueryParameters(values url.Values) error {
	a.PartyTypeID = values.Get("partyTypeId")
	a.SearchParam = values.Get("searchParam")
	a.TeamId = values.Get("teamId")
	return nil
}

func (h *Server) PickTeamParty(w http.ResponseWriter, req *http.Request){
	ctx := req.Context()
	iamClient := h.IAMClient(ctx)

	var listOptions TeamPartyOptions
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	partyList, err := iamClient.Parties().List(ctx, iam.PartyListOptions{
		PartyTypeID: listOptions.PartyTypeID,
		SearchParam: listOptions.SearchParam,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	partiesInTeam, err := iamClient.Memberships().List(ctx, iam.MembershipListOptions{
		TeamID:       listOptions.TeamId,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var returnList iam.PartyList
	for _, party := range partyList.Items {
		isMember := false
		for _, member := range partiesInTeam.Items {
			if party.ID == member.ID {
				isMember = true
			}
		}
		if !isMember {
			returnList.Items = append(returnList.Items, party)
		}
	}

	responseBytes, err := json.Marshal(returnList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}
