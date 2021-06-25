package iam

import (
	"net/http"
)

func (s *Server) ListRelationships(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	qry := req.URL.Query()

	listOptions := RelationshipListOptions{
		RelationshipTypeID: qry.Get("relationshipTypeId"),
		FirstPartyId:       qry.Get("firstPartyId"),
		SecondParty:        qry.Get("secondPartyId"),
		EitherParty:        qry.Get("eitherPartyId"),
	}

	ret, err := s.RelationshipStore.List(ctx, listOptions)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, ret)
}
