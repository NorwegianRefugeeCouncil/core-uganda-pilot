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
	"path"
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

	router.Path(server.AttributesEndpoint).Methods("GET").HandlerFunc(srv.ListAttributes)
	router.Path(server.AttributesEndpoint).Methods("POST").HandlerFunc(srv.PostAttribute)
	router.Path(path.Join(server.AttributesEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.GetAttribute)
	router.Path(path.Join(server.AttributesEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.PutAttribute)

	router.Path(server.IndividualsEndpoint).Methods("GET").HandlerFunc(srv.ListIndividuals)
	router.Path(server.IndividualsEndpoint).Methods("POST").HandlerFunc(srv.PostIndividual)
	router.Path(path.Join(server.IndividualsEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.GetIndividual)
	router.Path(path.Join(server.IndividualsEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.PutIndividual)

	router.Path(server.MembershipsEndpoint).Methods("GET").HandlerFunc(srv.ListMemberships)
	router.Path(server.MembershipsEndpoint).Methods("POST").HandlerFunc(srv.PostMembership)
	router.Path(path.Join(server.MembershipsEndpoint, "{v1}")).Methods("GET").HandlerFunc(srv.GetMembership)

	router.Path(server.PartiesEndpoint).Methods("GET").HandlerFunc(srv.ListParties)
	router.Path(server.PartiesEndpoint).Methods("POST").HandlerFunc(srv.PostParty)
	router.Path(path.Join(server.PartiesEndpoint, "/search")).Methods("POST").HandlerFunc(srv.SearchParties)
	router.Path(path.Join(server.PartiesEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.GetParty)
	router.Path(path.Join(server.PartiesEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.PutParty)

	router.Path(server.PartyTypesEndpoint).Methods("GET").HandlerFunc(srv.ListPartyTypes)
	router.Path(server.PartyTypesEndpoint).Methods("POST").HandlerFunc(srv.PostPartyType)
	router.Path(path.Join(server.PartyTypesEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.GetPartyType)
	router.Path(path.Join(server.PartyTypesEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.PutPartyType)

	router.Path(server.RelationshipsEndpoint).Methods("GET").HandlerFunc(srv.ListRelationships)
	router.Path(server.RelationshipsEndpoint).Methods("POST").HandlerFunc(srv.PostRelationship)
	router.Path(path.Join(server.RelationshipsEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.GetRelationship)
	router.Path(path.Join(server.RelationshipsEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.PutRelationship)
	router.Path(path.Join(server.RelationshipsEndpoint, "{id}")).Methods("DELETE").HandlerFunc(srv.DeleteRelationship)

	router.Path(server.RelationshipTypesEndpoint).Methods("GET").HandlerFunc(srv.ListRelationshipTypes)
	router.Path(server.RelationshipTypesEndpoint).Methods("POST").HandlerFunc(srv.PostRelationshipType)
	router.Path(path.Join(server.RelationshipTypesEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.GetRelationshipType)
	router.Path(path.Join(server.RelationshipTypesEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.PutRelationshipType)

	router.Path(server.TeamsEndpoint).Methods("GET").HandlerFunc(srv.ListTeams)
	router.Path(server.TeamsEndpoint).Methods("POST").HandlerFunc(srv.PostTeam)
	router.Path(path.Join(server.TeamsEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.GetTeam)
	router.Path(path.Join(server.TeamsEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.PutTeam)

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
	_, err = w.Write(responseBytes)
	if err != nil {
		return
	}
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
