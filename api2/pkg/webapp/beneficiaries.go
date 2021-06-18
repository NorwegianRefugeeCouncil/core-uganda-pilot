package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/cases/cases"
	"github.com/nrc-no/core-kafka/pkg/cases/casetypes"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/parties/beneficiaries"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"github.com/nrc-no/core-kafka/pkg/parties/relationships"
	"github.com/nrc-no/core-kafka/pkg/parties/relationshiptypes"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) Beneficiaries(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	attributesClient := attributes.NewClient("http://localhost:9000")
	attrs, err := attributesClient.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostBeneficiary(ctx, attrs.Items, "", w, req)
		return
	}

	list, err := h.beneficiaryClient.List(ctx, beneficiaries.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.template.ExecuteTemplate(w, "beneficiaries", map[string]interface{}{
		"Beneficiaries": list,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) Beneficiary(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var b *beneficiaries.Beneficiary
	var bList *beneficiaries.BeneficiaryList
	var ctList *casetypes.CaseTypeList
	var cList *cases.CaseList
	var relationshipsForBeneficiary *relationships.RelationshipList
	var relationshipTypes *relationshiptypes.RelationshipTypeList
	var attrs *attributes.AttributeList

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if id == "new" {
			b = beneficiaries.NewBeneficiary("")
			return nil
		}
		var err error
		b, err = h.beneficiaryClient.Get(waitCtx, id)
		return err
	})

	g.Go(func() error {
		var err error
		bList, err = h.beneficiaryClient.List(waitCtx, beneficiaries.ListOptions{})
		return err
	})

	g.Go(func() error {
		if id == "new" {
			relationshipsForBeneficiary = &relationships.RelationshipList{
				Items: []*relationships.Relationship{},
			}
			return nil
		}
		var err error
		relationshipsForBeneficiary, err = h.relationshipClient.List(waitCtx, relationships.ListOptions{EitherParty: id})
		return err
	})

	g.Go(func() error {
		var err error
		relationshipTypes, err = h.relationshipTypeClient.List(waitCtx, relationshiptypes.ListOptions{PartyType: partytypes.IndividualPartyType.ID})
		return err
	})

	g.Go(func() error {
		var err error
		attrs, err = h.attributeClient.List(waitCtx)
		return err
	})

	g.Go(func() error {
		var err error
		ctList, err = h.caseTypeClient.List(ctx, casetypes.ListOptions{})
		return err
	})

	g.Go(func() error {
		var err error
		cList, err = h.caseClient.List(ctx, cases.ListOptions{PartyID: id})
		return err
	})

	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	relationshipTypes = PrepRelationshipTypeDropdown(relationshipTypes)

	if req.Method == "POST" {
		h.PostBeneficiary(ctx, attrs.Items, id, w, req)
		return
	}

	type DisplayCase struct {
		Case     *cases.Case
		CaseType *casetypes.CaseType
	}

	ctMap := map[string]*casetypes.CaseType{}
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

	if err := h.template.ExecuteTemplate(w, "beneficiary", map[string]interface{}{
		"IsNew":             id == "new",
		"Beneficiary":       b,
		"Parties":           bList,
		"RelationshipTypes": relationshipTypes,
		"Relationships":     relationshipsForBeneficiary,
		"Attributes":        attrs,
		"Cases":             displayCases,
		"CaseTypes":         ctList,
		"BasePath": h.relationshipPartiesClient.BasePath,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func PrepRelationshipTypeDropdown(relationshipTypes *relationshiptypes.RelationshipTypeList) *relationshiptypes.RelationshipTypeList {
	var newList = relationshiptypes.RelationshipTypeList{}
	for _, relType := range relationshipTypes.Items {
		if relType.IsDirectional {
			for _, rule := range relType.Rules {
				if rule.FirstPartyType == partytypes.IndividualPartyType.ID {
					newList.Items = append(newList.Items, relType)
				}
				if rule.SecondPartyType == partytypes.IndividualPartyType.ID {
					newList.Items = append(newList.Items, relType.Reversed())
				}
			}
		} else {
			newList.Items = append(newList.Items, relType)
		}
	}
	return &newList
}

func (h *Handler) PostBeneficiary(
	ctx context.Context,
	attrs []*attributes.Attribute,
	id string,
	w http.ResponseWriter,
	req *http.Request) {

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b := beneficiaries.NewBeneficiary(id)

	attributeMap := map[string]*attributes.Attribute{}
	for _, attribute := range attrs {
		attributeMap[attribute.ID] = attribute
	}

	type RelationshipEntry struct {
		*relationships.Relationship
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
						Relationship: &relationships.Relationship{},
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
				rel.SecondParty = vals[0]
			case "relationshipTypeId":
				rel.RelationshipTypeID = vals[0]
			default:
				err := fmt.Errorf("unexpected relationship attribute: %s", attrName)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	var beneficiary *beneficiaries.Beneficiary

	// Update or create the beneficiary
	if id == "" {
		var err error
		beneficiary, err = h.beneficiaryClient.Create(ctx, b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		var err error
		if beneficiary, err = h.beneficiaryClient.Update(ctx, b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Update, create or delete the relationships

	for _, rel := range rels {

		if len(rel.ID) == 0 {
			// Create the relationship
			relationship := rel.Relationship
			relationship.FirstParty = beneficiary.ID
			if _, err := h.relationshipClient.Create(ctx, rel.Relationship); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if !rel.MarkedForDeletion {
			// Update the relationship
			relationship := rel.Relationship
			relationship.FirstParty = beneficiary.ID
			if _, err := h.relationshipClient.Update(ctx, rel.Relationship); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			// Delete the relationship
			if err := h.relationshipClient.Delete(ctx, rel.Relationship.ID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	}

	w.Header().Set("Location", "/beneficiaries/"+beneficiary.ID)
	w.WriteHeader(http.StatusSeeOther)

}
