package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/apps/seeder"
	"github.com/nrc-no/core/pkg/form"
	"github.com/nrc-no/core/pkg/registrationctrl"
	"github.com/nrc-no/core/pkg/sessionmanager"
	"github.com/nrc-no/core/pkg/teamstatusctrl"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/nrc-no/core/pkg/validation"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"strings"
)

type IndividualWithStatuses struct {
	Individual         *iam.Individual
	RegistrationStatus *registrationctrl.Status
	TeamStatusActions  []teamstatusctrl.TeamStatusAction
}

func (s *Server) Individuals(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	if req.Method == "POST" {
		s.Individual(w, req)
		return
	}

	var listOptions iam.IndividualListOptions
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.Error(w, err)
		return
	}
	listOptions.PartyTypeIDs = []string{iam.BeneficiaryPartyType.ID}

	list, err := iamClient.Individuals().List(ctx, listOptions)
	if err != nil {
		s.Error(w, err)
		return
	}

	var beneficiaryWithStatuses []IndividualWithStatuses
	for _, b := range list.Items {
		rc, err := s.GetRegistrationController(w, req, b)
		if err != nil {
			s.Error(w, err)
			return
		}

		tsc, err := s.GetTeamStatusController(req, b)
		if err != nil {
			s.Error(w, err)
			return
		}

		beneficiaryWithStatuses = append(beneficiaryWithStatuses, IndividualWithStatuses{
			Individual:         b,
			RegistrationStatus: rc.Status(),
			TeamStatusActions:  tsc.GetTeamStatusActions(),
		})
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "individuals", map[string]interface{}{
		"Individuals":             list,
		"IndividualsWithStatuses": beneficiaryWithStatuses,
		"Page":                    "list",
	}); err != nil {
		s.Error(w, err)
		return
	}

}

func (s *Server) IndividualCredentials(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		s.Error(w, err)
		return
	}

	individual, err := iamClient.Individuals().Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	if req.Method == "POST" {
		s.PostIndividualCredentials(w, req, individual.ID)
		return
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "individual_credentials", map[string]interface{}{
		"Page":       "credentials",
		"Individual": individual,
	}); err != nil {
		s.Error(w, err)
		return
	}

}

func (s *Server) PostIndividualCredentials(w http.ResponseWriter, req *http.Request, partyID string) {
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		s.Error(w, err)
		return
	}
	values := req.Form
	password := values.Get("password")

	if err := s.login.Login().SetCredentials(ctx, partyID, password); err != nil {
		s.Error(w, err)
		return
	}

	http.Redirect(w, req, "/individuals/"+partyID+"/credentials", http.StatusSeeOther)
}

