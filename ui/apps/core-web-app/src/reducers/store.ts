import { createAsyncThunk, createEntityAdapter, createSlice } from '@reduxjs/toolkit';
import { HasObjectMeta, RuntimeList, RuntimeObject, Status } from '@core/api-client';
import { Api } from '../data/api';


export function useInsertOne(){
  const db = usePouch()
}



type ResourceEnvelope = {
  resource: RuntimeObject & HasObjectMeta
  isNew: boolean
  isUpdated: boolean
}

type ErrorEnvelope = {
  status: Status
}

const getResourceId = (apiVersion: string, kind: string, name: string): string => {
  return apiVersion + '/' + kind + '/' + name;
};

const resourcesAdapter = createEntityAdapter<ResourceEnvelope>({
  selectId: (f) => getResourceId(f.resource.apiVersion, f.resource.kind, f.resource.metadata.name),
  sortComparer: (a, b) => {
    return getResourceId(a.resource.apiVersion, a.resource.kind, a.resource.metadata.name)
      .localeCompare(getResourceId(b.resource.apiVersion, b.resource.kind, b.resource.metadata.name));
  }
});

const INITIAL_STATE = {
  ...resourcesAdapter.getInitialState()
};

export type State = typeof INITIAL_STATE
export type StateSlice = { resources: State }

export const list = createAsyncThunk<RuntimeList, {
  group: string,
  resource: string,
  version: string
}, { rejectValue: Status }>(
  'resources/list',
  async ({ group, resource, version }, { rejectWithValue }) => {
    try {
      return await Api.restClient()
        .get()
        .group(group)
        .resource(resource)
        .version(version)
        .do<RuntimeList>()
        .toPromise();
    } catch (err) {
      return rejectWithValue(err as Status);
    }
  }
);

const resources = createSlice({
  name: 'resources',
  initialState: INITIAL_STATE,
  reducers: {
    addOne: resourcesAdapter.addOne,
    addMany: resourcesAdapter.addMany,
    setAll: resourcesAdapter.setAll,
    removeOne: resourcesAdapter.removeOne,
    removeAll: resourcesAdapter.removeAll,
    updateOne: resourcesAdapter.updateOne,
    updateMany: resourcesAdapter.updateMany,
    upsertOne: resourcesAdapter.upsertOne,
    upsertMany: resourcesAdapter.upsertMany
  },
  extraReducers: builder => {
    builder.addCase(list.fulfilled, (state, action) => {

      const stateIds: { [key: string]: boolean } = {};
      const payloadIds: { [key: string]: boolean } = {};

      let list = action.payload;
      let kind = list.kind;
      if (kind.endsWith('List')) {
        kind = kind.substring(0, kind.length - 4);
      }
      for (let item of list.items) {
        payloadIds[getResourceId(list.apiVersion, kind, item.metadata.name)] = true;
      }
      for (let item of state.ids) {
        stateIds[item] = true;
      }



    });
  }
});
