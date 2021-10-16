package form

import (
	"context"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

type Suite struct {
	suite.Suite
	done         chan struct{}
	server       *Server
	cli          Interface
	uidFn        func() string
	uidGenerator *utils.DelegateUidGenerator
	mongoCli     *mongo.Client
	databaseName string
}

func (s *Suite) NextUID(id string) {
	s.uidGenerator.SetGeneratorFn(func() string {
		return id
	})
}

func (s *Suite) SetupSuite() {
	s.done = make(chan struct{})
	s.uidFn = func() string {
		return "id"
	}
	s.uidGenerator = utils.NewDelegateUIDGenerator(s.uidFn)
	var err error
	s.mongoCli, err = mongo.Connect(context.Background())
	if err != nil {
		s.T().Fatal(err)
		return
	}
	s.databaseName = "test"
	s.mongoCli.Database(s.databaseName).Collection(collFormDefinitions).DeleteMany(context.Background(), bson.M{})
	s.server = NewServer(
		s.uidGenerator, func() (*mongo.Client, error) {
			return s.mongoCli, nil
		}, s.databaseName)
	if err := s.server.Start(s.done); err != nil {
		s.T().Fatal(err)
	}
	s.cli = s.server.NewClient()
}

func (s *Suite) TearDownSuite() {
	s.done <- struct{}{}
}

func (s *Suite) TestGetFormDefinitionShouldReturn() {
	ctx := context.Background()

	s.NextUID("TestGetFormDefinitionShouldReturn")
	def := &Definition{Name: "hello"}
	var err error
	if def, err = s.cli.Create(ctx, def, CreateOptions{}); !assert.NoError(s.T(), err) {
		return
	}

	res, err := s.cli.Get(ctx, def.ID, GetOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), def.ID, res.ID)
	assert.Equal(s.T(), def.Version, res.Version)

}

func (s *Suite) TestGetFormDefinitionShouldReturnLatest() {
	ctx := context.Background()

	s.NextUID("TestGetFormDefinitionShouldReturnLatest-1")
	def := &Definition{Name: "hello"}
	first, err := s.cli.Create(ctx, def, CreateOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}

	s.NextUID("TestGetFormDefinitionShouldReturnLatest-2")
	last, err := s.cli.Put(ctx, first, PutOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}

	got, err := s.cli.Get(ctx, last.ID, GetOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), last.ID, got.ID)
	assert.Equal(s.T(), last.Version, got.Version)

}

func (s *Suite) TestGetFormDefinitionWithVersionShouldReturnVersion() {
	ctx := context.Background()

	s.NextUID("TestGetFormDefinitionShouldReturnLatest-1")
	def := &Definition{Name: "hello"}
	first, err := s.cli.Create(ctx, def, CreateOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}

	s.NextUID("TestGetFormDefinitionShouldReturnLatest-2")
	last, err := s.cli.Put(ctx, first, PutOptions{})
	if !assert.NoError(s.T(), err) {
		return
	}

	got, err := s.cli.Get(ctx, last.ID, GetOptions{
		Version: first.Version,
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), first.ID, got.ID)
	assert.Equal(s.T(), first.Version, got.Version)

}

func (s *Suite) TestCreateFormDefinition() {

	ctx := context.Background()

	type args struct {
		preCreate              *Definition
		name                   string
		in                     Definition
		expectErr              bool
		expectName             string
		expectVersion          string
		expectIsCurrentVersion bool
		nextUid                string
	}
	tcs := []args{
		{
			name:                   "putCurrentVersion",
			preCreate:              &Definition{ID: "create-fd-test-1", Name: "hello"},
			nextUid:                "2",
			in:                     Definition{ID: "create-fd-test-1", Name: "hello", IsCurrentVersion: true},
			expectName:             "hello",
			expectVersion:          "2",
			expectIsCurrentVersion: true,
		}, {
			name:                   "putNonCurrentVersion",
			preCreate:              &Definition{ID: "create-fd-test-2", Name: "hello"},
			nextUid:                "2",
			in:                     Definition{ID: "create-fd-test-2", Name: "hello", IsCurrentVersion: false},
			expectName:             "hello",
			expectVersion:          "2",
			expectIsCurrentVersion: false,
		},
	}

	for _, tc := range tcs {
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.preCreate != nil {
				s.NextUID(tc.preCreate.ID)
				_, err := s.cli.Create(ctx, tc.preCreate, CreateOptions{})
				if !assert.NoError(t, err) {
					return
				}
			}

			s.NextUID(tc.nextUid)

			out, err := s.cli.Put(ctx, &tc.in, PutOptions{})
			if tc.expectErr {
				if !assert.Error(t, err) {
					return
				}
			} else {
				if !assert.NoError(t, err) {
					return
				}
				assert.Equal(t, tc.expectName, out.Name)
				assert.Equal(t, tc.expectVersion, out.Version)
				assert.Equal(t, tc.expectIsCurrentVersion, out.IsCurrentVersion)
			}
		})
	}

}
func (s *Suite) TestUpdateFormDefinition() {

	ctx := context.Background()

	type args struct {
		name                   string
		in                     Definition
		expectErr              bool
		expectName             string
		expectVersion          string
		expectIsCurrentVersion bool
		expectNotFound         bool
	}
	tcs := []args{
		{
			name:                   "updateFormDefinition",
			in:                     Definition{Name: "hello"},
			expectName:             "hello",
			expectVersion:          "updateFormDefinition",
			expectIsCurrentVersion: true,
			expectErr:              false,
			expectNotFound:         false,
		},
	}

	for _, tc := range tcs {
		s.T().Run(tc.name, func(t *testing.T) {
			s.NextUID(tc.name)
			out, err := s.cli.Create(ctx, &tc.in, CreateOptions{})
			if tc.expectErr {
				if !assert.Error(t, err) {
					return
				}
			} else {
				if !assert.NoError(t, err) {
					return
				}
				assert.Equal(t, tc.expectName, out.Name)
				assert.Equal(t, tc.expectVersion, out.Version)
				assert.Equal(t, tc.expectIsCurrentVersion, out.IsCurrentVersion)
			}
		})
	}

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
