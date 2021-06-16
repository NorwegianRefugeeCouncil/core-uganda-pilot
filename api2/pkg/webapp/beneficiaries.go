package webapp

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/parties/beneficiaries"
	"github.com/nrc-no/core-kafka/pkg/parties/relationships"
	"github.com/nrc-no/core-kafka/pkg/parties/relationshiptypes"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"
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
		relationshipsForBeneficiary, err = h.relationshipClient.List(waitCtx, relationships.ListOptions{Party: id})
		return err
	})

	g.Go(func() error {
		var err error
		relationshipTypes, err = h.relationshipTypeClient.List(waitCtx, relationshipTypes.ListOptions{PartyType: b.PartyTypes[0]})
		return err
	})

	g.Go(func() error {
		var err error
		attrs, err = h.attributeClient.List(waitCtx)
		return err
	})

	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostBeneficiary(ctx, attrs.Items, id, w, req)
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

	f := req.Form
	for key, vals := range f {

		// Making sure the key starts with "attribute[
		if !strings.HasPrefix(key, "attribute[") {
			err := fmt.Errorf("unexpected form value key: %s", key)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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

	if id == "" {
		newBenef, err := h.beneficiaryClient.Create(ctx, b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/beneficiaries/"+newBenef.ID)
		w.WriteHeader(http.StatusSeeOther)
	} else {
		if _, err := h.beneficiaryClient.Update(ctx, b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", "/beneficiaries/"+id)
		w.WriteHeader(http.StatusSeeOther)
	}
}
