package webapp

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/form"
	"github.com/nrc-no/core/pkg/i18n"
	"github.com/nrc-no/core/pkg/iam"
	"github.com/nrc-no/core/pkg/sessionmanager"
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) Attributes(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if req.Method == http.MethodPost {
		s.PostAttribute(ctx, &iam.PartyAttributeDefinition{}, w, req)

		return
	}

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)

		return
	}

	countryID := s.GetCountryFromLoginUser(w, req)

	list, err := iamClient.PartyAttributeDefinitions().List(ctx, iam.PartyAttributeDefinitionListOptions{
		CountryIDs: []string{iam.GlobalCountry.ID, countryID},
	})
	if err != nil {
		s.Error(w, err)

		return
	}

	partyTypes, err := iamClient.PartyTypes().List(ctx, iam.PartyTypeListOptions{})
	if err != nil {
		s.Error(w, err)

		return
	}

	notifications, err := s.flashes(req, w)
	if err != nil {
		s.Error(w, err)

		return
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "attributes", map[string]interface{}{
		"Attributes":    list,
		"PartyTypes":    partyTypes,
		"Notifications": notifications,
	}); err != nil {
		s.Error(w, err)

		return
	}
}

func (s *Server) NewAttribute(w http.ResponseWriter, req *http.Request) {
	notifications, err := s.flashes(req, w)
	if err != nil {
		s.Error(w, err)
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "attribute", map[string]interface{}{
		"PartyTypes": iam.PartyTypeList{
			Items: []*iam.PartyType{
				&iam.IndividualPartyType,
			},
		},
		"ControlTypes":  form.ControlTypes,
		"Notifications": notifications,
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

	if req.Method == http.MethodPost {
		s.PostAttribute(ctx, partyAttributeDefinition, w, req)

		return
	}

	notifications, err := s.flashes(req, w)
	if err != nil {
		s.Error(w, err)

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
		"Notifications": notifications,
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

	// form control
	name := values.Get("name")
	// TODO add dynamic locale
	label := i18n.Strings{{"en", values.Get("label")}}
	required := values.Get("isRequired") == "on"
	controlType := form.ControlType(values["controlType"][0])
	attribute.FormControl = *form.NewControl(name, controlType, label, required)

	// TODO infer country from subject
	// subject, ok := ctx.Value("Subject").(string)
	//if !ok {
	//	s.Error(w, errors.New("couldn't get subject id string"))
	//	return
	//}
	//user, err := iamClient.Individuals().Get(ctx, subject)
	//if err != nil {
	//	s.Error(w, err)
	//	return
	//}
	attribute.CountryID = iam.UgandaCountry.ID

	var storedAttribute *iam.PartyAttributeDefinition

	if isNew {
		storedAttribute, err = iamClient.PartyAttributeDefinitions().Create(ctx, attribute)
	} else {
		storedAttribute, err = iamClient.PartyAttributeDefinitions().Update(ctx, attribute)
	}

	if err != nil {
		if status, ok := err.(*validation.Status); ok {
			validatedElements := zipAttributeAndErrors(attribute, status.Errors)
			s.json(w, status.Code, validatedElements)
		} else {
			s.Error(w, err)
		}

		return
	}

	err = s.addFlash(req, w, &sessionmanager.FlashMessage{
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

// zipAttributeAndErrors returns a form.Form containing the validation information, ie the faulty form elements only
func zipAttributeAndErrors(attribute *iam.PartyAttributeDefinition, errorList validation.ErrorList) form.Form {
	var result form.Form

	var errs *validation.ErrorList

	ctrl := attribute.FormControl

	// name
	if errs = errorList.FindFamily("name"); errs.Length() > 0 {
		ctrl.Errors = errs
		result.Controls = append(result.Controls, ctrl)
	}
	// partyTypeIds
	if errs = errorList.FindFamily("partyTypeIds"); errs.Length() > 0 {
		ctrl.Errors = errs
		result.Controls = append(result.Controls, ctrl)
	}

	result.Controls = []form.Control{ctrl}

	return result
}
