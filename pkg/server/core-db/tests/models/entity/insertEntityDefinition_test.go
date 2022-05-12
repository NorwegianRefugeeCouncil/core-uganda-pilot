package tests

import (
	"context"

	"github.com/lib/pq"

	"github.com/nrc-no/core/pkg/server/core-db/types"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type Result struct {
	ID               string         `gorm:"column:id"`
	Name             string         `gorm:"column:name"`
	Description      string         `gorm:"column:description"`
	ConstraintCustom pq.StringArray `gorm:"column:constraint_custom"`
}

func (s *Suite) TestInsertEntityDefinitionSuccessSingle() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := s.DBFactory.Get()

	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}

	e := types.EntityDefinition{
		ID:          uuid.NewV4().String(),
		Name:        "test",
		Description: "test",
		Constraints: types.EntityConstraints{
			Custom: []string{
				"test",
			},
		},
	}

	createdEntity, err := s.entityModel.InsertEntityDefinition(ctx, db, e)

	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}

	assert.Equal(s.T(), e.ID, createdEntity.ID)
	assert.Equal(s.T(), e.Name, createdEntity.Name)
	assert.Equal(s.T(), e.Description, createdEntity.Description)
	assert.Equal(s.T(), e.Constraints.Custom, createdEntity.Constraints.Custom)

	var result []*Result
	db.Raw("SELECT * FROM entity_definition").Scan(&result)

	assert.Equal(s.T(), 1, len(result))

	assert.Equal(s.T(), e.ID, result[0].ID)
	assert.Equal(s.T(), e.Name, result[0].Name)
	assert.Equal(s.T(), e.Description, result[0].Description)
	assert.Equal(s.T(), e.Constraints.Custom, []string{result[0].ConstraintCustom[0]})
}

func (s *Suite) TestInsertEntityDefinitionSuccessMultiple() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := s.DBFactory.Get()

	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}

	e1 := types.EntityDefinition{

		ID:          uuid.NewV4().String(),
		Name:        "test",
		Description: "test",
		Constraints: types.EntityConstraints{
			Custom: []string{
				"test",
			},
		},
	}

	e2 := types.EntityDefinition{
		ID:          uuid.NewV4().String(),
		Name:        "test",
		Description: "test",
		Constraints: types.EntityConstraints{
			Custom: []string{
				"test",
			},
		},
	}

	createdEntity1, err := s.entityModel.InsertEntityDefinition(ctx, db, e1)

	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
		return
	}

	createdEntity2, err := s.entityModel.InsertEntityDefinition(ctx, db, e2)

	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
		return
	}

	assert.Equal(s.T(), e1.ID, createdEntity1.ID)
	assert.Equal(s.T(), e1.Name, createdEntity1.Name)
	assert.Equal(s.T(), e1.Description, createdEntity1.Description)
	assert.Equal(s.T(), e1.Constraints.Custom, createdEntity1.Constraints.Custom)

	assert.Equal(s.T(), e2.ID, createdEntity2.ID)
	assert.Equal(s.T(), e2.Name, createdEntity2.Name)
	assert.Equal(s.T(), e2.Description, createdEntity2.Description)
	assert.Equal(s.T(), e2.Constraints.Custom, createdEntity2.Constraints.Custom)

	var result []*Result
	db.Raw("SELECT * FROM entity_definition").Scan(&result)

	assert.Equal(s.T(), 2, len(result))

	assert.Equal(s.T(), e1.ID, result[0].ID)
	assert.Equal(s.T(), e1.Name, result[0].Name)
	assert.Equal(s.T(), e1.Description, result[0].Description)
	assert.Equal(s.T(), e1.Constraints.Custom, []string{result[0].ConstraintCustom[0]})

	assert.Equal(s.T(), e2.ID, result[1].ID)
	assert.Equal(s.T(), e2.Name, result[1].Name)
	assert.Equal(s.T(), e2.Description, result[1].Description)
	assert.Equal(s.T(), e2.Constraints.Custom, []string{result[1].ConstraintCustom[0]})
}

func (s *Suite) TestInsertEntityDefinitionSuccessMissingConstraint() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := s.DBFactory.Get()

	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}

	e := types.EntityDefinition{
		ID:          uuid.NewV4().String(),
		Name:        "test",
		Description: "test",
	}

	createdEntity, err := s.entityModel.InsertEntityDefinition(ctx, db, e)

	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}

	assert.Equal(s.T(), e.ID, createdEntity.ID)
	assert.Equal(s.T(), e.Name, createdEntity.Name)
	assert.Equal(s.T(), e.Description, createdEntity.Description)
	assert.Equal(s.T(), e.Constraints.Custom, createdEntity.Constraints.Custom)

	var result []*Result
	db.Raw("SELECT * FROM entity_definition").Scan(&result)

	assert.Equal(s.T(), 1, len(result))

	assert.Equal(s.T(), e.ID, result[0].ID)
	assert.Equal(s.T(), e.Name, result[0].Name)
	assert.Equal(s.T(), e.Description, result[0].Description)
	assert.Equal(s.T(), len(e.Constraints.Custom), len(result[0].ConstraintCustom))
}

func (s *Suite) TestInsertEntityDefinitionFailInvalidUUID() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := s.DBFactory.Get()

	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}

	e := types.EntityDefinition{
		ID:          "not-a-uuid",
		Name:        "test",
		Description: "test",
		Constraints: types.EntityConstraints{
			Custom: []string{
				"test",
			},
		},
	}

	_, err = s.entityModel.InsertEntityDefinition(ctx, db, e)

	assert.Error(s.T(), err)

	var result []*Result
	db.Raw("SELECT * FROM entity_definition").Scan(&result)

	assert.Equal(s.T(), 0, len(result))
}