func (s *Server) Individual(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	cmsClient, err := s.CMSClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	id, ok := mux.Vars(req)["id"]
	if (!ok || len(id) == 0) && req.Method != "POST" {
		err := fmt.Errorf("no id in path")
		s.Error(w, err)
		return
	}

	var individual *iam.Individual
	var parties *iam.IndividualList
	var caseTypes *cms.CaseTypeList
	var cList *cms.CaseList
	var partyTypes *iam.PartyTypeList
	var relationshipsForIndividual *iam.RelationshipList
	var relationshipTypes *iam.RelationshipTypeList
	var attrs *iam.AttributeList
	var teams *iam.TeamList
	var individualAssessment *cms.Case
	var situationAnalysis *cms.Case
	var identificationDocuments *iam.IdentificationDocumentList
	var identificationDocumentTypes *iam.IdentificationDocumentTypeList
	var identificationDocumentTypesMap = map[string]string{}

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if id == "new" {
			individual = iam.NewIndividual("")
			return nil
		}
		var err error
		individual, err = iamClient.Individuals().Get(waitCtx, id)
		return err
	})

	g.Go(func() error {
		var err error
		parties, err = iamClient.Individuals().List(waitCtx, iam.IndividualListOptions{})
		return err
	})

	g.Go(func() error {
		var err error
		teams, err = iamClient.Teams().List(waitCtx, iam.TeamListOptions{})
		return err
	})

	g.Go(func() error {
		var err error
		partyTypes, err = iamClient.PartyTypes().List(waitCtx, iam.PartyTypeListOptions{})
		return err
	})

	g.Go(func() error {
		var err error
		identificationDocumentTypes, err = iamClient.IdentificationDocumentTypes().List(waitCtx, iam.IdentificationDocumentTypeListOptions{})
		for _, idt := range identificationDocumentTypes.Items {
			identificationDocumentTypesMap[idt.ID] = idt.Name
		}
		return err
	})

	g.Go(func() error {
		if id == "new" {
			relationshipsForIndividual = &iam.RelationshipList{
				Items: []*iam.Relationship{},
			}
			return nil
		}
		var err error
		relationshipsForIndividual, err = iamClient.Relationships().List(waitCtx, iam.RelationshipListOptions{EitherPartyID: id})
		identificationDocuments, err = iamClient.IdentificationDocuments().List(waitCtx, iam.IdentificationDocumentListOptions{PartyIDs: []string{id}})
		return err
	})

	g.Go(func() error {
		var err error
		relationshipTypes, err = iamClient.RelationshipTypes().List(waitCtx, iam.RelationshipTypeListOptions{
			PartyTypeID: iam.IndividualPartyType.ID,
		})
		return err
	})

	g.Go(func() error {
		var err error
		attrs, err = iamClient.Attributes().List(waitCtx, iam.AttributeListOptions{
			PartyTypeIDs: []string{iam.IndividualPartyType.ID},
		})
		return err
	})

	g.Go(func() error {
		var err error
		caseTypes, err = cmsClient.CaseTypes().List(ctx, cms.CaseTypeListOptions{})
		return err
	})

	g.Go(func() error {
		var err error
		cList, err = cmsClient.Cases().List(ctx, cms.CaseListOptions{PartyIDs: []string{id}})
		return err
	})

	g.Go(func() error {
		var err error
		returnedCases, err := cmsClient.Cases().List(ctx, cms.CaseListOptions{
			PartyIDs:    []string{id},
			CaseTypeIDs: []string{seeder.UGIndividualResponseCaseType.ID},
		})
		if err != nil {
			return err
		}
		if len(returnedCases.Items) == 1 {
			individualAssessment = returnedCases.Items[0]
		}
		return err
	})

	g.Go(func() error {
		var err error
		returnedCases, err := cmsClient.Cases().List(ctx, cms.CaseListOptions{
			PartyIDs:    []string{id},
			CaseTypeIDs: []string{seeder.UGSituationalAnalysisCaseType.ID},
		})
		if err != nil {
			return err
		}
		if len(returnedCases.Items) == 1 {
			situationAnalysis = returnedCases.Items[0]
		}
		return err
	})

	g.Go(func() error {
		var err error
		cList, err = cmsClient.Cases().List(ctx, cms.CaseListOptions{PartyIDs: []string{id}})
		return err
	})

	if err := g.Wait(); err != nil {
		s.Error(w, err)
		return
	}

	filteredRelationshipTypes := PrepRelationshipTypeDropdown(relationshipTypes)

	if req.Method == "POST" {
		individual, err = s.PostIndividual(ctx, attrs, id, w, req)
		if err != nil {
			if status, ok := err.(*validation.Status); ok {
				validatedAttrs := zipAttributesAndErrors(&status.Errors, attrs)
				s.json(w, status.Code, validatedAttrs)
			} else {
				s.Error(w, err)
			}
			return
		}
	}

	type DisplayCase struct {
		Case     *cms.Case
		CaseType *cms.CaseType
	}

	ctMap := map[string]*cms.CaseType{}
	for _, item := range caseTypes.Items {
		ctMap[item.ID] = item
	}

	var displayCases []*DisplayCase
	for _, item := range cList.Items {
		d := DisplayCase{
			item,
			ctMap[item.CaseTypeID],
		}
		displayCases = append(displayCases, &d)
	}

	rc, err := s.GetRegistrationController(w, req, individual)
	if err != nil {
		s.Error(w, err)
		return
	}
	status := rc.Status()
	progressLabel := status.Label
	progress := status.Progress

	// Write Individual attribute values (if any) to Attributes object
	for _, attribute := range attrs.Items {
		values := individual.GetAttribute(attribute.ID)
		attribute.Attributes.Value = values
	}

	// mark cases readonly if needed
	for _, kase := range []*cms.Case{individualAssessment, situationAnalysis} {
		if kase != nil && (kase.Done || status.CurrentStage == -1) {
			kase.Template.MarkAsReadonly()
		}
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "individual", map[string]interface{}{
		"IsNew":                          id == "new",
		"Individual":                     individual,
		"Parties":                        parties,
		"Teams":                          teams,
		"PartyTypes":                     partyTypes,
		"RelationshipTypes":              relationshipTypes,
		"FilteredRelationshipTypes":      filteredRelationshipTypes,
		"Relationships":                  relationshipsForIndividual,
		"Attributes":                     attrs,
		"Cases":                          displayCases,
		"CaseTypes":                      caseTypes,
		"FullNameAttribute":              iam.FullNameAttribute,
		"DisplayNameAttribute":           iam.DisplayNameAttribute,
		"Page":                           "general",
		"Constants":                      s.Constants,
		"IndividualPartyTypeID":          iam.IndividualPartyType.ID,
		"HouseholdPartyTypeID":           iam.HouseholdPartyType.ID,
		"TeamPartyTypeID":                iam.TeamPartyType.ID,
		"IndividualAssessment":           individualAssessment,
		"SituationAnalysis":              situationAnalysis,
		"ProgressLabel":                  progressLabel,
		"Progress":                       progress,
		"IdentificationDocuments":        identificationDocuments,
		"IdentificationDocumentTypes":    identificationDocumentTypes,
		"IdentificationDocumentTypesMap": identificationDocumentTypesMap,
	}); err != nil {
		s.Error(w, err)
		return
	}

}

