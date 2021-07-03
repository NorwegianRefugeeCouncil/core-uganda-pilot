package webapp

import (
	"github.com/nrc-no/core/pkg/apps/cms"
	"net/http"
)

func (s *Server) PostComment(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	caseID := req.Form.Get("case_id")
	body := req.Form.Get("body")
	comment := cms.Comment{
		CaseID: caseID,
		Body:   body,
	}

	_, err := s.CMSClient(ctx).Comments().Create(ctx, &comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/cases/"+caseID, http.StatusSeeOther)

}