package iam

import (
	"net/http"
)

func (s *Server) PutAttribute(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	var id string

	if !s.GetPathParam("id", w, req, &id) {
		return
	}

	var payload Attribute
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	attribute, err := s.AttributeStore.Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	attribute.Name = payload.Name
	attribute.Translations = payload.Translations
	attribute.PartyTypeIDs = payload.PartyTypeIDs
	attribute.IsPersonallyIdentifiableInfo = payload.IsPersonallyIdentifiableInfo

	if err := s.AttributeStore.Update(ctx, attribute); err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, payload)

}
