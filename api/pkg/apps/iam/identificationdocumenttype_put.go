package iam

import "net/http"

func (s *Server) putIdentificationDocumentType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	var payload IdentificationDocumentType
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	identificationDocumentType, err := s.identificationDocumentTypeStore.get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}

	identificationDocumentType.Name = payload.Name

	if err := s.identificationDocumentTypeStore.update(ctx, identificationDocumentType); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, payload)
}
