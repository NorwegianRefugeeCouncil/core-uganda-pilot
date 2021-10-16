package form

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/nrc-no/core/pkg/validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net"
	"net/http"
	"net/url"
	"strings"
)

const (
	collFormDefinitions = "form_definitions"
)

type Server struct {
	uidGenerator  utils.UIDGenerator
	mongoClientFn func() (*mongo.Client, error)
	databaseName  string
	listener      net.Listener
}

func NewServer(
	uidGenerator utils.UIDGenerator,
	mongoClientFn func() (*mongo.Client, error),
	databaseName string,
) *Server {
	srv := &Server{
		uidGenerator:  uidGenerator,
		mongoClientFn: mongoClientFn,
		databaseName:  databaseName,
	}
	return srv
}

func (s *Server) Start(done chan struct{}) error {
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		return err
	}
	s.listener = listener
	go func() {
		if err := http.Serve(s.listener, s); err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Print(err)
		}
	}()
	go func() {
		<-done
		s.listener.Close()
	}()
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	path := req.URL.Path
	if strings.HasPrefix(path, server.FormDefinitionsEndpoint) {
		path = strings.TrimPrefix(path, server.FormDefinitionsEndpoint)
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if req.Method == "GET" && path == "/" {

	} else if req.Method == "GET" {
		s.HandleGetFormDefinition(w, req)
	} else if req.Method == "POST" && path == "/" {
		s.HandlePostFormDefinition(w, req)
	} else if req.Method == "POST" {
		s.HandleValidateForm(w, req)
	} else if req.Method == "PUT" {
		s.HandlePutFormDefinition(w, req)
	}

}

func (s *Server) HandleValidateForm(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var submission Submission
	if err := utils.BindJSON(req, &submission); err != nil {
		utils.ErrorResponse(w, fmt.Errorf("failed to unmarshal json payload: %v", err))
		return
	}

	mongoClient, err := s.mongoClientFn()
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("failed to connect to database: %v", err))
		return
	}

	definition, err := s.getDefinition(ctx, mongoClient, submission.FormDefinitionID, submission.FormDefinitionVersion)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("failed to get form definition: %v", err))
		return
	}

	if err := validateSubmission(definition, &submission); err != nil {
		utils.ErrorResponse(w, fmt.Errorf("validation failed: %v", err))
		return
	}

	responseBytes, err := json.Marshal(submission)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("failed to marshal form: %v", err))
		return
	}

	w.WriteHeader(200)
	_, _ = w.Write(responseBytes)

}

func (s *Server) HandleGetFormDefinition(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	id := req.URL.Path
	if strings.HasPrefix(id, server.FormDefinitionsEndpoint) {
		id = strings.TrimPrefix(id, server.FormDefinitionsEndpoint)
	}

	if strings.HasPrefix(id, "/") {
		id = strings.TrimPrefix(id, "/")
	}

	mongoClient, err := s.mongoClientFn()
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("failed to connect to database: %v", err))
		return
	}

	collection := mongoClient.Database(s.databaseName).Collection(collFormDefinitions)

	filter := bson.M{
		"id": id,
	}

	version := getVersionFromQueryParam(req.URL.Query())
	if len(version) > 0 {
		filter["version"] = version
	} else {
		filter["isCurrentVersion"] = true
	}

	result := collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			utils.ErrorResponse(w, fmt.Errorf("form definition not found"))
			return
		} else {
			utils.ErrorResponse(w, fmt.Errorf("failed to get form definition"))
			return
		}
	}

	var out Definition
	if err := result.Decode(&out); err != nil {
		utils.ErrorResponse(w, fmt.Errorf("failed to unmarshal form definition: %v", err))
		return
	}

	utils.JSONResponse(w, http.StatusOK, &out)

}

func (s *Server) HandlePostFormDefinition(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var definition Definition
	if err := utils.BindJSON(req, &definition); err != nil {
		utils.ErrorResponse(w, fmt.Errorf("failed to unmarshal form definition: %v", err))
		return
	}
	definition.Version = s.uidGenerator.GenUID()
	definition.ID = s.uidGenerator.GenUID()
	definition.IsCurrentVersion = true

	mongoClient, err := s.mongoClientFn()
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("failed to connect to database: %v", err))
		return
	}

	collection := s.getCollection(mongoClient)

	_, err = collection.InsertOne(ctx, definition)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("failed to save form definition: %v", err))
		return
	}

	utils.JSONResponse(w, http.StatusOK, definition)

}
func (s *Server) HandlePutFormDefinition(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var definition Definition
	if err := utils.BindJSON(req, &definition); err != nil {
		utils.ErrorResponse(w, fmt.Errorf("failed to unmarshal form definition: %v", err))
		return
	}
	definition.Version = s.uidGenerator.GenUID()

	mongoClient, err := s.mongoClientFn()
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("failed to connect to database: %v", err))
		return
	}

	collection := s.getCollection(mongoClient)

	if definition.IsCurrentVersion {
		_, err = collection.UpdateMany(ctx, bson.M{
			"id":               definition.ID,
			"isCurrentVersion": true,
		}, bson.M{
			"$set": bson.M{
				"isCurrentVersion": false,
			},
		})
		if err != nil {
			utils.ErrorResponse(w, fmt.Errorf("failed to update current version: %v", err))
			return
		}
	}

	_, err = collection.InsertOne(ctx, definition)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("failed to save form definition: %v", err))
		return
	}

	utils.JSONResponse(w, http.StatusOK, definition)

}

func getVersionFromQueryParam(v url.Values) string {
	return v.Get("version")
}

func validateSubmission(definition *Definition, submission *Submission) error {
	return nil
}

func (s *Server) getDefinition(
	ctx context.Context,
	mongoClient *mongo.Client,
	id string,
	revision string,
) (*Definition, error) {

	collection := s.getCollection(mongoClient)

	filter := bson.M{
		"id": id,
	}

	if len(revision) > 0 {
		filter["revision"] = revision
	}

	result := collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, validation.Status{
				Status:  validation.Failure,
				Code:    http.StatusNotFound,
				Message: "object not found",
			}
		} else {
			return nil, fmt.Errorf("failed to find form definition: %v", result.Err())
		}
	}

	var definition Definition
	if err := result.Decode(&definition); err != nil {
		return nil, fmt.Errorf("failed to decode form definition: %v", err)
	}

	return &definition, nil

}

func (s *Server) getCollection(mongoClient *mongo.Client) *mongo.Collection {
	collection := mongoClient.
		Database(s.databaseName).
		Collection(collFormDefinitions)
	return collection
}

func (s *Server) NewClient() Interface {
	cli := NewClientFromConfig(&rest.Config{
		Scheme:     "http",
		Host:       s.listener.Addr().String(),
		HTTPClient: http.DefaultClient,
	})
	return cli
}

type Submission struct {
	ID                    string                 `bson:"id"`
	FormDefinitionID      string                 `bson:"formDefinitionId"`
	FormDefinitionVersion string                 `bson:"formDefinitionVersion"`
	Data                  map[string]interface{} `bson:"data"`
}

type Definition struct {
	ID               string `bson:"id"`
	Version          string `bson:"version"`
	Name             string `bson:"name"`
	IsCurrentVersion bool   `bson:"isCurrentVersion"`
}

type DefinitionList struct {
	Items []*Definition
}
