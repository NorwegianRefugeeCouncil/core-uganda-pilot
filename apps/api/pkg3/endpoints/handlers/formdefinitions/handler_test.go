package formdefinitions

import (
	"context"
	goes "github.com/EventStore/EventStore-Client-Go/client"
	"github.com/nrc-no/core/apps/api/pkg3/apis"
	"github.com/nrc-no/core/apps/api/pkg3/client/nrc"
	"github.com/nrc-no/core/apps/api/pkg3/client/rest"
	"github.com/nrc-no/core/apps/api/pkg3/server"
	"github.com/nrc-no/core/apps/api/pkg3/storage/eventstoredb"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateFormDefinition(t *testing.T) {
	ctx := context.TODO()

	// Create API server
	apiServer := server.NewServer()

	// Create HTTP server
	httpServer := httptest.NewServer(apiServer)
	defer httpServer.Close()

	// Create client
	nrcClient, err := nrc.NewForConfig(&rest.Config{
		ContentConfig: rest.DefaultContentConfig,
		Host:          httpServer.URL,
	})
	if err != nil {
		t.Errorf("unable to create rest client: %v", err)
		return
	}

	// Create eventdb client
	eventStoreDBClient, err := goes.NewClient(&goes.Configuration{
		Address:    "localhost:2113",
		DisableTLS: true,
	})
	if err != nil {
		t.Errorf("failed to create eventstoredb client: %v", err)
		return
	}
	if err := eventStoreDBClient.Connect(); err != nil {
		t.Errorf("failed to connect to evenstore: %v", err)
		return
	}
	defer eventStoreDBClient.Close()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(options.Credential{
		Username: "root",
		Password: "example",
	}))
	if err != nil {
		t.Errorf("could not connect to mongo: %v", err)
		return
	}
	defer mongoClient.Disconnect(ctx)

	// Create storage
	store := eventstoredb.NewStore(eventStoreDBClient, mongoClient)

	go func() {
		if err := store.Watch(context.TODO()); err != nil {
			t.Errorf("failed to start watch: %v", err)
			return
		}
	}()

	time.Sleep(1 * time.Second)

	// Install FormDefinitions api
	Install(apiServer.Container, store)

	var formDefinition = apis.FormDefinition{
		TypeMeta: apis.TypeMeta{
			Kind:       "FormDefinition",
			APIVersion: apiVersion,
		},
		Spec: apis.FormDefinitionSpec{
			Group: groupName,
			Names: apis.CustomResourceNames{
				Plural:   "customresources",
				Singular: "customresource",
				Kind:     "CustomResource",
			},
			Versions: []apis.FormDefinitionVersion{
				{
					Name: "v1",
					Schema: apis.FormSchema{
						FormSchema: apis.FormSchemaDefinition{
							Root: apis.FormElement{
								Type: "text",
								ID:   uuid.NewV4().String(),
								Key:  "key",
								Description: []apis.TranslatedString{
									{
										Locale: "en",
										Value:  "Description",
									},
								},
								Name: []apis.TranslatedString{
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

	out, err := nrcClient.FormDefinitions().Create(ctx, &formDefinition)
	if err != nil {
		t.Errorf("could not create form definition: %v", err)
		return
	}

	assert.Equal(t, "nrc.no/v1", out.APIVersion)
	assert.Equal(t, "FormDefinition", out.Kind)
	assert.Equal(t, "nrc.no", out.Spec.Group)
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

	out.Spec.Names.Plural = "abc"
	updated, err := nrcClient.FormDefinitions().Update(ctx, out)
	if err != nil {
		t.Errorf("unable to update form definition: %v", err)
		return
	}

	t.Logf("%#v", updated)

	time.Sleep(10 * time.Second)

}

func TestUpdateFormDefinition(t *testing.T) {

}
