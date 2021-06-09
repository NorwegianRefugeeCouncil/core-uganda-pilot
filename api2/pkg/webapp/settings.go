package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/expressions"
	"github.com/nrc-no/core-kafka/pkg/subjects/api"
	"github.com/nrc-no/core-kafka/pkg/subjects/attributes"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strings"
)

func (h *Handler) Settings(w http.ResponseWriter, req *http.Request) {
	if err := h.template.ExecuteTemplate(w, "settings", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Attributes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	cli := attributes.NewClient("http://localhost:9000")

	if req.Method == "POST" {
		h.PostAttribute(ctx, cli, &api.Attribute{}, w, req)
		return
	}

	list, err := cli.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.template.ExecuteTemplate(w, "attributes", map[string]interface{}{
		"Attributes": list,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) NewAttribute(w http.ResponseWriter, req *http.Request) {
	if err := h.template.ExecuteTemplate(w, "attribute", map[string]interface{}{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Attribute(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	cli := attributes.NewClient("http://localhost:9000")

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("No id in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a, err := cli.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostAttribute(ctx, cli, a, w, req)
	}

	if err := h.template.ExecuteTemplate(w, "attribute", map[string]interface{}{
		"Attribute": a,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) PostAttribute(ctx context.Context, cli *attributes.Client, attribute *api.Attribute, w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	values := req.Form

	translationMap := map[string]*api.AttributeTranslation{}

	for key, v := range values {
		if !strings.HasPrefix(key, "translations.") {
			continue
		}
		parts := strings.Split(key, ".")
		if len(parts) != 3 {
			http.Error(w, "unexpected translation key. Expected 'translation.{locale}.{short/long}' format", http.StatusInternalServerError)
			return
		}

		locale := parts[1]
		part := parts[2]

		if part != "long" && part != "short" {
			http.Error(w, "unexpected translation key. Expected 'translation.{locale}.{short/long}' format", http.StatusInternalServerError)
			return
		}

		if _, ok := translationMap[locale]; !ok {
			translationMap[locale] = &api.AttributeTranslation{
				Locale: locale,
			}
		}
		t := translationMap[locale]

		if part == "long" {
			t.LongFormulation = v[0]
		} else if part == "short" {
			t.ShortFormulation = v[0]
		} else {
			http.Error(w, "unexpected translation key. Expected 'translation.{locale}.{short/long}' format", http.StatusInternalServerError)
			return
		}

	}

	var translations []api.AttributeTranslation
	for _, translation := range translationMap {
		translations = append(translations, *translation)
	}

	isNew := false
	if len(attribute.ID) == 0 {
		attribute.ID = uuid.NewV4().String()
		isNew = true
	}

	attribute.Name = values.Get("name")
	attribute.ValueType = expressions.ValueType{}
	attribute.SubjectType = api.SubjectType(values.Get("subjectType"))
	attribute.Translations = translations
	attribute.IsPersonallyIdentifiableInfo = values.Get("isPersonallyIdentifiableInfo") == "true"

	var out *api.Attribute

	if isNew {
		var err error
		out, err = cli.Create(ctx, attribute)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		var err error
		out, err = cli.Update(ctx, attribute)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Location", "/settings/attributes/"+out.ID)
	w.WriteHeader(http.StatusSeeOther)

}
