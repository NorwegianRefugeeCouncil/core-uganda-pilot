import { createAsyncThunk, createEntityAdapter, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { FormDefinition, FormDefinitionList, Status } from '@core/api-client';
import { Api } from '../data/api';

const formDefinitionsAdapter = createEntityAdapter<FormDefinition>({
  selectId: (f) => f.metadata.name,
  sortComparer: (a, b) => a.metadata.name.localeCompare(b.metadata.name)
});

export const listFormDefinitions = createAsyncThunk<FormDefinitionList, void, { rejectValue: Status }>(
  'formDefinitions/list',
  async (_, { rejectWithValue }) => {
    try {
      return await Api.core().v1().formDefinitions().list().toPromise();
    } catch (err) {
      return rejectWithValue(err as Status);
    }
  }
);

export const getFormDefinition = createAsyncThunk<FormDefinition, string, { rejectValue: Status }>(
  'formDefinitions/get',
  async (id, { rejectWithValue }) => {
    try {
      return await Api.core().v1().formDefinitions().get(id).toPromise();
    } catch (err) {
      return rejectWithValue(err as Status);
    }
  }
);

export const createFormDefinition = createAsyncThunk<FormDefinition, FormDefinition>(
  'formDefinitions/create',
  async (formDefinition: FormDefinition, { rejectWithValue }) => {
    try {
      return await Api.core().v1().formDefinitions().create(formDefinition).toPromise();
    } catch (err) {
      return rejectWithValue(err);
    }
  }
);

export const updateFormDefinition = createAsyncThunk<FormDefinition, FormDefinition>(
  'formDefinitions/update',
  async (formDefinition: FormDefinition, { rejectWithValue }) => {
    try {
      return await Api.core().v1().formDefinitions().update(formDefinition).toPromise();
    } catch (err) {
      return rejectWithValue(err);
    }
  }
);

export const deleteFormDefinition = createAsyncThunk<void, string>(
  'formDefinitions/update',
  async (id: string, { rejectWithValue }) => {
    try {
      return await Api.core().v1().formDefinitions().delete(id).toPromise();
    } catch (err) {
      return rejectWithValue(err);
    }
  }
);

const INITIAL_STATE = {
  error: undefined as any,
  pending: false,
  ...formDefinitionsAdapter.getInitialState()
};

export type State = typeof INITIAL_STATE;
export type StateSlice = { formDefinitions: State }

const formDefinitionsSlice = createSlice({
  name: 'formDefinitions',
  initialState: INITIAL_STATE,
  reducers: {
    formDefinitionReceived: formDefinitionsAdapter.addOne,
    formDefinitionsReceived(state, action: PayloadAction<FormDefinitionList>) {
      formDefinitionsAdapter.setAll(state, action.payload.items ? action.payload.items : []);
    },
    formDefinitionsUpdated: formDefinitionsAdapter.updateOne,
    formDefinitionRemoved: formDefinitionsAdapter.removeOne
  },
  extraReducers: (builder) => {
    builder.addCase(listFormDefinitions.fulfilled, (state, { payload }) => {
      state = {...state}
      return {
        ...formDefinitionsAdapter.addMany(state, payload.items ? payload.items : []),
        pending: false,
        error: undefined
      };
    });
    builder.addCase(listFormDefinitions.rejected, (state, action) => {
      return { ...state, error: action.payload, pending: false };
    });
    builder.addCase(listFormDefinitions.pending, (state, action) => {
      return { ...state, pending: true, error: undefined };
    });
    builder.addCase(getFormDefinition.fulfilled, (state, { payload }) => {
      state = {...state}
      return {
        ...formDefinitionsAdapter.addOne(state, payload),
        pending: false,
        error: undefined
      };
    });
    builder.addCase(getFormDefinition.rejected, (state, action) => {
      return { ...state, error: action.payload, pending: false };
    });
    builder.addCase(getFormDefinition.pending, (state, action) => {
      return { ...state, pending: true, error: undefined };
    });
  }
});

export const {
  formDefinitionReceived,
  formDefinitionsReceived,
  formDefinitionsUpdated,
  formDefinitionRemoved
} = formDefinitionsSlice.actions;

export default formDefinitionsSlice.reducer;


export const selectAllFormDefinitions = (state: StateSlice) => {
  return state.formDefinitions.ids.map(id => state.formDefinitions.entities[id] as FormDefinition);
};

export const selectFormDefinitionById = (state: StateSlice, id: string) => {
  return state.formDefinitions.entities[id];
};

export const {
  selectAll,
  selectById,
  selectEntities,
  selectIds,
  selectTotal
} = formDefinitionsAdapter.getSelectors();
