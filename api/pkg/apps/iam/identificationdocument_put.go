package iam

import "net/http"

func (s *Server) putIdentificationDocument(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	var payload IdentificationDocument
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	identificationDocument, err := s.identificationDocumentStore.get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}

	identificationDocument.PartyID = payload.PartyID
	identificationDocument.DocumentNumber = payload.DocumentNumber
	identificationDocument.IdentificationDocumentTypeID = payload.IdentificationDocumentTypeID

	if err := s.identificationDocumentStore.update(ctx, identificationDocument); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, payload)
}
