package iam

import (
	"net/http"
)

func (s *Server) putAttribute(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	var payload Attribute
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	attribute, err := s.attributeStore.get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}

	attribute.Name = payload.Name
	attribute.Translations = payload.Translations
	attribute.PartyTypeIDs = payload.PartyTypeIDs
	attribute.IsPersonallyIdentifiableInfo = payload.IsPersonallyIdentifiableInfo

	if err := s.attributeStore.update(ctx, attribute); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, payload)

}
