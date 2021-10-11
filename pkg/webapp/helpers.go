package webapp

import (
	"github.com/nrc-no/core/internal/sessionmanager"
	"github.com/nrc-no/core/pkg/cms"
	iam2 "github.com/nrc-no/core/pkg/iam"
	"net/http"
)

func (s *Server) retrieveParties(req *http.Request) (*iam2.PartyList, error) {
	cli, err := s.IAMClient(req)
	if err != nil {
		return nil, err
	}
	parties, err := cli.Parties().List(req.Context(), iam2.PartyListOptions{})
	if err != nil {
		return nil, err
	}
	return parties, nil
}

func (s *Server) retrieveTeams(req *http.Request) (*iam2.TeamList, error) {
	cli, err := s.IAMClient(req)
	if err != nil {
		return nil, err
	}
	teams, err := cli.Teams().List(req.Context(), iam2.TeamListOptions{})
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (s *Server) retrieveCaseTypes(req *http.Request) (*cms.CaseTypeList, error) {
	cli, err := s.CMSClient(req)
	if err != nil {
		return nil, err
	}
	caseTypes, err := cli.CaseTypes().List(req.Context(), cms.CaseTypeListOptions{})
	if err != nil {
		return nil, err
	}
	return caseTypes, nil
}

func (s *Server) validationErrorNotification(req *http.Request, w http.ResponseWriter) {
	if err := s.sessionManager.AddNotification(req, w, &sessionmanager.Notification{
		Message: "There seems to be an problem with the data you have submitted. See below for errors.",
		Theme:   "danger",
	}); err != nil {
		s.Error(w, err)
	}
}
