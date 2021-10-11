package webapp

import (
	"github.com/nrc-no/core/pkg/cms"
	"net/http"
)

func (s *Server) PostComment(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if err := req.ParseForm(); err != nil {
		s.Error(w, err)
		return
	}

	caseID := req.Form.Get("case_id")
	body := req.Form.Get("body")
	comment := cms.Comment{
		CaseID: caseID,
		Body:   body,
	}

	cmsClient, err := s.CMSClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	_, err = cmsClient.Comments().Create(ctx, &comment)
	if err != nil {
		s.Error(w, err)
		return
	}

	http.Redirect(w, req, "/cases/"+caseID, http.StatusSeeOther)

}
