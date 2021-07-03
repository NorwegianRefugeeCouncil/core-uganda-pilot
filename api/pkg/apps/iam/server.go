package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/ory/hydra-client-go/client/admin"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
)

type Server struct {
	environment           string
	router                *mux.Router
	AttributeStore        *AttributeStore
	PartyStore            *PartyStore
	PartyTypeStore        *PartyTypeStore
	RelationshipStore     *RelationshipStore
	RelationshipTypeStore *RelationshipTypeStore
	IndividualStore       *IndividualStore
	TeamStore             *TeamStore
	MembershipStore       *MembershipStore
	HydraAdmin            admin.ClientService
	mongoClient           *mongo.Client
}

func NewServer(ctx context.Context, o *server.GenericServerOptions) (*Server, error) {

	relationshipStore, err := NewRelationshipStore(ctx, o.MongoClient, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	partyStore, err := NewPartyStore(ctx, o.MongoClient, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	attributeStore, err := NewAttributeStore(ctx, o.MongoClient, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	partyTypeStore, err := NewPartyTypeStore(ctx, o.MongoClient, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	relationshipTypeStore, err := NewRelationshipTypeStore(ctx, o.MongoClient, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		environment:           o.Environment,
		mongoClient:           o.MongoClient,
		AttributeStore:        attributeStore,
		PartyStore:            partyStore,
		PartyTypeStore:        partyTypeStore,
		RelationshipStore:     relationshipStore,
		RelationshipTypeStore: relationshipTypeStore,
		IndividualStore:       NewIndividualStore(o.MongoClient, o.MongoDatabase),
		TeamStore:             NewTeamStore(partyStore),
		MembershipStore:       NewMembershipStore(relationshipStore),
		HydraAdmin:            o.HydraAdminClient.Admin,
	}

	router := mux.NewRouter()
	router.Use(srv.WithAuth())

	router.Path("/apis/iam/v1/attributes").Methods("GET").HandlerFunc(srv.ListAttributes)
	router.Path("/apis/iam/v1/attributes").Methods("POST").HandlerFunc(srv.PostAttribute)
	router.Path("/apis/iam/v1/attributes/{id}").Methods("GET").HandlerFunc(srv.GetAttribute)
	router.Path("/apis/iam/v1/attributes/{id}").Methods("PUT").HandlerFunc(srv.PutAttribute)

	router.Path("/apis/iam/v1/individuals").Methods("GET").HandlerFunc(srv.ListIndividuals)
	router.Path("/apis/iam/v1/individuals").Methods("POST").HandlerFunc(srv.PostIndividual)
	router.Path("/apis/iam/v1/individuals/{id}").Methods("GET").HandlerFunc(srv.GetIndividual)
	router.Path("/apis/iam/v1/individuals/{id}").Methods("PUT").HandlerFunc(srv.PutIndividual)

	router.Path("/apis/iam/v1/memberships").Methods("GET").HandlerFunc(srv.ListMemberships)
	router.Path("/apis/iam/v1/memberships").Methods("POST").HandlerFunc(srv.PostMembership)
	router.Path("/apis/iam/v1/memberships/{v1}").Methods("GET").HandlerFunc(srv.GetMembership)

	router.Path("/apis/iam/v1/parties").Methods("GET").HandlerFunc(srv.ListParties)
	router.Path("/apis/iam/v1/parties").Methods("POST").HandlerFunc(srv.PostParty)
	router.Path("/apis/iam/v1/parties/search").Methods("POST").HandlerFunc(srv.SearchParties)
	router.Path("/apis/iam/v1/parties/{id}").Methods("GET").HandlerFunc(srv.GetParty)
	router.Path("/apis/iam/v1/parties/{id}").Methods("PUT").HandlerFunc(srv.PutParty)

	router.Path("/apis/iam/v1/partytypes").Methods("GET").HandlerFunc(srv.ListPartyTypes)
	router.Path("/apis/iam/v1/partytypes").Methods("POST").HandlerFunc(srv.PostPartyType)
	router.Path("/apis/iam/v1/partytypes/{id}").Methods("GET").HandlerFunc(srv.GetPartyType)
	router.Path("/apis/iam/v1/partytypes/{id}").Methods("PUT").HandlerFunc(srv.PutPartyType)

	router.Path("/apis/iam/v1/relationships").Methods("GET").HandlerFunc(srv.ListRelationships)
	router.Path("/apis/iam/v1/relationships").Methods("POST").HandlerFunc(srv.PostRelationship)
	router.Path("/apis/iam/v1/relationships/{id}").Methods("GET").HandlerFunc(srv.GetRelationship)
	router.Path("/apis/iam/v1/relationships/{id}").Methods("PUT").HandlerFunc(srv.PutRelationship)
	router.Path("/apis/iam/v1/relationships/{id}").Methods("DELETE").HandlerFunc(srv.DeleteRelationship)

	router.Path("/apis/iam/v1/relationshiptypes").Methods("GET").HandlerFunc(srv.ListRelationshipTypes)
	router.Path("/apis/iam/v1/relationshiptypes").Methods("POST").HandlerFunc(srv.PostRelationshipType)
	router.Path("/apis/iam/v1/relationshiptypes/{id}").Methods("GET").HandlerFunc(srv.GetRelationshipType)
	router.Path("/apis/iam/v1/relationshiptypes/{id}").Methods("PUT").HandlerFunc(srv.PutRelationshipType)

	router.Path("/apis/iam/v1/teams").Methods("GET").HandlerFunc(srv.ListTeams)
	router.Path("/apis/iam/v1/teams").Methods("POST").HandlerFunc(srv.PostTeam)
	router.Path("/apis/iam/v1/teams/{id}").Methods("GET").HandlerFunc(srv.GetTeam)
	router.Path("/apis/iam/v1/teams/{id}").Methods("PUT").HandlerFunc(srv.PutTeam)

	srv.router = router

	return srv, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Server) JSON(w http.ResponseWriter, status int, data interface{}) {
	responseBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}

func (s *Server) GetPathParam(param string, w http.ResponseWriter, req *http.Request, into *string) bool {
	id, ok := mux.Vars(req)[param]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("path parameter '%s' not found in path", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	*into = id
	return true
}

func (s *Server) Error(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (s *Server) Bind(req *http.Request, into interface{}) error {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bodyBytes, &into); err != nil {
		return err
	}

	return nil
}
