package iam

import (
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
}

func NewServer(o *ServerOptions) (*Server, error) {
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

	router := mux.NewRouter()

	relationshipStore := NewRelationshipStore(mongoClient, o.MongoDatabase)
	partyStore := NewPartyStore(mongoClient, o.MongoDatabase)

	srv := &Server{
		router:                router,
		AttributeStore:        NewAttributeStore(mongoClient, o.MongoDatabase),
		PartyStore:            partyStore,
		PartyTypeStore:        NewPartyTypeStore(mongoClient, o.MongoDatabase),
		RelationshipStore:     relationshipStore,
		RelationshipTypeStore: NewRelationshipTypeStore(mongoClient, o.MongoDatabase),
		IndividualStore:       NewIndividualStore(mongoClient, o.MongoDatabase),
		StaffStore:            NewStaffStore(relationshipStore),
		TeamStore:             NewTeamStore(partyStore),
		OrganizationStore:     NewOrganizationStore(partyStore),
		MembershipStore:       NewMembershipStore(relationshipStore),
	}

	attributes := router.Path("/apis/v1/attributes")
	attributes.Path("").Methods("GET").HandlerFunc(srv.ListAttributes)
	attributes.Path("").Methods("POST").HandlerFunc(srv.PostAttribute)
	attributes.Path("/{id}").Methods("GET").HandlerFunc(srv.GetAttribute)
	attributes.Path("/{id}").Methods("PUT").HandlerFunc(srv.PutAttribute)

	individuals := router.Path("/apis/v1/individuals")
	individuals.Path("").Methods("GET").HandlerFunc(srv.ListIndividuals)
	individuals.Path("").Methods("POST").HandlerFunc(srv.PostIndividual)
	individuals.Path("/{id}").Methods("GET").HandlerFunc(srv.GetIndividual)
	individuals.Path("/{id}").Methods("PUT").HandlerFunc(srv.PutIndividual)

	memberships := router.Path("/apis/v1/memberships")
	memberships.Path("").Methods("GET").HandlerFunc(srv.ListMemberships)
	memberships.Path("").Methods("POST").HandlerFunc(srv.PostMembership)
	memberships.Path("/{id}").Methods("GET").HandlerFunc(srv.GetMembership)

	organizations := router.Path("/apis/v1/organizations")
	organizations.Path("").Methods("GET").HandlerFunc(srv.ListOrganizations)
	organizations.Path("").Methods("POST").HandlerFunc(srv.PostOrganization)
	organizations.Path("/{id}").Methods("GET").HandlerFunc(srv.GetOrganization)
	organizations.Path("/{id}").Methods("PUT").HandlerFunc(srv.PutOrganization)

	parties := router.Path("/apis/v1/parties")
	parties.Path("").Methods("GET").HandlerFunc(srv.ListParties)
	parties.Path("").Methods("POST").HandlerFunc(srv.PostParty)
	parties.Path("/{id}").Methods("GET").HandlerFunc(srv.GetParty)
	parties.Path("/{id}").Methods("PUT").HandlerFunc(srv.PutParty)

	partyTypes := router.Path("/apis/v1/partytypes")
	partyTypes.Path("").Methods("GET").HandlerFunc(srv.ListPartyTypes)
	partyTypes.Path("").Methods("POST").HandlerFunc(srv.PostPartyType)
	partyTypes.Path("/{id}").Methods("GET").HandlerFunc(srv.GetPartyType)
	partyTypes.Path("/{id}").Methods("PUT").HandlerFunc(srv.PutPartyType)

	relationships := router.Path("/apis/v1/relationships")
	relationships.Path("").Methods("GET").HandlerFunc(srv.ListRelationships)
	relationships.Path("").Methods("POST").HandlerFunc(srv.PostRelationship)
	relationships.Path("/{id}").Methods("GET").HandlerFunc(srv.GetRelationship)
	relationships.Path("/{id}").Methods("PUT").HandlerFunc(srv.PutRelationship)
	relationships.Path("/{id}").Methods("DELETE").HandlerFunc(srv.DeleteRelationship)

	relationshipTypes := router.Path("/apis/v1/relationshiptypes")
	relationshipTypes.Path("").Methods("GET").HandlerFunc(srv.ListRelationshipTypes)
	relationshipTypes.Path("").Methods("POST").HandlerFunc(srv.PostRelationshipType)
	relationshipTypes.Path("/{id}").Methods("GET").HandlerFunc(srv.GetRelationshipType)
	relationshipTypes.Path("/{id}").Methods("PUT").HandlerFunc(srv.PutRelationshipType)

	staff := router.Path("/apis/v1/staff")
	staff.Path("").Methods("GET").HandlerFunc(srv.ListStaff)
	staff.Path("").Methods("POST").HandlerFunc(srv.PostStaff)
	staff.Path("/{id}").Methods("GET").HandlerFunc(srv.GetStaff)

	teams := router.Path("/apis/v1/teams")
	teams.Path("").Methods("GET").HandlerFunc(srv.ListTeams)
	teams.Path("").Methods("POST").HandlerFunc(srv.PostTeam)
	teams.Path("/{id}").Methods("GET").HandlerFunc(srv.GetTeam)
	teams.Path("/{id}").Methods("PUT").HandlerFunc(srv.PutTeam)

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
