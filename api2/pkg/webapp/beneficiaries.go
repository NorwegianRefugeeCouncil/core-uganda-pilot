package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/subjects/api"
	"github.com/nrc-no/core-kafka/pkg/subjects/attributes"
	"github.com/nrc-no/core-kafka/pkg/subjects/beneficiaries"
	"github.com/nrc-no/core-kafka/pkg/subjects/relationships"
	"github.com/nrc-no/core-kafka/pkg/subjects/relationshiptypes"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) Beneficiaries(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	cli := beneficiaries.NewClient("http://localhost:9000")
	list, err := cli.List(ctx)
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

	beneficiaryClient := beneficiaries.NewClient("http://localhost:9000")
	relationshipTypeClient := relationshiptypes.NewClient("http://localhost:9000")
	relationshipClient := relationships.NewClient("http://localhost:9000")

	var b *api.Beneficiary
	var bList *api.BeneficiaryList
	var relationshipsForBeneficiary *api.RelationshipList
	var relationshipTypes *api.RelationshipTypeList

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if id == "new" {
			b = &api.Beneficiary{}
			return nil
		}
		var err error
		b, err = beneficiaryClient.Get(waitCtx, id)
		return err
	})

	g.Go(func() error {
		var err error
		bList, err = beneficiaryClient.List(waitCtx)
		return err
	})

	g.Go(func() error {
		if id == "new" {
			relationshipsForBeneficiary = &api.RelationshipList{
				Items: []*api.Relationship{},
			}
			return nil
		}
		var err error
		relationshipsForBeneficiary, err = relationshipClient.List(waitCtx, relationships.ListOptions{Party: id})
		return err
	})

	g.Go(func() error {
		var err error
		relationshipTypes, err = relationshipTypeClient.List(waitCtx)
		return err
	})

	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostBeneficiary(ctx, beneficiaryClient, id, w, req)
		return
	}

	attributesClient := attributes.NewClient("http://localhost:9000")
	attrs, err := attributesClient.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.template.ExecuteTemplate(w, "beneficiary", map[string]interface{}{
		"Beneficiary":       b,
		"Parties":           bList,
		"RelationshipTypes": relationshipTypes,
		"Relationships":     relationshipsForBeneficiary,
		"Attributes":        attrs,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) PostBeneficiary(ctx context.Context, cli *beneficiaries.Client, id string, w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = &api.Beneficiary{
		ID:         "",
		Attributes: map[string]*api.AttributeValue{},
	}

	var attrs []*api.AttributeValue

	f := req.Form
	for key, vals := range f {

		// Making sure the key starts with "attribute[
		if !strings.HasPrefix(key, "attribute[") {
			err := fmt.Errorf("unexpected form value key: %s", key)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// splitting the attribute key
		keyParts := strings.Split(key, ".")

		// making sure the first part of the key ends with ]
		// like in attribute[{index}]
		if !strings.HasSuffix(keyParts[0], "]") {
			err := fmt.Errorf("unexpected form value key: %s", key)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Retrieve the index of the attribute
		attributeIndexStr := key[10 : len(keyParts[0])-1]

		// Convert to number
		attributeIdx, err := strconv.Atoi(attributeIndexStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Append an element to the attribute slice if needed
		if attrs == nil {
			attrs = []*api.AttributeValue{}
		}
		for len(attrs)-1 < attributeIdx {
			attrs = append(attrs, &api.AttributeValue{})
		}

		// Retrieve the attribute from the slice
		attr := attrs[attributeIdx]

		if len(keyParts) == 2 {

			// Happens when we are dealing with first level attributes
			// Key in the format attribute[0].translations[0].locale ..

			if keyParts[1] == "name" {
				attr.Name = vals[0]
			} else if keyParts[1] == "valueType" {
				// TODO: attr.Va = vals[0]
			} else if keyParts[1] == "subjectType" {
				attr.SubjectType = api.SubjectType(vals[0])
			} else if keyParts[1] == "isPersonallyIdentifiableInfo" {
				attr.IsPersonallyIdentifiableInfo = vals[0] == "true"
			} else if keyParts[1] == "value" {
				attr.Value = vals[0]
			} else if keyParts[1] == "id" {
				attr.ID = vals[0]
			}

		} else if len(keyParts) == 3 {

			// Happens when we are dealing with translations
			// Key in the format attribute[0].translations[0].locale ..

			if !strings.HasPrefix(keyParts[1], "translations[") {
				err := fmt.Errorf("unexpected form value key: %s", key)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !strings.HasSuffix(keyParts[1], "]") {
				err := fmt.Errorf("unexpected form value key: %s", key)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			translationIdxStr := keyParts[1][13 : len(keyParts[1])-1]
			translationIdx, err := strconv.Atoi(translationIdxStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if attr.Translations == nil {
				attr.Translations = []api.AttributeTranslation{}
			}
			for len(attr.Translations)-1 < translationIdx {
				attr.Translations = append(attr.Translations, api.AttributeTranslation{})
			}
			translation := &attr.Translations[translationIdx]

			if keyParts[2] == "locale" {
				translation.Locale = vals[0]
			} else if keyParts[2] == "short" {
				translation.ShortFormulation = vals[0]
			} else if keyParts[2] == "long" {
				translation.LongFormulation = vals[0]
			}

		}

	}

	beneficiaryAttributes := map[string]*api.AttributeValue{}
	for _, value := range attrs {
		beneficiaryAttributes[value.Name] = value
	}

	_, err := cli.Update(ctx, &api.Beneficiary{
		ID:         id,
		Attributes: beneficiaryAttributes,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/beneficiaries/"+id)
	w.WriteHeader(http.StatusSeeOther)

}
