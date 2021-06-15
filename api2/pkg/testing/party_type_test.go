package testing

import (
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestPartyTypeCrud() {
	name := uuid.NewV4().String()
	created, err := s.server.PartyTypeClient.Create(s.ctx, &partytypes.PartyType{
		Name:      name,
		IsBuiltIn: false,
	})
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), name, created.Name)
	assert.NotEmpty(s.T(), created.ID)
	assert.False(s.T(), created.IsBuiltIn)

	get, err := s.server.PartyTypeClient.Get(s.ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.Equal(s.T(), get, created) {
		return
	}

	list, err := s.server.PartyTypeClient.List(s.ctx)
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Contains(s.T(), list.Items, get)

}
