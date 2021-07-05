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
	HydraHTTPClient       *http.Client
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

	var hydraAdmin admin.ClientService
	if o.HydraAdminClient != nil {
		hydraAdmin = o.HydraAdminClient.Admin
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
		HydraAdmin:            hydraAdmin,
		HydraHTTPClient:       o.HydraHTTPClient,
	}

	router := mux.NewRouter()
	router.Use(srv.WithAuth())

	var attributesEP = server.Endpoints["attributes"]
	router.Path(attributesEP).Methods("GET").HandlerFunc(srv.ListAttributes)
	router.Path(attributesEP).Methods("POST").HandlerFunc(srv.PostAttribute)
	router.Path(attributesEP + "/{id}").Methods("GET").HandlerFunc(srv.GetAttribute)
	router.Path(attributesEP + "/{id}").Methods("PUT").HandlerFunc(srv.PutAttribute)

	var individualsEP = server.Endpoints["individuals"]
	router.Path(individualsEP).Methods("GET").HandlerFunc(srv.ListIndividuals)
	router.Path(individualsEP).Methods("POST").HandlerFunc(srv.PostIndividual)
	router.Path(individualsEP + "/{id}").Methods("GET").HandlerFunc(srv.GetIndividual)
	router.Path(individualsEP + "/{id}").Methods("PUT").HandlerFunc(srv.PutIndividual)

	var membershipsEP = server.Endpoints["memberships"]
	router.Path(membershipsEP).Methods("GET").HandlerFunc(srv.ListMemberships)
	router.Path(membershipsEP).Methods("POST").HandlerFunc(srv.PostMembership)
	router.Path(membershipsEP + "/{v1}").Methods("GET").HandlerFunc(srv.GetMembership)

	var partiesEP = server.Endpoints["parties"]
	router.Path(partiesEP).Methods("GET").HandlerFunc(srv.ListParties)
	router.Path(partiesEP).Methods("POST").HandlerFunc(srv.PostParty)
	router.Path(partiesEP + "/search").Methods("POST").HandlerFunc(srv.SearchParties)
	router.Path(partiesEP + "/{id}").Methods("GET").HandlerFunc(srv.GetParty)
	router.Path(partiesEP + "/{id}").Methods("PUT").HandlerFunc(srv.PutParty)

	var partyTypesEP = server.Endpoints["partytypes"]
	router.Path(partyTypesEP).Methods("GET").HandlerFunc(srv.ListPartyTypes)
	router.Path(partyTypesEP).Methods("POST").HandlerFunc(srv.PostPartyType)
	router.Path(partyTypesEP + "/{id}").Methods("GET").HandlerFunc(srv.GetPartyType)
	router.Path(partyTypesEP + "/{id}").Methods("PUT").HandlerFunc(srv.PutPartyType)

	var relationshipsEP = server.Endpoints["relationships"]
	router.Path(relationshipsEP).Methods("GET").HandlerFunc(srv.ListRelationships)
	router.Path(relationshipsEP).Methods("POST").HandlerFunc(srv.PostRelationship)
	router.Path(relationshipsEP + "/{id}").Methods("GET").HandlerFunc(srv.GetRelationship)
	router.Path(relationshipsEP + "/{id}").Methods("PUT").HandlerFunc(srv.PutRelationship)
	router.Path(relationshipsEP + "/{id}").Methods("DELETE").HandlerFunc(srv.DeleteRelationship)

	var relationshipTypesEP = server.Endpoints["relationshiptypes"]
	router.Path(relationshipTypesEP).Methods("GET").HandlerFunc(srv.ListRelationshipTypes)
	router.Path(relationshipTypesEP).Methods("POST").HandlerFunc(srv.PostRelationshipType)
	router.Path(relationshipTypesEP + "/{id}").Methods("GET").HandlerFunc(srv.GetRelationshipType)
	router.Path(relationshipTypesEP + "/{id}").Methods("PUT").HandlerFunc(srv.PutRelationshipType)

	var teamsEP = server.Endpoints["teams"]
	router.Path(teamsEP).Methods("GET").HandlerFunc(srv.ListTeams)
	router.Path(teamsEP).Methods("POST").HandlerFunc(srv.PostTeam)
	router.Path(teamsEP + "/{id}").Methods("GET").HandlerFunc(srv.GetTeam)
	router.Path(teamsEP + "/{id}").Methods("PUT").HandlerFunc(srv.PutTeam)

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
