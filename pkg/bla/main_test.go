package bla

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/bla/client"
	"github.com/nrc-no/core/pkg/bla/types"
	"github.com/nrc-no/core/pkg/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type Suite struct {
	suite.Suite
	server *Server
	done   chan struct{}
	client client.Client
}

func (s *Suite) SetupSuite() {
	host := "localhost"
	port := 5435
	user := "postgres"
	password := "postgres"
	dbname := "core"
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	s.done = make(chan struct{}, 1)
	s.server = NewServer(psqlconn)

	if err := s.server.Start(s.done, "127.0.0.1:"); err != nil {
		s.T().Fatal(err)
	}

	s.client = client.NewClientFromConfig(rest.Config{
		Scheme:     "http",
		Host:       s.server.Address(),
		HTTPClient: http.DefaultClient,
	})

}

func (s *Suite) TestCreateDatabase() {
	ctx := context.Background()
	var db types.Database
	if err := s.client.CreateDatabase(ctx, &types.Database{
		Name: "testdb",
	}, &db); !assert.NoError(s.T(), err) {
		return
	}
	_ = s.client.DeleteDatabase(ctx, db.ID)
}

func (s *Suite) TestCreateForm() {
	ctx := context.Background()

	db, doneDb, err := s.WithDatabase(ctx)
	if !assert.NoError(s.T(), err) {
		return
	}
	defer doneDb()

	folder, doneFolder, err := s.WithFolder(ctx, db.ID)
	if !assert.NoError(s.T(), err) {
		return
	}
	defer doneFolder()

	var fd types.FormDefinition
	if err := s.client.CreateForm(ctx, &types.FormDefinition{
		ID:         "",
		DatabaseID: db.ID,
		Name:       "testform",
		FolderID:   folder.ID,
		Fields: []*types.FieldDefinition{
			{
				Name: "testfield",
				FieldType: types.FieldType{
					Text: &types.FieldTypeText{},
				},
			},
		},
	}, &fd); !assert.NoError(s.T(), err) {
		return
	}
	assert.NotEmpty(s.T(), fd.ID)
	assert.Equal(s.T(), "testform", fd.Name)
	assert.Equal(s.T(), folder.ID, fd.FolderID)
	assert.Equal(s.T(), db.ID, db.ID)
	if !assert.Equal(s.T(), 1, len(fd.Fields)) {
		return
	}
	assert.NotEmpty(s.T(), fd.Fields[0].ID)
	assert.Equal(s.T(), "testfield", fd.Fields[0].Name)
	assert.NotNil(s.T(), fd.Fields[0].FieldType.Text)
	assert.Nil(s.T(), fd.Fields[0].FieldType.Reference)
	assert.Nil(s.T(), fd.Fields[0].FieldType.SubForm)
}

func (s *Suite) TearDownSuite() {
	s.done <- struct{}{}
}

func (s *Suite) WithFolder(ctx context.Context, dbId string) (*types.Folder, func(), error) {
	var folder types.Folder
	if err := s.client.CreateFolder(ctx, &types.Folder{
		DatabaseID: dbId,
		Name:       "testdb",
	}, &folder); err != nil {
		return nil, nil, err
	}
	return &folder, func() {
		s.client.DeleteFolder(ctx, folder.ID)
	}, nil
}

func (s *Suite) WithDatabase(ctx context.Context) (*types.Database, func(), error) {
	var db types.Database
	if err := s.client.CreateDatabase(ctx, &types.Database{
		Name: "testdb",
	}, &db); err != nil {
		return nil, nil, err
	}
	return &db, func() {
		s.client.DeleteDatabase(ctx, db.ID)
	}, nil
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
