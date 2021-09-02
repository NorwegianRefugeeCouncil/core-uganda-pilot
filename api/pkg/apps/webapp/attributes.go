package webapp

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/form"
	"github.com/nrc-no/core/pkg/sessionmanager"
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) Attributes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	if req.Method == "POST" {
		s.PostAttribute(ctx, &iam.PartyAttributeDefinition{}, w, req)
		return
	}

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	list, err := iamClient.PartyAttributeDefinitions().List(ctx, iam.PartyAttributeDefinitionListOptions{})
	if err != nil {
		s.Error(w, err)
		return
	}

	partyTypes, err := iamClient.PartyTypes().List(ctx, iam.PartyTypeListOptions{})
	if err != nil {
		s.Error(w, err)
		return
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "attributes", map[string]interface{}{
		"Attributes": list,
		"PartyTypes": partyTypes,
	}); err != nil {
		s.Error(w, err)
		return
	}
}

func (s *Server) NewAttribute(w http.ResponseWriter, req *http.Request) {
	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "attribute", map[string]interface{}{
		"PartyTypes": iam.PartyTypeList{
			Items: []*iam.PartyType{
				&iam.IndividualPartyType,
			},
		},
	}); err != nil {
		s.Error(w, err)
		return
	}
}

func (s *Server) Attribute(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := errors.New("no id in path")
		s.Error(w, err)
		return
	}

	partyAttributeDefinition, err := iamClient.PartyAttributeDefinitions().Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	if req.Method == "POST" {
		s.PostAttribute(ctx, partyAttributeDefinition, w, req)
		return
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "attribute", map[string]interface{}{
		"PartyAttributeDefinition": partyAttributeDefinition,
		"ControlTypes":             form.ControlTypes,
		"PartyTypes": iam.PartyTypeList{
			Items: []*iam.PartyType{
				&iam.IndividualPartyType,
			},
		},
	}); err != nil {
		s.Error(w, err)
		return
	}

}

func (s *Server) PostAttribute(ctx context.Context, attribute *iam.PartyAttributeDefinition, w http.ResponseWriter, req *http.Request) {
	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	if err := req.ParseForm(); err != nil {
		s.Error(w, err)
		return
	}

	values := req.Form

	isNew := false
	if len(attribute.ID) == 0 {
		attribute.ID = uuid.NewV4().String()
		isNew = true
	}

	attribute.ID = values.Get("id")
	attribute.PartyTypeIDs = values["partyTypeIds"]
	attribute.IsPersonallyIdentifiableInfo = values.Get("isPii") == "on"

	var storedAttribute *iam.PartyAttributeDefinition
	if isNew {
		storedAttribute, err = iamClient.PartyAttributeDefinitions().Create(ctx, attribute)
	} else {
		storedAttribute, err = iamClient.PartyAttributeDefinitions().Update(ctx, attribute)
	}
	if err != nil {
		if status, ok := err.(*validation.Status); ok {
			validatedElements := zipNewAttributeFormcontrolsAndErrors(attribute, status.Errors)
			s.json(w, status.Code, validatedElements)
		} else {
			s.Error(w, err)
		}
		return
	}

	err = s.sessionManager.AddNotification(req, w, &sessionmanager.Notification{
		Message: "New attribute definition successfully saved",
		Theme:   "success",
	})
	if err != nil {
		s.Error(w, err)
		return
	}

	w.Header().Set("Location", "/settings/attributes/"+storedAttribute.ID)
	w.WriteHeader(http.StatusSeeOther)

}

func zipNewAttributeFormcontrolsAndErrors(attribute *iam.PartyAttributeDefinition, errorList validation.ErrorList) form.ValuedForm {
	var submittedForm form.ValuedForm
	var validated []form.ValuedControl
	var errs *validation.ErrorList

	// name
	if errs = errorList.FindFamily("id"); errs.Length() > 0 {
		validated = append(validated, form.ValuedControl{
			Control: &attribute.FormControl,
			Errors:  errs,
		})
	}
	// partyTypeIds
	if errs = errorList.FindFamily("partyTypeIds"); errs.Length() > 0 {
		validated = append(validated, form.ValuedControl{
			Control: &attribute.FormControl,
			Errors:  errs,
		})
	}

	submittedForm.Controls = validated
	return submittedForm
}
