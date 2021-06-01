import {
  createAsyncThunk,
  createEntityAdapter,
  createSelector,
  createSlice,
  EntityState,
  PayloadAction
} from '@reduxjs/toolkit';
import * as uuid from 'uuid';

import { HasObjectMeta, isInCondition, Operation, setCondition, TypeMeta } from '../api';

export type Entity = TypeMeta & HasObjectMeta

export const getEntityId = (entity: Entity): string => {
  return [
    entity.apiVersion,
    entity.kind,
    entity.metadata.name
  ].join('/');
};

const adapter = createEntityAdapter<Entity>({
  selectId: (a) => {
    return getEntityId(a);
  },
  sortComparer: (a, b) => {
    return getEntityId(a).localeCompare(getEntityId(b));
  }
});

export type State = {
  // represents the operations done by the
  // user who are not yet synced to the server
  operations: Operation[],
  // index of the last committed operation
  firstUncommittedIndex: number
} & EntityState<Entity>

export type StateSlice = {
  'entities': State
}

export const INITIAL_STATE: State = {
  ...adapter.getInitialState(),
  firstUncommittedIndex: 0,
  operations: []
};

export const submit = createAsyncThunk<any,
  void,
  {
    state: StateSlice,
  }>(
  'entities/submitOperation',
  async (arg, { getState }): Promise<void> => {
    await new Promise(resolve => {
      setTimeout(() => {
        resolve(null);
      }, 300);
    });
  },
  {
    condition: (arg, { getState }) => {
      const state = getState();
      if (state.entities.firstUncommittedIndex === 0) {
        return false;
      }
      for (let operation of state.entities.operations) {
        if (isInCondition(operation, 'OperationPending', 'True')) {
          return false;
        }
      }
      const first = selectNextQueuedOperation(getState());
      if (isInCondition(first, 'OperationPending', 'True')) {
        return false;
      }
      if (isInCondition(first, 'OperationAccepted', 'True')) {
        return false;
      }
      return true;
    },
    dispatchConditionRejection: false
  }
);

export const entitiesSlice = createSlice({
  name: 'entities',
  initialState: INITIAL_STATE,
  reducers: {
    entityAdded: adapter.addOne,
    entitiesAdded: adapter.addMany,
    entityRemoved: adapter.removeOne,
    entitiesRemoved: adapter.removeAll,
    entityUpdated: adapter.updateOne,
    entitiesUpdated: adapter.updateMany,
    entityUpserted: adapter.upsertOne,
    entitiesUpserted: adapter.upsertMany,
    entitiesSet: adapter.setAll,
    commit: (state) => {
      const id = uuid.v4();
      for (let i = state.firstUncommittedIndex; i < state.operations.length; i++) {
        state.operations[i].spec.batchId = id;
      }
      state.firstUncommittedIndex = state.operations.length;
    },
    rollback: (state) => {
      state.operations.splice(state.firstUncommittedIndex);
    },
    push: (state, action: PayloadAction<{ operation: Operation }>) => {
      state.operations.push(action.payload.operation);
    },
    clearCommitted: (state) => {
      state.operations.splice(0, state.firstUncommittedIndex);
      state.firstUncommittedIndex = 0;
    },
    clearUncommitted: (state) => {
      state.operations.splice(state.firstUncommittedIndex);
    },
    acknowledgeAcceptedOperations: (state) => {
      let count = 0;
      for (let i = state.firstUncommittedIndex - 1; i >= 0; i--) {
        if (isInCondition(state.operations[i], 'OperationAccepted', 'True')) {
          state.operations.splice(i, 1);
          count++;
        }
      }
      state.firstUncommittedIndex -= count;
    }
  }, extraReducers: builder => {
    builder.addCase(submit.fulfilled, (state, action) => {

      const batch = selectNextQueuedOperationBatch({ entities: state });
      for (let operation of batch) {
        setCondition(operation, 'OperationAccepted', 'True');
        setCondition(operation, 'OperationPending', 'False');
        adapter.addOne(state, operation.spec.entity);
      }

    });
    builder.addCase(submit.pending, (state, action) => {
      const batch = selectNextQueuedOperationBatch({ entities: state });
      for (let operation of batch) {
        setCondition(operation, 'OperationPending', 'True');
      }
    });
    builder.addCase(submit.rejected, (state, action) => {
      const batch = selectNextQueuedOperationBatch({ entities: state });
      for (let operation of batch) {
        setCondition(operation, 'OperationAccepted', 'False');
        setCondition(operation, 'OperationPending', 'False');
      }
    });
  }
});

export const entitiesReducer = entitiesSlice.reducer;

export const {
  rollback,
  push,
  commit,
  clearCommitted,
  clearUncommitted,
  acknowledgeAcceptedOperations,
  entitiesAdded,
  entitiesRemoved,
  entitiesUpdated,
  entitiesUpserted,
  entityAdded,
  entityRemoved,
  entityUpdated,
  entitiesSet,
  entityUpserted
} = entitiesSlice.actions;

