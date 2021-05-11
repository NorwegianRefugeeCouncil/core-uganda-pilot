package testing

import (
	"context"
	"encoding/json"
	v12 "github.com/nrc-no/core/apps/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/watch"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"time"
)

func (s *MainTestSuite) TestFormDefinitionCRUD() {
	t := s.T()

	watchCtx, watchCancel := context.WithCancel(s.ctx)
	defer watchCancel()

	var events []watch.Event

	watcher, err := s.nrcClient.FormDefinitions().Watch(watchCtx, v1.ListOptions{})
	if err != nil {
		s.T().Errorf("cannot start watch: %v", err)
		return
	}

	go func() {
		for event := range watcher.ResultChan() {
			s.T().Logf("event: %#v", event)
			events = append(events, event)
			if len(events) == 2 {
				watchCancel()
			}
		}
	}()

	var formDefinition = v12.FormDefinition{
		TypeMeta: v1.TypeMeta{
			Kind:       "FormDefinition",
			APIVersion: "core/v1",
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

	time.Sleep(30 * time.Second)

	// Asserting equality of input & output
	assert.Equal(t, "core/v1", out.APIVersion)
	assert.Equal(t, "FormDefinition", out.Kind)
	assert.Equal(t, "core.nrc.no", out.Spec.Group)
	assert.Equal(t, "customresources", out.Spec.Names.Plural)
	assert.Equal(t, "customresource", out.Spec.Names.Singular)
	assert.Equal(t, "CustomResource", out.Spec.Names.Kind)
	if !assert.Equal(t, 1, len(out.Spec.Versions)) {
		return
	}
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

	select {
	case <-watchCtx.Done():
	}

	if assert.Len(t, events, 2) {
		assert.Equal(t, "insert", events[0].Type)
		assert.Equal(t, "replace", events[1].Type)
	}

	list, err := s.nrcClient.FormDefinitions().List(s.ctx)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, list) {
		return
	}
	if !assert.NotNil(t, list.Items) {
		return
	}
	bytes, err := json.MarshalIndent(list, "", " ")
	if err == nil {
		t.Log(string(bytes))
	}

	assert.GreaterOrEqual(t, len(list.Items), 1)

}
