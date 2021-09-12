package webapp

import (
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) GetProfile(w http.ResponseWriter, req *http.Request)(*Claims) {
	session, err := s.sessionManager.Get(req)
	if err != nil {
		s.Error(w, err)
	}

	profileIntf, ok := session.Values["profile"]
	if !ok {
		logrus.Errorf("profile not found")
	}

	claims, ok := profileIntf.(*Claims)
	if !ok {
		logrus.Errorf("claims not found")
	}
	return claims
}

func (s *Server) GetTeamFromLoginUser(w http.ResponseWriter, req *http.Request) string {
	ctx := req.Context()

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return ""
	}

	claims := s.GetProfile(w, req)

	membership, err := iamClient.Memberships().List(ctx, iam.MembershipListOptions{
		IndividualID: claims.Subject,
	})

	var teamID string
	if len((*membership).Items) == 1 {
		teamID = (*membership).Items[0].TeamID
	} else {
		logrus.Errorf("User %s has no team or more than one team", claims.Subject)
	}
	return teamID
}

func (s *Server) GetCountryFromLoginUser(w http.ResponseWriter, req *http.Request) string {
	ctx := req.Context()

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return ""
	}
	teamID := s.GetTeamFromLoginUser(w, req)
	nationality, err:= iamClient.Nationalities().List(ctx, iam.NationalityListOptions{TeamID: teamID})

	var countryID string
	if len((*nationality).Items) == 1 {
		countryID = (*nationality).Items[0].CountryID
	} else {
		logrus.Errorf("Team %s has no nationality or more than one nationality", teamID)
	}
	return countryID

}