// selector for the the state slice
export const selectState = createSelector<StateSlice, State, State>(state => state.entities, res => res);

export const selectEntities = createSelector(selectState, state => {
  return state.ids.map(id => state.entities[id]);
});

// selector for the operations
export const selectOperations = createSelector(selectState, res => res.operations);

// selector for the operation commit index
export const selectCommittedSelector = createSelector(selectState, res => res.firstUncommittedIndex);

// selector for the committed operations
export const selectCommittedOperations = createSelector(
  selectOperations,
  selectCommittedSelector,
  (operations, firstUncommittedIndex) => {
    const result: Operation[] = [];
    for (let i = 0; i < firstUncommittedIndex; i++) {
      result.push(operations[i]);
    }
    return result;
  }
);

export const selectNextQueuedOperationBatch = createSelector(
  selectCommittedOperations,
  res => {
    let batchId = '';
    let ops: Operation[] = [];
    for (let op of res) {
      if (isInCondition(op, 'OperationAccepted', 'True')) {
        continue;
      }
      if (!batchId) {
        batchId = op.spec.batchId;
      }
      if (op.spec.batchId === batchId) {
        ops.push(op);
      } else {
        break;
      }
    }
    return ops;
  }
);

export const selectOperationBatches = createSelector(
  selectOperations,
  res => {
    let currentBatchId = '';
    let ops: Operation[][] = [];
    for (let op of res) {
      if (!currentBatchId) {
        currentBatchId = op.spec.batchId;
      }
      if (op.spec.batchId === currentBatchId) {
        if (ops.length === 0){
          ops.push([])
        }
        ops[ops.length - 1].push(op);
      } else {
        ops.push([op]);
        currentBatchId = op.spec.batchId;
      }
    }
    return ops;
  }
);

export const selectNextQueuedOperation = createSelector(
  selectCommittedOperations,
  res => {
    for (let op of res) {
      if (isInCondition(op, 'OperationAccepted', 'True')) {
        continue;
      }
      return op;
    }
  }
);

// selector for the uncommitted operations
export const selectUncommittedOperations = createSelector(
  selectOperations,
  selectCommittedSelector,
  (operations, firstUncommittedIndex) => {
    const result: Operation[] = [];
    for (let i = firstUncommittedIndex; i < operations.length; i++) {
      result.push(operations[i]);
    }
    return result;
  }
);

// selector for the committed state
export const selectCommittedState = createSelector(
  selectState,
  selectCommittedOperations,
  (state, operations) => {
    for (let operation of operations) {
      state = applyOperationToState(state, operation) as State;
    }
    return state;
  }
);

export const selectCommittedEntities = createSelector(selectCommittedState, state => {
  return state.ids.map(id => state.entities[id]);
});

// selector for the uncommitted state
export const selectUncommittedState = createSelector(
  selectCommittedState,
  selectUncommittedOperations,
  (state, operations) => {
    for (let operation of operations) {
      state = applyOperationToState(state, operation) as State;
    }
    return state;
  }
);

export const selectUncommittedEntities = createSelector(selectUncommittedState, state => {
  return state.ids.map(id => state.entities[id]);
});

export const selectEntity = (id: string) => {
  return createSelector(
    selectState,
    res => res.entities[id]
  );
};

export const selectCommittedEntity = (id: string) => {
  return createSelector(
    selectEntity(id),
    selectCommittedOperations,
    (entity, operations) => {
      applyOperationsToEntity(entity, operations);
      return entity;
    }
  );
};

export const selectUncommittedEntity = (id: string) => {
  return createSelector(
    selectEntity(id),
    selectUncommittedOperations,
    (entity, operations) => {
      entity = applyOperationsToEntity(entity, operations);
      return entity;
    }
  );
};

const applyOperationsToEntity = (entity: Entity, operations: Operation[]): Entity => {
  const id = getEntityId(entity);
  let entityState: EntityState<Entity> = {
    ids: [id],
    entities: {
      [id]: entity
    }
  };
  if (operations.length === 0) {
    return entity;
  }
  for (let operation of operations) {
    const operationEntityId = getEntityId(operation.spec.entity);
    if (operationEntityId !== id) {
      continue;
    }
    entityState = applyOperationToState(entityState, operation);
  }
  return entityState.entities[id];
};

const applyOperationToState = (state: EntityState<Entity>, operation: Operation): EntityState<Entity> => {
  if (operation.spec.operationType === 'delete') {
    return adapter.removeOne(state, getEntityId(operation.spec.entity));
  }
  if (operation.spec.operationType === 'create') {
    return adapter.addOne(state, operation.spec.entity);
  }
  if (operation.spec.operationType === 'update') {
    const id = getEntityId(operation.spec.entity);
    let changes = operation.spec.entity;
    return adapter.updateOne(state, { id, changes: changes });
  }
};
