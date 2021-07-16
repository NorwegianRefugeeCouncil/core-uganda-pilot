package webapp

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/sessionmanager"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func (s *Server) RelationshipTypes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	iamClient := s.IAMClient(ctx)

	r := &iam.RelationshipType{}

	if req.Method == "POST" {
		s.PostRelationshipType(ctx, r, w, req)
		return
	}

	relationshipTypes, err := iamClient.RelationshipTypes().List(ctx, iam.RelationshipTypeListOptions{})
	if err != nil {
		s.Error(w, err)
		return
	}

	if err := s.renderFactory.New(req).ExecuteTemplate(w, "relationshiptypes", map[string]interface{}{
		"RelationshipTypes": relationshipTypes,
	}); err != nil {
		s.Error(w, err)
		return
	}
}

func (s *Server) NewRelationshipType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	iamClient := s.IAMClient(ctx)

	p, err := iamClient.PartyTypes().List(ctx, iam.PartyTypeListOptions{})
	if err != nil {
		s.Error(w, err)
		return
	}

	if err := s.renderFactory.New(req).ExecuteTemplate(w, "relationshiptype", map[string]interface{}{
		"PartyTypes": p,
	}); err != nil {
		s.Error(w, err)
		return
	}
}

func (s *Server) RelationshipType(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	iamClient := s.IAMClient(ctx)

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id found in path")
		s.Error(w, err)
		return
	}

	r, err := iamClient.RelationshipTypes().Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	p, err := iamClient.PartyTypes().List(ctx, iam.PartyTypeListOptions{})
	if err != nil {
		s.Error(w, err)
		return
	}

	if req.Method == "POST" {
		s.PostRelationshipType(ctx, r, w, req)
		return
	}

	if err := s.renderFactory.New(req).ExecuteTemplate(w, "relationshiptype", map[string]interface{}{
		"RelationshipType": r,
		"PartyTypes":       p,
	}); err != nil {
		s.Error(w, err)
		return
	}
}

func (s *Server) PostRelationshipType(
	ctx context.Context,
	r *iam.RelationshipType,
	w http.ResponseWriter,
	req *http.Request,
) {

	iamClient := s.IAMClient(ctx)

	if err := req.ParseForm(); err != nil {
		s.Error(w, err)
		return
	}

	formValues := req.Form

	isNew := false
	if len(r.ID) == 0 {
		r.ID = uuid.NewV4().String()
		isNew = true
	}

	r.Name = formValues.Get("name")

	if formValues.Get("isDirectional") == "true" {
		r.IsDirectional = true
	} else {
		r.IsDirectional = false
	}
	r.FirstPartyRole = formValues.Get("firstPartyRole")
	r.SecondPartyRole = formValues.Get("secondPartyRole")

	r.Rules = []iam.RelationshipTypeRule{
		{
			PartyTypeRule: &iam.PartyTypeRule{
				FirstPartyTypeID:  formValues.Get("rules[0].firstPartyTypeId"),
				SecondPartyTypeID: formValues.Get("rules[0].secondPartyTypeId"),
			},
		},
	}

	if isNew {
		out, err := iamClient.RelationshipTypes().Create(ctx, r)
		if err != nil {
			s.Error(w, err)
			return
		}
		s.sessionManager.AddNotification(ctx, &sessionmanager.Notification{
			Message: fmt.Sprintf("Relationship type \"%s\" successfully updated", r.Name),
			Theme:   "success",
		})
		w.Header().Set("Location", "/settings/relationshiptypes/"+out.ID)
		w.WriteHeader(http.StatusSeeOther)
	} else {
		out, err := iamClient.RelationshipTypes().Update(ctx, r)
		if err != nil {
			s.Error(w, err)
			return
		}
		s.sessionManager.AddNotification(ctx, &sessionmanager.Notification{
			Message: fmt.Sprintf("Relationship type \"%s\" successfully updated", r.Name),
			Theme:   "success",
		})
		w.Header().Set("Location", "/settings/relationshiptypes/"+out.ID)
		w.WriteHeader(http.StatusSeeOther)
	}

}
