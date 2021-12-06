package sqlmanager

import (
	"errors"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/schema"
)

type sqlState struct {
	Tables schema.SQLTables
}

type Interface interface {
	PutForms(formDefinitions *types.FormDefinitionList) (Interface, error)
	GetStatements() []schema.DDL
}

type sqlManager struct {
	State      sqlState
	Statements []schema.DDL
}

func New() Interface {
	return sqlManager{}
}

func (s sqlManager) GetStatements() []schema.DDL {
	return s.Statements
}

func (s sqlManager) PutForms(formDefinitions *types.FormDefinitionList) (Interface, error) {
	actions := sqlActions{}
	for _, item := range formDefinitions.Items {
		formActions, err := getSQLActionsForForm(item)
		if err != nil {
			return nil, err
		}
		actions = append(actions, formActions...)
	}
	return s.handleActions(actions)
}

func (s sqlManager) handleActions(actions sqlActions) (sqlManager, error) {
	walk := s
	var err error
	for _, action := range actions {
		walk, err = walk.handleAction(action)
		if err != nil {
			return sqlManager{}, err
		}
	}
	return walk, nil
}

func (s sqlManager) handleAction(action sqlAction) (sqlManager, error) {
	if action.createColumn != nil {
		return s.handleCreateColumn(*action.createColumn)
	}
	if action.createTable != nil {
		return s.handleCreateTable(*action.createTable)
	}
	if action.createUniqueConstraint != nil {
		return s.handleCreateConstraint(*action.createUniqueConstraint)
	}
	return sqlManager{}, errors.New("could not handle action")
}
