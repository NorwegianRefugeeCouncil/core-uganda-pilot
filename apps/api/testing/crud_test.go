package testing

import (
	"context"
	v12 "github.com/nrc-no/core/apps/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"time"
)

func (s *MainTestSuite) TestFormDefinitionCRUD() {
	t := s.T()

	watchCtx, watchCancel := context.WithCancel(s.ctx)
	defer watchCancel()

	watch, err := s.nrcClient.FormDefinitions().Watch(watchCtx)
	if err != nil {
		s.T().Errorf("cannot start watch: %v", err)
		return
	}

	go func() {
		for event := range watch.ResultChan() {
			s.T().Logf("%#v", event)
			watchCancel()
		}
	}()

	var formDefinition = v12.FormDefinition{
		TypeMeta: v1.TypeMeta{
			Kind:       "FormDefinition",
			APIVersion: "core.nrc.no/v1",
		},
		Spec: v12.FormDefinitionSpec{
			Group: "core.nrc.no",
			Names: v12.CustomResourceNames{
				Plural:   "customresources",
				Singular: "customresource",
				Kind:     "CustomResource",
			},
			Versions: []v12.FormDefinitionVersion{
				{
					Name: "v1",
					Schema: v12.FormSchema{
						FormSchema: v12.FormSchemaDefinition{
							Root: v12.FormElement{
								Type: "text",
								ID:   uuid.NewV4().String(),
								Key:  "key",
								Description: []v12.TranslatedString{
									{
										Locale: "en",
										Value:  "Description",
									},
								},
								Name: []v12.TranslatedString{
									{
										Locale: "en",
										Value:  "Name",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	out, err := s.nrcClient.FormDefinitions().Create(s.ctx, &formDefinition)
	if err != nil {
		s.T().Errorf("could not create form definition: %v", err)
		return
	}

	// Asserting equality of input & output
	assert.Equal(t, "core.nrc.no/v1", out.APIVersion)
	assert.Equal(t, "FormDefinition", out.Kind)
	assert.Equal(t, "core.nrc.no", out.Spec.Group)
	assert.Equal(t, "customresources", out.Spec.Names.Plural)
	assert.Equal(t, "customresource", out.Spec.Names.Singular)
	assert.Equal(t, "CustomResource", out.Spec.Names.Kind)
	assert.Equal(t, 1, len(out.Spec.Versions))
	assert.Equal(t, "v1", out.Spec.Versions[0].Name)
	assert.Equal(t, "key", out.Spec.Versions[0].Schema.FormSchema.Root.Key)
	assert.NotEmpty(t, out.Spec.Versions[0].Schema.FormSchema.Root.ID)
	assert.Equal(t, 1, len(out.Spec.Versions[0].Schema.FormSchema.Root.Name))
	assert.Equal(t, "en", out.Spec.Versions[0].Schema.FormSchema.Root.Name[0].Locale)
	assert.Equal(t, "Name", out.Spec.Versions[0].Schema.FormSchema.Root.Name[0].Value)
	assert.Equal(t, "text", out.Spec.Versions[0].Schema.FormSchema.Root.Type)
	assert.Equal(t, 1, len(out.Spec.Versions[0].Schema.FormSchema.Root.Description))
	assert.Equal(t, "en", out.Spec.Versions[0].Schema.FormSchema.Root.Description[0].Locale)
	assert.Equal(t, "Description", out.Spec.Versions[0].Schema.FormSchema.Root.Description[0].Value)
	assert.NotEmpty(t, out.ObjectMeta.UID)
	assert.NotEqual(t, time.Time{}, out.ObjectMeta.CreationTimestamp)
	assert.Nil(t, out.ObjectMeta.DeletionTimestamp)
	assert.Equal(t, 1, out.ObjectMeta.ResourceVersion)

	// Update form definition
	out.Spec.Names.Plural = "abc"
	updated, err := s.nrcClient.FormDefinitions().Update(s.ctx, out)
	if err != nil {
		t.Errorf("unable to update form definition: %v", err)
		return
	}
	assert.Equal(t, "abc", updated.Spec.Names.Plural)

	// Should update version
	assert.Equal(t, 2, updated.ResourceVersion)

	t.Logf("\n%v", updated)

	time.Sleep(2 * time.Second)
}
