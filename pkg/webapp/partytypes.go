package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/iam"
	"github.com/nrc-no/core/pkg/sessionmanager"
	"net/http"
)

func (s *Server) PartyTypes(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	if req.Method == http.MethodPost {
		s.PostPartyType(ctx, &iam.PartyType{}, w, req)
		return
	}

	partyTypes, err := iamClient.PartyTypes().List(ctx, iam.PartyTypeListOptions{})
	if err != nil {
		s.Error(w, err)
		return
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "partytypes", map[string]interface{}{
		"PartyTypes": partyTypes,
	}); err != nil {
		s.Error(w, err)
		return
	}
}

func (s *Server) PartyType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)

		return
	}

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id found in path")
		s.Error(w, err)

		return
	}

	var partyType = &iam.PartyType{}

	if id != "new" {
		var err error

		partyType, err = iamClient.PartyTypes().Get(ctx, id)
		if err != nil {
			s.Error(w, err)

			return
		}
	}

	if req.Method == http.MethodPost {
		s.PostPartyType(ctx, partyType, w, req)

		return
	}

	notifications, err := s.flashes(req, w)
	if err != nil {
		s.Error(w, err)

		return
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "partytype", map[string]interface{}{
		"PartyType":     partyType,
		"Notifications": notifications,
	}); err != nil {
		s.Error(w, err)
		return
	}
}

func (s *Server) PostPartyType(
	ctx context.Context,
	partyType *iam.PartyType,
	w http.ResponseWriter,
	req *http.Request,
) {
	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	isNew := false
	if len(partyType.ID) == 0 {
		isNew = true
	}

	if err := req.ParseForm(); err != nil {
		s.Error(w, err)
		return
	}

	name := req.Form.Get("name")
	partyType.Name = name
	partyType.IsBuiltIn = req.Form.Get("isBuiltIn") == "true"

	if isNew {
		created, err := iamClient.PartyTypes().Create(ctx, partyType)
		if err != nil {
			s.Error(w, err)

			return
		}

		if err := s.addFlash(req, w, &sessionmanager.FlashMessage{
			Message: fmt.Sprintf("Party type \"%s\" successfully updated", partyType.Name),
			Theme:   "success",
		}); err != nil {
			s.Error(w, err)

			return
		}

		w.Header().Set("Location", "/settings/partytypes/"+created.ID)
		w.WriteHeader(http.StatusSeeOther)

		return
	}

	updated, err := iamClient.PartyTypes().Update(ctx, partyType)
	if err != nil {
		s.Error(w, err)

		return
	}

	if err := s.addFlash(req, w, &sessionmanager.FlashMessage{
		Message: fmt.Sprintf("Party type \"%s\" successfully updated", partyType.Name),
		Theme:   "success",
	}); err != nil {
		s.Error(w, err)

		return
	}

	w.Header().Set("Location", "/settings/partytypes/"+updated.ID)
	w.WriteHeader(http.StatusSeeOther)
}
