package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"net/http"
)

type ServerOptions struct {
	ListenAddress string
	MongoHosts    []string
	MongoDatabase string
	MongoUsername string
	MongoPassword string
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		ListenAddress: ":9001",
		MongoHosts:    []string{"mongo://localhost:27017"},
	}
}

func (o *ServerOptions) WithMongoHosts(hosts []string) *ServerOptions {
	o.MongoHosts = hosts
	return o
}
func (o *ServerOptions) WithMongoDatabase(mongoDatabase string) *ServerOptions {
	o.MongoDatabase = mongoDatabase
	return o
}
func (o *ServerOptions) WithMongoUsername(mongoUsername string) *ServerOptions {
	o.MongoUsername = mongoUsername
	return o
}
func (o *ServerOptions) WithMongoPassword(mongoPassword string) *ServerOptions {
	o.MongoPassword = mongoPassword
	return o
}
func (o *ServerOptions) WithListenAddress(address string) *ServerOptions {
	o.ListenAddress = address
	return o
}

func (o *ServerOptions) Flags(fs pflag.FlagSet) {
	fs.StringVar(&o.ListenAddress, "listen-address", o.ListenAddress, "Server listen address")
	fs.StringSliceVar(&o.MongoHosts, "mongo-url", o.MongoHosts, "Mongo url")
	fs.StringVar(&o.MongoDatabase, "mongo-database", o.MongoDatabase, "Mongo database")
	fs.StringVar(&o.MongoUsername, "mongo-username", o.MongoUsername, "Mongo username")
	fs.StringVar(&o.MongoPassword, "mongo-password", o.MongoPassword, "Mongo password")
}

type Server struct {
	router                *mux.Router
	AttributeStore        *AttributeStore
	PartyStore            *PartyStore
	PartyTypeStore        *PartyTypeStore
	RelationshipStore     *RelationshipStore
	RelationshipTypeStore *RelationshipTypeStore
	IndividualStore       *IndividualStore
	StaffStore            *StaffStore
	OrganizationStore     *OrganizationStore
	TeamStore             *TeamStore
	MembershipStore       *MembershipStore
	mongoClient           *mongo.Client
}

func NewServer(ctx context.Context, o *ServerOptions) (*Server, error) {

	mongoClient, err := mongo.NewClient(
		options.Client().
			SetHosts(o.MongoHosts).
			SetAuth(options.Credential{
				Username: o.MongoUsername,
				Password: o.MongoPassword,
			}))
	if err != nil {
		return nil, err
	}

	if err := mongoClient.Connect(ctx); err != nil {
		return nil, err
	}

	router := mux.NewRouter()

	relationshipStore, err := NewRelationshipStore(ctx, mongoClient, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	partyStore, err := NewPartyStore(ctx, mongoClient, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	attributeStore, err := NewAttributeStore(ctx, mongoClient, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	partyTypeStore, err := NewPartyTypeStore(ctx, mongoClient, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	relationshipTypeStore, err := NewRelationshipTypeStore(ctx, mongoClient, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		router:                router,
		mongoClient:           mongoClient,
		AttributeStore:        attributeStore,
		PartyStore:            partyStore,
		PartyTypeStore:        partyTypeStore,
		RelationshipStore:     relationshipStore,
		RelationshipTypeStore: relationshipTypeStore,
		IndividualStore:       NewIndividualStore(mongoClient, o.MongoDatabase),
		StaffStore:            NewStaffStore(relationshipStore),
		TeamStore:             NewTeamStore(partyStore),
		OrganizationStore:     NewOrganizationStore(partyStore),
		MembershipStore:       NewMembershipStore(relationshipStore),
	}

	router.Use(srv.WithAuth(ctx))

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

	router.Path("/apis/iam/v1/organizations").Methods("GET").HandlerFunc(srv.ListOrganizations)
	router.Path("/apis/iam/v1/organizations").Methods("POST").HandlerFunc(srv.PostOrganization)
	router.Path("/apis/iam/v1/organizations/{id}").Methods("GET").HandlerFunc(srv.GetOrganization)
	router.Path("/apis/iam/v1/organizations/{id}").Methods("PUT").HandlerFunc(srv.PutOrganization)

	router.Path("/apis/iam/v1/parties").Methods("GET").HandlerFunc(srv.ListParties)
	router.Path("/apis/iam/v1/parties").Methods("POST").HandlerFunc(srv.PostParty)
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

	router.Path("/apis/iam/v1/staff").Methods("GET").HandlerFunc(srv.ListStaff)
	router.Path("/apis/iam/v1/staff").Methods("POST").HandlerFunc(srv.PostStaff)
	router.Path("/apis/iam/v1/staff/{id}").Methods("GET").HandlerFunc(srv.GetStaff)

	router.Path("/apis/iam/v1/teams").Methods("GET").HandlerFunc(srv.ListTeams)
	router.Path("/apis/iam/v1/teams").Methods("POST").HandlerFunc(srv.PostTeam)
	router.Path("/apis/iam/v1/teams/{id}").Methods("GET").HandlerFunc(srv.GetTeam)
	router.Path("/apis/iam/v1/teams/{id}").Methods("PUT").HandlerFunc(srv.PutTeam)

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
