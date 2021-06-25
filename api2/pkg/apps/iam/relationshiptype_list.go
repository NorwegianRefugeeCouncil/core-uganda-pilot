package iam

import (
	"net/http"
)

func (s *Server) ListRelationshipTypes(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	qry := req.URL.Query()

	listOptions := RelationshipTypeListOptions{
		PartyType: qry.Get("partyType"),
	}

	ret, err := s.RelationshipTypeStore.List(ctx, listOptions)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, ret)
}
