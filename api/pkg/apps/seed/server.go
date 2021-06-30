package seed

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/apps/login"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"strings"
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
	router      *mux.Router
	mongoClient *mongo.Client
}

func (s Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
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

	srv := &Server{
		router:      router,
		mongoClient: mongoClient,
	}

	router.Path("/seed").Methods("POST").HandlerFunc(srv.SeedHandler)

	return srv, nil

}

func (s *Server) SeedHandler(w http.ResponseWriter, req *http.Request) {
	databaseName := "core" // TODO not sure where to best record this dynamically
	ctx := req.Context()

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var body map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	action := body["action"]

	switch action {
	case "RESET":
		if err := Clear(ctx, databaseName, s.mongoClient); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = Seed(ctx, databaseName, s.mongoClient)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "SEED":
		if err := Seed(ctx, databaseName, s.mongoClient); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Clear(ctx context.Context, databaseName string, mongoClient *mongo.Client) error {
	collectionNames := []string{"parties", "credentials", "relationships", "caseTypes", "cases"}
	for _, collectionName := range collectionNames {
		_, err := mongoClient.Database(databaseName).Collection(collectionName).DeleteMany(ctx, bson.D{})
		if err != nil {
			return err
		}
	}
	return nil
}

func Seed(ctx context.Context, databaseName string, mongoClient *mongo.Client) error {

	for _, obj := range individuals {
		if err := seedMongo(ctx, mongoClient, databaseName, "parties", bson.M{"id": obj.ID}, obj.Party); err != nil {
			return err
		}
		if obj.HasPartyType(iam.StaffPartyType.ID) {
			hash, err := login.HashAndSalt(bcrypt.MinCost, []byte("password"))
			if err != nil {
				return err
			}
			if _, err := mongoClient.Database(databaseName).Collection("credentials").UpdateOne(ctx,
				bson.M{
					"partyId": obj.ID,
				},
				bson.M{
					"$set": bson.M{
						"partyId": obj.ID,
						"hash":    hash,
					},
				},
				options.Update().SetUpsert(true)); err != nil {
				return err
			}
		}
	}

	for _, obj := range teams {
		if err := seedMongo(ctx, mongoClient, databaseName, "parties", bson.M{"id": obj.ID}, iam.MapTeamToParty(&obj)); err != nil {
			return err
		}
	}

	for _, obj := range relationships {
		if err := seedMongo(ctx, mongoClient, databaseName, "relationships", bson.M{"id": obj.ID}, obj); err != nil {
			return err
		}
	}

	for _, obj := range memberships {
		if err := seedMongo(ctx, mongoClient, databaseName, "relationships", bson.M{"id": obj.ID}, iam.MapMembershipToRelationship(&obj)); err != nil {
			return err
		}
	}

	for _, obj := range caseTypes {
		if err := seedMongo(ctx, mongoClient, databaseName, "caseTypes", bson.M{"id": obj.ID}, obj); err != nil {
			return err
		}
	}

	for _, obj := range cases {
		if err := seedMongo(ctx, mongoClient, databaseName, "cases", bson.M{"id": obj.ID}, obj); err != nil {
			return err
		}
	}

	return nil
}

func caseType(id, name, partyTypeID, teamID string) cms.CaseType {
	ct := cms.CaseType{
		ID:          id,
		Name:        name,
		PartyTypeID: partyTypeID,
		TeamID:      teamID,
	}
	caseTypes = append(caseTypes, ct)
	return ct
}

func team(id, name string) iam.Team {
	t := iam.Team{
		ID:   id,
		Name: name,
	}
	teams = append(teams, t)
	return t
}

func individual(id, firstName, lastName string) iam.Individual {
	var i = iam.Individual{
		Party: &iam.Party{
			ID: id,
			PartyTypeIDs: []string{
				iam.IndividualPartyType.ID,
			},
			Attributes: map[string][]string{
				iam.FirstNameAttribute.ID: {firstName},
				iam.LastNameAttribute.ID:  {lastName},
				iam.EMailAttribute.ID:     {strings.ToLower(firstName) + "." + strings.ToLower(lastName) + "@email.com"},
			},
		},
	}
	individuals = append(individuals, i)
	return i
}

func staff(individual iam.Individual) iam.Individual {
	individual.AddPartyType(iam.StaffPartyType.ID)
	return individual
}

func membership(id string, individual iam.Individual, team iam.Team) iam.Membership {
	m := iam.Membership{
		ID:           id,
		TeamID:       team.ID,
		IndividualID: individual.ID,
	}
	memberships = append(memberships, m)
	return m
}

func kase(id, caseTypeID, partyID, teamID, description string, done bool) cms.Case {
	k := cms.Case{
		ID:          id,
		CaseTypeID:  caseTypeID,
		PartyID:     partyID,
		TeamID:      teamID,
		Description: description,
		Done:        done,
	}
	cases = append(cases, k)
	return k
}

func seedMongo(ctx context.Context, mongoClient *mongo.Client, databaseName, collectionName string, filter interface{}, document interface{}) error {
	logrus.Infof("seeding collection %s.%s with object: %#v", databaseName, collectionName, document)
	collection := mongoClient.Database(databaseName).Collection(collectionName)
	if _, err := collection.InsertOne(ctx, document); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			return err
		}
		if _, err := collection.ReplaceOne(ctx, filter, document); err != nil {
			return err
		}
	}
	return nil
}
