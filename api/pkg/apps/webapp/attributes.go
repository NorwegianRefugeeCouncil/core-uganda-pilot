package webapp

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/form"
	"github.com/nrc-no/core/pkg/sessionmanager"
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"strings"
)

func (s *Server) Attributes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	if req.Method == "POST" {
		s.PostAttribute(ctx, &iam.Attribute{}, w, req)
		return
	}

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	list, err := iamClient.Attributes().List(ctx, iam.AttributeListOptions{})
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
		err := fmt.Errorf("No id in path")
		s.Error(w, err)
		return
	}

	attribute, err := iamClient.Attributes().Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	if req.Method == "POST" {
		s.PostAttribute(ctx, attribute, w, req)
		return
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "attribute", map[string]interface{}{
		"Attribute": attribute,
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

func (s *Server) PostAttribute(ctx context.Context, attribute *iam.Attribute, w http.ResponseWriter, req *http.Request) {
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

	translationMap := s.makeTranslationMap(values, w)

	var translations []iam.AttributeTranslation
	for _, translation := range translationMap {
		translations = append(translations, *translation)
	}

	isNew := false
	if len(attribute.ID) == 0 {
		attribute.ID = uuid.NewV4().String()
		isNew = true
	}

	attribute.Name = values.Get("name")
	attribute.PartyTypeIDs = values["partyTypeIds"]
	attribute.Translations = translations
	attribute.IsPersonallyIdentifiableInfo = values.Get("isPii") == "on"

	var storageAction string
	var storedAttribute *iam.Attribute
	if isNew {
		storedAttribute, err = iamClient.Attributes().Create(ctx, attribute)
		storageAction = "created"
	} else {
		storedAttribute, err = iamClient.Attributes().Update(ctx, attribute)
		storageAction = "updated"
	}
	if err != nil {
		if status, ok := err.(*validation.Status); ok {
			validatedElements := zipNewAttributeFormElementsAndErrors(attribute, status.Errors)
			s.json(w, status.Code, validatedElements)
		} else {
			s.Error(w, err)
		}
		return
	}

	err = s.sessionManager.AddNotification(req, w, &sessionmanager.Notification{
		Message: fmt.Sprintf("Attribute \"%s\" successfully %s.", attribute.Name, storageAction),
		Theme:   "success",
	})
	if err != nil {
		s.Error(w, err)
		return
	}

	w.Header().Set("Location", "/settings/attributes/"+storedAttribute.ID)
	w.WriteHeader(http.StatusSeeOther)

}

func zipNewAttributeFormElementsAndErrors(attribute *iam.Attribute, errorList validation.ErrorList) []form.FormElement {
	var validated []form.FormElement
	var errs *validation.ErrorList

	// name
	if errs = errorList.FindFamily("name"); errs.Length() > 0 {
		validated = append(validated, form.FormElement{
			Type:       form.Text,
			Attributes: form.FormElementAttributes{Name: "name"},
			Errors:     errs,
		})
	}
	// partyTypeIds
	if errs = errorList.FindFamily("partyTypeIds"); errs.Length() > 0 {
		validated = append(validated, form.FormElement{
			Type:       form.Dropdown,
			Attributes: form.FormElementAttributes{Name: "partyTypeId"},
			Errors:     errs,
		})
	}
	// translations
	if errs = errorList.Find("translations"); errs.Length() == 1 {
		validated = append(validated, form.FormElement{
			Type:       form.Dropdown,
			Attributes: form.FormElementAttributes{Name: "language-picker"},
			Errors:     errs,
		})
	} else if errs = errorList.FindFamily("translations"); errs.Length() > 0 {
		for _, translation := range attribute.Translations {
			short := fmt.Sprintf("translation.%s.short", translation.Locale)
			long := fmt.Sprintf("translation.%s.long", translation.Locale)
			for _, s := range []string{short, long} {
				if e := errs.FindFamily(s); e != nil {
					validated = append(validated, form.FormElement{
						Type:       form.Text,
						Attributes: form.FormElementAttributes{Name: s},
						Errors:     e,
					})
				}
			}
		}
	}

	return validated
}

func (s *Server) makeTranslationMap(values url.Values, w http.ResponseWriter) map[string]*iam.AttributeTranslation {
	translationMap := map[string]*iam.AttributeTranslation{}

	unexpectedKeyErrMsg := "unexpected translation key. Expected 'translation.{locale}.{short/long}' format"

	for key, v := range values {
		if !strings.HasPrefix(key, "translations.") {
			continue
		}
		parts := strings.Split(key, ".")
		if len(parts) != 3 {
			s.Error(w, errors.New(unexpectedKeyErrMsg))
			return nil
		}

		locale := parts[1]
		part := parts[2]

		if part != "long" && part != "short" {
			s.Error(w, errors.New(unexpectedKeyErrMsg))
			return nil
		}

		if _, ok := translationMap[locale]; !ok {
			translationMap[locale] = &iam.AttributeTranslation{
				Locale: locale,
			}
		}
		t := translationMap[locale]

		if part == "long" {
			t.LongFormulation = v[0]
		} else if part == "short" {
			t.ShortFormulation = v[0]
		} else {
			s.Error(w, errors.New(unexpectedKeyErrMsg))
			return nil
		}

	}
	return translationMap
}
