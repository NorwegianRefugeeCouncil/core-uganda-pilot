import {
  entityAdded,
  entitiesSlice,
  INITIAL_STATE,
  Entity,
  selectCommittedState,
  selectEntity,
  getEntityId, selectCommittedEntity, selectUncommittedEntity, State, StateSlice, commit, push, selectUncommittedState
} from './redux';
import configureMockStore from 'redux-mock-store';
import { HasObjectMeta, TypeMeta, Operation, OperationType } from '../api';

const mockStore = configureMockStore();

type TestEntity = HasObjectMeta & TypeMeta & {
  spec: {
    property: string
  }
}

type EntityBuilder = (entity: TestEntity) => TestEntity

const withName: (name: string) => EntityBuilder = name => {
  return entity => ({
    ...entity,
    metadata: {
      ...entity.metadata,
      name
    }
  });
};

const withProperty: (value: string) => EntityBuilder = value => {
  return entity => ({
    ...entity,
    spec: {
      ...entity.spec,
      property: value
    }
  });
};

const withBase: (base: TestEntity) => EntityBuilder = value => {
  return () => ({
    kind: value.kind,
    apiVersion: value.apiVersion,
    metadata: {
      ...value.metadata
    },
    spec: {
      ...value.spec
    }
  });
};

const anEntity = (...builders: EntityBuilder[]): TestEntity => {
  let entity: TestEntity = {
    apiVersion: 'test.com/v1',
    kind: 'TestResource',
    metadata: {
      name: ''
    },
    spec: {
      property: ''
    }
  };
  for (let builder of builders) {
    entity = builder(entity);
  }
  return entity;
};

type OperationBuilder = (operation: Operation) => Operation

const withType = (operationType: OperationType): OperationBuilder => operation => {
  return {
    ...operation,
    spec: {
      ...operation.spec,
      operationType
    }
  };
};

const withEntity = (entity: HasObjectMeta & TypeMeta): OperationBuilder => operation => {
  return {
    ...operation,
    spec: {
      ...operation.spec,
      entity
    }
  };
};

const anOperation = (...builders: OperationBuilder[]): Operation => {
  let op: Operation = {
    apiVersion: 'offline.nrc.no/v1',
    kind: 'Operation',
    metadata: {
      name: ''
    },
    spec: {}
  };
  for (let builder of builders) {
    op = builder(op);
  }
  return op;
};

const slice = (state: State): StateSlice => {
  return {
    entities: state
  };
};

const reducer = entitiesSlice.reducer;

const stateWithEntities = (...entities: Entity[]): { initialState: State, stateSlice: StateSlice } => {
  let state = INITIAL_STATE;
  for (let entity of entities) {
    state = reducer(state, entityAdded(entity));
  }
  return {
    initialState: state,
    stateSlice: slice(state)
  };
};

const getProjections = (state: State) => {
  const stateSlice = slice(state);
  const committed = selectCommittedState(stateSlice);
  const uncommitted = selectUncommittedState(stateSlice);
  return {
    stateSlice,
    committedState: committed,
    uncommittedState: uncommitted
  };
};

const getEntityProjections = (state: State, id: string) => {
  const stateSlice = slice(state);
  const actualEntity = selectEntity(id)(stateSlice);
  const committedEntity = selectCommittedEntity(id)(stateSlice);
  const uncommittedEntity = selectUncommittedEntity(id)(stateSlice);
  return {
    actualEntity,
    committedEntity,
    uncommittedEntity
  };
};

describe('redux', () => {

  const reducer = entitiesSlice.reducer;

  it('should add an entity to the state', function() {
    const state = INITIAL_STATE;
    const entity = anEntity(withName('test'));
    const id = getEntityId(entity);
    const actualState = reducer(state, entityAdded(entity));

    const actualSlice = slice(actualState);
    const committedState = selectCommittedState(actualSlice);
    const uncommittedState = selectUncommittedState(actualSlice);

    const actualEntity = selectEntity(id)(actualSlice);
    const committedEntity = selectCommittedEntity(id)(actualSlice);
    const uncommittedEntity = selectUncommittedEntity(id)(actualSlice);

    expect(actualEntity).toEqual(entity);
    expect(committedEntity).toEqual(entity);
    expect(uncommittedEntity).toEqual(entity);
    expect(actualState).toEqual(committedState);
    expect(actualState).toEqual(uncommittedState);

  });

  describe('push operation', () => {

    it('should update the uncommitted state', function() {

      const initialEntity = anEntity(withName('test'), withProperty('A'));
      const updatedEntity = anEntity(withBase(initialEntity), withProperty('B'));
      const entityId = getEntityId(initialEntity);

      const { initialState } = stateWithEntities(initialEntity);

      const operation = anOperation(withEntity(updatedEntity), withType('update'));
      const updatedState = reducer(initialState, push({ operation }));

      let { committedState, uncommittedState } = getProjections(updatedState);
      let { actualEntity, committedEntity, uncommittedEntity } = getEntityProjections(updatedState, entityId);

      expect(updatedState).toEqual(committedState);
      expect(uncommittedState).not.toEqual(committedState);

      expect(actualEntity).toEqual(initialEntity);
      expect(committedEntity).toEqual(initialEntity);
      expect(uncommittedEntity).toEqual(updatedEntity);

      expect(updatedState.firstUncommittedIndex).toBe(0);
      expect(updatedState.operations).toContain(operation);
    });

  });

});