// zipAttributesAndErrors returns a slice of form.FormElement populated with validated attributes.
func zipAttributesAndErrors(errorList *validation.ErrorList, attributes *iam.AttributeList) []form.FormElement {
	var formElements []form.FormElement
	for _, attribute := range attributes.Items {
		errs := errorList.FindFamily(attribute.Attributes.Name)
		if len(*errs) > 0 && !shouldIgnoreValidationError(attribute) {
			formElements = append(formElements, form.FormElement{
				Type:       attribute.Type,
				Attributes: attribute.Attributes,
				Validation: attribute.Validation,
				Errors:     errs,
			})
		}
	}
	return formElements
}

func shouldIgnoreValidationError(attribute *iam.Attribute) bool {
	// Ignore validation errors on empty optional fields
	if !attribute.Validation.Required && utils.AllEmpty(attribute.Attributes.Value) {
		return true
	}
	return false
}

func PrepRelationshipTypeDropdown(relationshipTypes *iam.RelationshipTypeList) *iam.RelationshipTypeList {

	// TODO:
	// Currently Core only works with individuals
	// this function should be changed when this is no longer
	// the case

	var newList = iam.RelationshipTypeList{}
	for _, relType := range relationshipTypes.Items {

		relTypeOnlyForIndividuals := true

		for _, rule := range relType.Rules {
			if !(rule.PartyTypeRule.FirstPartyTypeID == iam.IndividualPartyType.ID && rule.PartyTypeRule.SecondPartyTypeID == iam.IndividualPartyType.ID) {
				relTypeOnlyForIndividuals = false
			}
		}

		if relTypeOnlyForIndividuals {
			newList.Items = append(newList.Items, relType)
			if relType.IsDirectional {
				newList.Items = append(newList.Items, relType.Mirror())
			}
		}
	}
	return &newList
}

