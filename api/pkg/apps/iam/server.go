package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
	"io/ioutil"
	"net/http"
	"path"
)

type Server struct {
	environment           string
	router                *mux.Router
	attributeStore        *AttributeStore
	partyStore            *PartyStore
	partyTypeStore        *PartyTypeStore
	relationshipStore     *RelationshipStore
	relationshipTypeStore *RelationshipTypeStore
	individualStore       *IndividualStore
	teamStore             *TeamStore
	membershipStore       *MembershipStore
	hydraAdmin            admin.ClientService
	mongoClientFn         utils.MongoClientFn
	hydraHTTPClient       *http.Client
}

func NewServer(ctx context.Context, o *server.GenericServerOptions) (*Server, error) {

	relationshipStore, err := newRelationshipStore(ctx, o.MongoClientFn, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	partyStore, err := newPartyStore(ctx, o.MongoClientFn, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	attributeStore, err := newAttributeStore(ctx, o.MongoClientFn, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	partyTypeStore, err := newPartyTypeStore(ctx, o.MongoClientFn, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	relationshipTypeStore, err := newRelationshipTypeStore(ctx, o.MongoClientFn, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	var hydraAdmin admin.ClientService
	if o.HydraAdminClient != nil {
		hydraAdmin = o.HydraAdminClient.Admin
	}

	srv := &Server{
		environment:           o.Environment,
		mongoClientFn:         o.MongoClientFn,
		attributeStore:        attributeStore,
		partyStore:            partyStore,
		partyTypeStore:        partyTypeStore,
		relationshipStore:     relationshipStore,
		relationshipTypeStore: relationshipTypeStore,
		individualStore:       NewIndividualStore(o.MongoClientFn, o.MongoDatabase),
		teamStore:             NewTeamStore(partyStore),
		membershipStore:       NewMembershipStore(relationshipStore),
		hydraAdmin:            hydraAdmin,
		hydraHTTPClient:       o.HydraHTTPClient,
	}

	router := mux.NewRouter()
	router.Use(srv.withAuth())

	router.Path(server.AttributesEndpoint).Methods("GET").HandlerFunc(srv.listAttributes)
	router.Path(server.AttributesEndpoint).Methods("POST").HandlerFunc(srv.postAttributes)
	router.Path(path.Join(server.AttributesEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.getAttribute)
	router.Path(path.Join(server.AttributesEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.putAttribute)

	router.Path(server.IndividualsEndpoint).Methods("GET").HandlerFunc(srv.listIndividuals)
	router.Path(server.IndividualsEndpoint).Methods("POST").HandlerFunc(srv.postIndividual)
	router.Path(path.Join(server.IndividualsEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.getIndividual)
	router.Path(path.Join(server.IndividualsEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.putIndividual)

	router.Path(server.MembershipsEndpoint).Methods("GET").HandlerFunc(srv.listMemberships)
	router.Path(server.MembershipsEndpoint).Methods("POST").HandlerFunc(srv.postMembership)
	router.Path(path.Join(server.MembershipsEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.getMembership)

	router.Path(server.PartiesEndpoint).Methods("GET").HandlerFunc(srv.listParties)
	router.Path(server.PartiesEndpoint).Methods("POST").HandlerFunc(srv.postParty)
	router.Path(path.Join(server.PartiesEndpoint, "/search")).Methods("POST").HandlerFunc(srv.searchParties)
	router.Path(path.Join(server.PartiesEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.getParty)
	router.Path(path.Join(server.PartiesEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.putParty)

	router.Path(server.PartyTypesEndpoint).Methods("GET").HandlerFunc(srv.listPartyTypes)
	router.Path(server.PartyTypesEndpoint).Methods("POST").HandlerFunc(srv.postPartyType)
	router.Path(path.Join(server.PartyTypesEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.getPartyType)
	router.Path(path.Join(server.PartyTypesEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.putPartyType)

	router.Path(server.RelationshipsEndpoint).Methods("GET").HandlerFunc(srv.listRelationships)
	router.Path(server.RelationshipsEndpoint).Methods("POST").HandlerFunc(srv.postRelationship)
	router.Path(path.Join(server.RelationshipsEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.getRelationship)
	router.Path(path.Join(server.RelationshipsEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.putRelationship)
	router.Path(path.Join(server.RelationshipsEndpoint, "{id}")).Methods("DELETE").HandlerFunc(srv.deleteRelationship)

	router.Path(server.RelationshipTypesEndpoint).Methods("GET").HandlerFunc(srv.listRelationshipTypes)
	router.Path(server.RelationshipTypesEndpoint).Methods("POST").HandlerFunc(srv.postRelationshipType)
	router.Path(path.Join(server.RelationshipTypesEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.getRelationshipType)
	router.Path(path.Join(server.RelationshipTypesEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.putRelationshipType)

	router.Path(server.TeamsEndpoint).Methods("GET").HandlerFunc(srv.listTeams)
	router.Path(server.TeamsEndpoint).Methods("POST").HandlerFunc(srv.postTeam)
	router.Path(path.Join(server.TeamsEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.getTeam)
	router.Path(path.Join(server.TeamsEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.putTeam)

	srv.router = router

	return srv, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Server) json(w http.ResponseWriter, status int, data interface{}) {
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

func (s *Server) getPathParam(param string, w http.ResponseWriter, req *http.Request, into *string) bool {
	id, ok := mux.Vars(req)[param]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("path parameter '%s' not found in path", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	*into = id
	return true
}

func (s *Server) error(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (s *Server) bind(req *http.Request, into interface{}) error {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bodyBytes, &into); err != nil {
		return err
	}

	return nil
}

func (s *Server) ResetDB(ctx context.Context, databaseName string) error {
	mongoClient, err := s.mongoClientFn(ctx)
	if err != nil {
		return err
	}
	if err := mongoClient.Database(databaseName).Drop(ctx); err != nil {
		return err
	}
	if err := s.Init(ctx); err != nil {
		return err
	}
	return nil
}
