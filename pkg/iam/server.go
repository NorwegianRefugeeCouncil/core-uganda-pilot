package iam

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
	"net/http"
	"path"
)

type Server struct {
	environment                     string
	router                          *mux.Router
	partyAttributeDefinitionStore   *PartyAttributeDefinitionStore
	partyStore                      *PartyStore
	partyTypeStore                  *PartyTypeStore
	relationshipStore               *RelationshipStore
	relationshipTypeStore           *RelationshipTypeStore
	identificationDocumentStore     *IdentificationDocumentStore
	identificationDocumentTypeStore *IdentificationDocumentTypeStore
	individualStore                 *IndividualStore
	teamStore                       *TeamStore
	countryStore                    *CountryStore
	membershipStore                 *MembershipStore
	nationalityStore                *NationalityStore
	hydraAdmin                      admin.ClientService
	mongoClientFn                   utils.MongoClientFn
	hydraHTTPClient                 *http.Client
}

func NewServerOrDie(ctx context.Context, o *server.GenericServerOptions) *Server {
	srv, err := NewServer(ctx, o)
	if err != nil {
		panic(err)
	}
	return srv
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

	identificationDocumentStore, err := newIdentificationDocumentStore(ctx, o.MongoClientFn, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	identificationDocumentTypeStore, err := newIdentificationDocumentTypeStore(ctx, o.MongoClientFn, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	hydraAdmin := o.HydraAdminClient.Admin

	srv := &Server{
		environment:                     o.Environment,
		mongoClientFn:                   o.MongoClientFn,
		partyAttributeDefinitionStore:   attributeStore,
		countryStore:                    NewCountryStore(partyStore),
		partyStore:                      partyStore,
		partyTypeStore:                  partyTypeStore,
		relationshipStore:               relationshipStore,
		relationshipTypeStore:           relationshipTypeStore,
		identificationDocumentStore:     identificationDocumentStore,
		identificationDocumentTypeStore: identificationDocumentTypeStore, individualStore: NewIndividualStore(o.MongoClientFn, o.MongoDatabase),
		teamStore:        NewTeamStore(partyStore),
		membershipStore:  NewMembershipStore(relationshipStore),
		nationalityStore: NewNationalityStore(relationshipStore),
		hydraAdmin:       hydraAdmin,
		hydraHTTPClient:  o.HydraHTTPClient,
	}

	router := mux.NewRouter()
	router.Use(srv.withAuth())

	router.Path(server.AttributesEndpoint).Methods(http.MethodGet).HandlerFunc(srv.listPartyAttributeDefinitions)
	router.Path(server.AttributesEndpoint).Methods(http.MethodPost).HandlerFunc(srv.postPartyAttributeDefinition)
	router.Path(path.Join(server.AttributesEndpoint, "{id}")).Methods(http.MethodGet).HandlerFunc(srv.getPartyAttributeDefinition)
	router.Path(path.Join(server.AttributesEndpoint, "{id}")).Methods(http.MethodPut).HandlerFunc(srv.putPartyAttributeDefinition)

	router.Path(server.IndividualsEndpoint).Methods(http.MethodGet).HandlerFunc(srv.listIndividuals)
	router.Path(server.IndividualsEndpoint).Methods(http.MethodPost).HandlerFunc(srv.postIndividual)
	router.Path(path.Join(server.IndividualsEndpoint, "{id}")).Methods(http.MethodGet).HandlerFunc(srv.getIndividual)
	router.Path(path.Join(server.IndividualsEndpoint, "{id}")).Methods(http.MethodPut).HandlerFunc(srv.putIndividual)

	router.Path(server.MembershipsEndpoint).Methods(http.MethodGet).HandlerFunc(srv.listMemberships)
	router.Path(server.MembershipsEndpoint).Methods(http.MethodPost).HandlerFunc(srv.postMembership)
	router.Path(path.Join(server.MembershipsEndpoint, "{id}")).Methods(http.MethodGet).HandlerFunc(srv.getMembership)

	router.Path(server.NationalitiesEndpoint).Methods(http.MethodGet).HandlerFunc(srv.listNationalities)
	router.Path(server.NationalitiesEndpoint).Methods(http.MethodPost).HandlerFunc(srv.postNationality)
	router.Path(path.Join(server.NationalitiesEndpoint, "{id}")).Methods(http.MethodGet).HandlerFunc(srv.getNationality)

	router.Path(server.PartiesEndpoint).Methods(http.MethodGet).HandlerFunc(srv.listParties)
	router.Path(server.PartiesEndpoint).Methods(http.MethodPost).HandlerFunc(srv.postParty)
	router.Path(path.Join(server.PartiesEndpoint, "/search")).Methods(http.MethodPost).HandlerFunc(srv.searchParties)
	router.Path(path.Join(server.PartiesEndpoint, "{id}")).Methods(http.MethodGet).HandlerFunc(srv.getParty)
	router.Path(path.Join(server.PartiesEndpoint, "{id}")).Methods(http.MethodPut).HandlerFunc(srv.putParty)

	router.Path(server.PartyTypesEndpoint).Methods(http.MethodGet).HandlerFunc(srv.listPartyTypes)
	router.Path(server.PartyTypesEndpoint).Methods(http.MethodPost).HandlerFunc(srv.postPartyType)
	router.Path(path.Join(server.PartyTypesEndpoint, "{id}")).Methods(http.MethodGet).HandlerFunc(srv.getPartyType)
	router.Path(path.Join(server.PartyTypesEndpoint, "{id}")).Methods(http.MethodPut).HandlerFunc(srv.putPartyType)

	router.Path(server.RelationshipsEndpoint).Methods(http.MethodGet).HandlerFunc(srv.listRelationships)
	router.Path(server.RelationshipsEndpoint).Methods(http.MethodPost).HandlerFunc(srv.postRelationship)
	router.Path(path.Join(server.RelationshipsEndpoint, "{id}")).Methods(http.MethodGet).HandlerFunc(srv.getRelationship)
	router.Path(path.Join(server.RelationshipsEndpoint, "{id}")).Methods(http.MethodPut).HandlerFunc(srv.putRelationship)
	router.Path(path.Join(server.RelationshipsEndpoint, "{id}")).Methods(http.MethodDelete).HandlerFunc(srv.deleteRelationship)

	router.Path(server.RelationshipTypesEndpoint).Methods(http.MethodGet).HandlerFunc(srv.listRelationshipTypes)
	router.Path(server.RelationshipTypesEndpoint).Methods(http.MethodPost).HandlerFunc(srv.postRelationshipType)
	router.Path(path.Join(server.RelationshipTypesEndpoint, "{id}")).Methods(http.MethodGet).HandlerFunc(srv.getRelationshipType)
	router.Path(path.Join(server.RelationshipTypesEndpoint, "{id}")).Methods(http.MethodPut).HandlerFunc(srv.putRelationshipType)

	router.Path(server.TeamsEndpoint).Methods(http.MethodGet).HandlerFunc(srv.listTeams)
	router.Path(server.TeamsEndpoint).Methods(http.MethodPost).HandlerFunc(srv.postTeam)
	router.Path(path.Join(server.TeamsEndpoint, "{id}")).Methods(http.MethodGet).HandlerFunc(srv.getTeam)
	router.Path(path.Join(server.TeamsEndpoint, "{id}")).Methods(http.MethodPut).HandlerFunc(srv.putTeam)

	router.Path(server.CountriesEndpoint).Methods(http.MethodGet).HandlerFunc(srv.listCountries)
	router.Path(server.CountriesEndpoint).Methods(http.MethodPost).HandlerFunc(srv.postCountry)
	router.Path(path.Join(server.CountriesEndpoint, "{id}")).Methods(http.MethodGet).HandlerFunc(srv.getCountry)
	router.Path(path.Join(server.CountriesEndpoint, "{id}")).Methods(http.MethodPut).HandlerFunc(srv.putCountry)

	router.Path(server.IdentificationDocumentsEndpoint).Methods(http.MethodGet).HandlerFunc(srv.listIdentificationDocuments)
	router.Path(server.IdentificationDocumentsEndpoint).Methods(http.MethodPost).HandlerFunc(srv.postIdentificationDocument)
	router.Path(path.Join(server.IdentificationDocumentsEndpoint, "{id}")).Methods(http.MethodGet).HandlerFunc(srv.getIdentificationDocument)
	router.Path(path.Join(server.IdentificationDocumentsEndpoint, "{id}")).Methods(http.MethodPut).HandlerFunc(srv.putIdentificationDocument)
	router.Path(path.Join(server.IdentificationDocumentsEndpoint, "{id}")).Methods(http.MethodDelete).HandlerFunc(srv.deleteIdentificationDocument)

	router.Path(server.IdentificationDocumentTypesEndpoint).Methods(http.MethodGet).HandlerFunc(srv.listIdentificationDocumentTypes)
	router.Path(server.IdentificationDocumentTypesEndpoint).Methods(http.MethodPost).HandlerFunc(srv.postIdentificationDocumentType)
	router.Path(path.Join(server.IdentificationDocumentTypesEndpoint, "{id}")).Methods(http.MethodGet).HandlerFunc(srv.getIdentificationDocumentType)
	router.Path(path.Join(server.IdentificationDocumentTypesEndpoint, "{id}")).Methods(http.MethodPut).HandlerFunc(srv.putIdentificationDocumentType)

	srv.router = router

	return srv, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Server) json(w http.ResponseWriter, status int, data interface{}) {
	utils.JSONResponse(w, status, data)
}

func (s *Server) getPathParam(param string, w http.ResponseWriter, req *http.Request, into *string) bool {
	return utils.GetPathParam(param, w, req, into)
}

func (s *Server) error(w http.ResponseWriter, err error) {
	utils.ErrorResponse(w, err)
}

func (s *Server) bind(req *http.Request, into interface{}) error {
	return utils.BindJSON(req, into)
}

func (s *Server) ResetDB(ctx context.Context, databaseName string) error {
	mongoClient, err := s.mongoClientFn(ctx)
	if err != nil {
		return err
	}
	if err := mongoClient.Database(databaseName).Drop(ctx); err != nil {
		return err
	}
	return s.Init(ctx)
}