func (s *Server) PostIndividual(ctx context.Context, attrs *iam.AttributeList, id string, w http.ResponseWriter, req *http.Request) (*iam.Individual, error) {

	iamClient, err := s.IAMClient(req)
	if err != nil {
		return nil, err
	}

	if err := req.ParseForm(); err != nil {
		return nil, err
	}

	var individual *iam.Individual
	if len(id) == 0 {
		individual = iam.NewIndividual("")
	} else {
		individual, err = iamClient.Individuals().Get(ctx, id)
	}

	attributeMap := map[string]*iam.Attribute{}
	for _, attribute := range attrs.Items {
		attributeMap[attribute.ID] = attribute
	}

	type RelationshipEntry struct {
		*iam.Relationship
		MarkedForDeletion bool
	}

	//goland:noinspection GoPreferNilSlice
	var rels = []*RelationshipEntry{}

	// Validate Attributes
	attrErrs := validation.ErrorList{}
	f := req.Form
	for key, vals := range f {

		// Populate the Party.attributes
		// as well as the Attributes object (for validation)
		if attr := attrs.FindByName(key); attr != nil {
			individual.Attributes[attr.ID] = vals
			attr.Attributes.Value = vals
			attrErrs = append(attrErrs, iam.ValidateAttribute(attr, validation.NewPath(""))...)
		}

		// Retrieve party relationships
		if strings.HasPrefix(key, "relationships[") {

			keyParts := strings.Split(key, ".")
			if len(keyParts) != 2 {
				err := fmt.Errorf("unexpected form value key: %s", key)
				return nil, err
			}

			if !strings.HasSuffix(keyParts[0], "]") {
				err := fmt.Errorf("unexpected form value key: %s", key)
				return nil, err
			}

			relationshipIdxStr := keyParts[0][14 : len(keyParts[0])-1]
			relationshipIdx, err := strconv.Atoi(relationshipIdxStr)
			if err != nil {
				return nil, err
			}

			var rel *RelationshipEntry
			for {
				if (relationshipIdx + 1) > len(rels) {
					rels = append(rels, &RelationshipEntry{
						Relationship: &iam.Relationship{},
					})
				} else {
					rel = rels[relationshipIdx]
					break
				}
			}

			attrName := keyParts[1]

			switch attrName {
			case "markedForDeletion":
				rel.MarkedForDeletion = vals[0] == "true"
			case "id":
				rel.ID = vals[0]
			case "secondPartyId":
				rel.SecondPartyID = vals[0]
			case "relationshipTypeId":
				rel.RelationshipTypeID = vals[0]
			default:
				err := fmt.Errorf("unexpected relationship attribute: %s", attrName)
				return nil, err
			}
		}
	}

	// Verify attribute validation and act accordingly
	if len(attrErrs) > 0 {
		status := attrErrs.Status(http.StatusUnprocessableEntity, "invalid case")
		return nil, &status
	}

	// Update or create the individual
	var storageAction string
	if id == "" {
		individual.PartyTypeIDs = append(individual.PartyTypeIDs, iam.BeneficiaryPartyType.ID)
		individual, err = iamClient.Individuals().Create(ctx, individual)
		if err != nil {
			return nil, err
		}
		err = s.createDefaultIndividualIntakeCases(req, individual)
		if err != nil {
			return nil, err
		}
		storageAction = "created"
	} else {
		if individual, err = iamClient.Individuals().Update(ctx, individual); err != nil {
			return nil, err
		}
		storageAction = "updated"
	}

	// Update, create or delete the relationships

	for _, rel := range rels {

		if len(rel.ID) == 0 {
			// Create the relationship
			relationship := rel.Relationship
			relationship.FirstPartyID = individual.ID
			if _, err := iamClient.Relationships().Create(ctx, rel.Relationship); err != nil {
				return nil, err
			}
		} else if !rel.MarkedForDeletion {
			// Update the relationship
			relationship := rel.Relationship
			relationship.FirstPartyID = individual.ID
			if _, err := iamClient.Relationships().Update(ctx, rel.Relationship); err != nil {
				return nil, err
			}
		} else {
			// Delete the relationship
			if err := iamClient.Relationships().Delete(ctx, rel.Relationship.ID); err != nil {
				return nil, err
			}
		}

	}

	// Set flash notification
	if err := s.sessionManager.AddNotification(req, w, &sessionmanager.Notification{
		Message: fmt.Sprintf("Individual \"%s\" successfully %s", individual.String(), storageAction),
		Theme:   "success",
	}); err != nil {
		return nil, err
	}

	if storageAction == "created" {
		w.Header().Set("Location", "/individuals/"+individual.ID)
		w.WriteHeader(http.StatusSeeOther)
	}
	return individual, nil
}

func (s *Server) createDefaultIndividualIntakeCases(req *http.Request, individual *iam.Individual) error {
	var situationAnalysisCaseType = &seeder.UGSituationalAnalysisCaseType
	var individualResponseCaseType = &seeder.UGIndividualResponseCaseType

	cmsClient, err := s.CMSClient(req)
	if err != nil {
		return err
	}
	creatorId := ""
	ctx := req.Context()
	subject := ctx.Value("Subject")
	if subject != nil {
		creatorId = subject.(string)
	}
	// Create UG Intake Cases for new individual
	if situationAnalysisCaseType != nil {
		_, err = cmsClient.Cases().Create(ctx, &cms.Case{
			CaseTypeID: situationAnalysisCaseType.ID,
			PartyID:    individual.ID,
			Done:       false,
			TeamID:     situationAnalysisCaseType.TeamID,
			CreatorID:  creatorId,
			IntakeCase: situationAnalysisCaseType.IntakeCaseType,
		})
		if err != nil {
			return err
		}
	}
	if individualResponseCaseType != nil {
		_, err = cmsClient.Cases().Create(ctx, &cms.Case{
			CaseTypeID: individualResponseCaseType.ID,
			PartyID:    individual.ID,
			Done:       false,
			TeamID:     individualResponseCaseType.TeamID,
			CreatorID:  creatorId,
			IntakeCase: situationAnalysisCaseType.IntakeCaseType,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
