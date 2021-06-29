package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/sessionmanager"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"strings"
)

func (h *Server) Individuals(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	iamCli := h.IAMClient(ctx)
	attrs, err := iamCli.Attributes().List(ctx, iam.AttributeListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostIndividual(ctx, attrs.Items, "", w, req)
		return
	}

	list, err := iamCli.Individuals().List(ctx, iam.IndividualListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "individuals", map[string]interface{}{
		"Individuals": list,
		"Page":        "list",
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Server) IndividualCredentials(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	iamClient := h.IAMClient(ctx)

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	individual, err := iamClient.Individuals().Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostIndividualCredentials(w, req, individual.ID)
		return
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "individual_credentials", map[string]interface{}{
		"Page":       "credentials",
		"Individual": individual,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Server) PostIndividualCredentials(w http.ResponseWriter, req *http.Request, partyID string) {
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	values := req.Form
	password := values.Get("password")

	if err := h.login.Login().SetCredentials(ctx, partyID, password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/individuals/"+partyID+"/credentials", http.StatusSeeOther)
}

func (h *Server) Individual(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	cmsClient := h.CMSClient(ctx)
	iamClient := h.IAMClient(ctx)

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var b *iam.Individual
	var bList *iam.IndividualList
	var ctList *cms.CaseTypeList
	var cList *cms.CaseList
	var partyTypes *iam.PartyTypeList
	var relationshipsForIndividual *iam.RelationshipList
	var relationshipTypes *iam.RelationshipTypeList
	var attrs *iam.AttributeList

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if id == "new" {
			b = iam.NewIndividual("")
			return nil
		}
		var err error
		b, err = iamClient.Individuals().Get(waitCtx, id)
		return err
	})

	g.Go(func() error {
		var err error
		bList, err = iamClient.Individuals().List(waitCtx, iam.IndividualListOptions{})
		return err
	})

	g.Go(func() error {
		var err error
		partyTypes, err = iamClient.PartyTypes().List(waitCtx, iam.PartyTypeListOptions{})
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
		ctList, err = cmsClient.CaseTypes().List(ctx, cms.CaseTypeListOptions{})
		return err
	})

	g.Go(func() error {
		var err error
		cList, err = cmsClient.Cases().List(ctx, cms.CaseListOptions{PartyID: id})
		return err
	})

	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	relationshipTypes = PrepRelationshipTypeDropdown(relationshipTypes)

	if req.Method == "POST" {
		h.PostIndividual(ctx, attrs.Items, id, w, req)
		return
	}

	type DisplayCase struct {
		Case     *cms.Case
		CaseType *cms.CaseType
	}

	ctMap := map[string]*cms.CaseType{}
	for _, item := range ctList.Items {
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

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "individual", map[string]interface{}{
		"IsNew":              id == "new",
		"Individual":         b,
		"Parties":            bList,
		"PartyTypes":         partyTypes,
		"RelationshipTypes":  relationshipTypes,
		"Relationships":      relationshipsForIndividual,
		"Attributes":         attrs,
		"Cases":              displayCases,
		"CaseTypes":          ctList,
		"FirstNameAttribute": iam.FirstNameAttribute,
		"LastNameAttribute":  iam.LastNameAttribute,
		"Page":               "general",
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func PrepRelationshipTypeDropdown(relationshipTypes *iam.RelationshipTypeList) *iam.RelationshipTypeList {
	var newList = iam.RelationshipTypeList{}
	for _, relType := range relationshipTypes.Items {
		if relType.IsDirectional {
			for _, rule := range relType.Rules {
				if rule.PartyTypeRule.FirstPartyTypeID == iam.IndividualPartyType.ID {
					newList.Items = append(newList.Items, relType)
				}
				if rule.PartyTypeRule.SecondPartyTypeID == iam.IndividualPartyType.ID {
					newList.Items = append(newList.Items, relType.Mirror())
				}
			}
		} else {
			newList.Items = append(newList.Items, relType)
		}
	}
	return &newList
}

func (h *Server) PostIndividual(
	ctx context.Context,
	attrs []*iam.Attribute,
	id string,
	w http.ResponseWriter,
	req *http.Request) {

	iamClient := h.IAMClient(ctx)

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b := iam.NewIndividual(id)

	attributeMap := map[string]*iam.Attribute{}
	for _, attribute := range attrs {
		attributeMap[attribute.ID] = attribute
	}

	type RelationshipEntry struct {
		*iam.Relationship
		MarkedForDeletion bool
	}

	//goland:noinspection GoPreferNilSlice
	var rels = []*RelationshipEntry{}

	f := req.Form
	for key, vals := range f {

		// Populate the Party.attributes
		if strings.HasPrefix(key, "attribute[") {

			if !strings.HasSuffix(key, "]") {
				err := fmt.Errorf("unexpected form value key: %s", key)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			attrId, err := uuid.FromString(key[10 : len(key)-1])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			attr, ok := attributeMap[attrId.String()]
			if !ok {
				err := fmt.Errorf("attribute with id %s not found", attrId)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			b.Attributes[attr.ID] = vals

		}

		// Retrieve party relationships
		if strings.HasPrefix(key, "relationships[") {

			keyParts := strings.Split(key, ".")
			if len(keyParts) != 2 {
				err := fmt.Errorf("unexpected form value key: %s", key)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !strings.HasSuffix(keyParts[0], "]") {
				err := fmt.Errorf("unexpected form value key: %s", key)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			relationshipIdxStr := keyParts[0][14 : len(keyParts[0])-1]
			relationshipIdx, err := strconv.Atoi(relationshipIdxStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
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
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	var individual *iam.Individual

	// Update or create the individual
	if id == "" {
		var err error
		individual, err = iamClient.Individuals().Create(ctx, b)
		if err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.sessionManager.AddNotification(ctx, &sessionmanager.Notification{
			Message: fmt.Sprintf("%s successfully created", b.String()),
			Theme:   "success",
		})
	} else {
		var err error
		if individual, err = iamClient.Individuals().Update(ctx, b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.sessionManager.AddNotification(ctx, &sessionmanager.Notification{
			Message: fmt.Sprintf("%s successfully updated", b.String()),
			Theme:   "success",
		})
	}

	// Update, create or delete the relationships

	for _, rel := range rels {

		if len(rel.ID) == 0 {
			// Create the relationship
			relationship := rel.Relationship
			relationship.FirstPartyID = individual.ID
			if _, err := iamClient.Relationships().Create(ctx, rel.Relationship); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if !rel.MarkedForDeletion {
			// Update the relationship
			relationship := rel.Relationship
			relationship.FirstPartyID = individual.ID
			if _, err := iamClient.Relationships().Update(ctx, rel.Relationship); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			// Delete the relationship
			if err := iamClient.Relationships().Delete(ctx, rel.Relationship.ID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	}

	w.Header().Set("Location", "/individuals/"+individual.ID)
	w.WriteHeader(http.StatusSeeOther)

}
