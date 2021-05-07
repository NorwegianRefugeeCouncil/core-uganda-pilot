package formdefinitions

import (
	"encoding/json"
	"github.com/nrc-no/core/apps/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/apps/api/pkg/apis/meta"
	"github.com/nrc-no/core/apps/api/pkg/util/exceptions"
	"github.com/nrc-no/core/apps/api/pkg/util/validation/field"
	"io/ioutil"
	"net/http"
	"time"
)

// Post formDefinition
func (h *Handler) Post(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.scope.Error(err, w, req)
		return
	}

	var formDefinition v1.FormDefinition
	if err := json.Unmarshal(bodyBytes, &formDefinition); err != nil {
		h.scope.Error(err, w, req)
		return
	}

	obj, err := meta.Accessor(&formDefinition)
	if err != nil {
		h.scope.Error(err, w, req)
		return
	}

	var fieldErrors field.ErrorList
	if formDefinition.TypeMeta.Kind == "" {
		fieldErrors = append(fieldErrors, field.Required(field.NewPath("kind"), "kind is required"))
	}
	if formDefinition.TypeMeta.APIVersion == "" {
		fieldErrors = append(fieldErrors, field.Required(field.NewPath("apiVersion"), "apiVersion is required"))
	}
	if formDefinition.Spec.Group == "" {
		fieldErrors = append(fieldErrors, field.Required(field.NewPath("spec.group"), "group is required"))
	}
	if formDefinition.Spec.Names.Singular == "" {
		fieldErrors = append(fieldErrors, field.Required(field.NewPath("spec.names.singular"), "singular name is required"))
	}
	if formDefinition.Spec.Names.Plural == "" {
		fieldErrors = append(fieldErrors, field.Required(field.NewPath("spec.names.plural"), "plural name is required"))
	}
	if formDefinition.Spec.Names.Kind == "" {
		fieldErrors = append(fieldErrors, field.Required(field.NewPath("spec.names.kind"), "kind name is required"))
	}
	if len(formDefinition.Spec.Versions) == 0 {
		fieldErrors = append(fieldErrors, field.Required(field.NewPath("spec.versions"), "versions must not be empty"))
	}
	versionsPath := field.NewPath("spec.versions")
	for i, version := range formDefinition.Spec.Versions {
		versionPath := versionsPath.Index(i)
		if len(version.Name) == 0 {
			fieldErrors = append(fieldErrors, field.Required(versionPath.Child("name"), "version name is required"))
		}
		root := version.Schema.FormSchema.Root
		rootPath := versionPath.Child("schema", "formSchema", "root")
		if root.Type != "section" {
			if root.Key == "" {
				fieldErrors = append(fieldErrors, field.Required(rootPath.Child("key"), "root key is required when root is not section"))
			}
		}
	}

	if len(fieldErrors) > 0 {
		err = exceptions.NewInvalid(h.scope.Kind.GroupKind(), "", fieldErrors)
		h.scope.Error(err, w, req)
		return
	}

	obj.SetCreationTimestamp(time.Now().UTC())
	obj.SetDeletionTimestamp(nil)
	formDefinition.SetAPIVersion(h.scope.Kind.GroupVersion().String())

	var out v1.FormDefinition
	if err := h.storage.Create(ctx, &formDefinition, &out); err != nil {
		h.scope.Error(err, w, req)
		return
	}

	responseBytes, err := json.Marshal(&out)
	if err != nil {
		h.scope.Error(err, w, req)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)

}
