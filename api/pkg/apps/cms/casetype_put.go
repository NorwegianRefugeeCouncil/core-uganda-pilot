package cms

import "net/http"

func (s *Server) PutCaseType(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var id string
	if !s.GetPathParam("id", w, req, &id) {
		return
	}

	var payload CaseType
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	caseType, err := s.caseTypeStore.Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	caseType.Name = payload.Name
	caseType.PartyTypeID = payload.PartyTypeID
	caseType.TeamID = payload.TeamID
	caseType.Template = payload.Template

	if err := s.caseTypeStore.Update(ctx, caseType); err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, caseType)

}
