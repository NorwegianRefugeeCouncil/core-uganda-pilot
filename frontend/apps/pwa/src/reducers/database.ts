import { createAsyncThunk, createEntityAdapter, createSlice, EntityState } from '@reduxjs/toolkit';
import { Database, DatabaseList, DatabaseListRequest, Response } from 'core-api-client';

import { RootState } from '../app/store';
import client from '../app/client';

const adapter = createEntityAdapter<Database>({
  // Assume IDs are stored in a field other than `book.id`
  selectId: (db) => db.id,
  // Keep the "all IDs" array sorted based on book titles
  sortComparer: (a, b) => a.name.localeCompare(b.name),
});

export const fetchDatabases = createAsyncThunk<Response<DatabaseListRequest, DatabaseList>>(
  'databases/fetch',
  async (_, thunkAPI) => {
    try {
      const response = await client.listDatabases({});
      if (response.success) {
        return response;
      }
      return thunkAPI.rejectWithValue(response);
    } catch (err) {
      return thunkAPI.rejectWithValue(err);
    }
  },
);

export interface DatabaseState extends EntityState<Database> {
  fetchPending: boolean;
  fetchError: any;
  fetchSuccess: boolean;
}

export const databasesSlice = createSlice({
  name: 'databases',
  initialState: {
    ...adapter.getInitialState(),
    fetchPending: false,
    fetchError: undefined as any,
    fetchSuccess: true,
  },
  reducers: {
    addOne: adapter.addOne,
    addMany: adapter.addMany,
    removeAll: adapter.removeAll,
    removeMany: adapter.removeMany,
    removeOne: adapter.removeOne,
    updateMany: adapter.updateMany,
    updateOne: adapter.updateOne,
    upsertOne: adapter.upsertOne,
    upsertMany: adapter.upsertMany,
    setOne: adapter.setOne,
    setMany: adapter.setMany,
    setAll: adapter.setAll,
  },
  extraReducers: (builder) => {
    builder.addCase(fetchDatabases.pending, (state, action) => {
      state.fetchSuccess = false;
      state.fetchPending = true;
    });
    builder.addCase(fetchDatabases.rejected, (state, action) => {
      state.fetchSuccess = false;
      state.fetchPending = false;
      state.fetchError = action.payload;
    });
    builder.addCase(fetchDatabases.fulfilled, (state, action) => {
      state.fetchSuccess = true;
      state.fetchPending = false;
      state.fetchError = undefined;
      if (action.payload.response?.items) {
        adapter.setAll(state, action.payload.response.items);
      }
    });
  },
});

export const databaseActions = databasesSlice.actions;
const selectors = adapter.getSelectors();

export const databaseSelectors = {
  ...selectors,
};

export const databaseGlobalSelectors = adapter.getSelectors<RootState>((state) => state.databases);

export default databasesSlice.reducer;
