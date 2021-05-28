import {
  CauseType,
  FormDefinition,
  FormDefinitionVersion,
  FormElement, Path, pathFrom,
  Status,
  TranslatedStrings
} from '@core/api-client';
import { createAction, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { isArray } from 'util';

export const KEY = 'formBuilder';

export type State = {
  error: Partial<Status>,
  selectedVersion: string,
  formDefinition: FormDefinition
}

export type StateSlice = {
  formBuilder: State
}

export const INITIAL_STATE: State = {
  error: {
    details: {
      causes: []
    }
  },
  selectedVersion: '',
  formDefinition: {
    kind: 'FormDefinition',
    apiVersion: 'core.nrc.no/v1',
    metadata: {},
    spec: {
      group: '',
      names: {
        kind: '',
        plural: '',
        singular: ''
      },
      versions: [
        {
          name: '',
          storage: true,
          served: true,
          schema: {
            formSchema: {
              root: {}
            }
          }
        }
      ]
    }
  }
};

type setValuePayload = {
  path: string
  value: any
}
export const setValue = createAction<setValuePayload>('formBuilder/setValue');
const handleSetValue = (state: State, action: PayloadAction<setValuePayload>) => {
  const obj = state.formDefinition;
  const path = pathFrom(action.payload.path);
  path.set(obj, action.payload.value);
};

type setIndexedValuePayload = {
  path: string
  key: string
  value: any
}
export const setIndexedValue = createAction<setIndexedValuePayload>('formBuilder/setIndexedValue');
const handleSetIndexedValue = (state: State, action: PayloadAction<setIndexedValuePayload>) => {
  const obj = state.formDefinition;
  const path = pathFrom(action.payload.path);
  path.setIndexed(obj, action.payload.key, action.payload.value);
  state.formDefinition = obj;
};

type removeIndexedValuePayload = {
  path: string
  key: string
  keyValue: string
}
export const removeIndexedValue = createAction<removeIndexedValuePayload>('formBuilder/removeIndexedValue');
const handleRemoveIndexedValue = (state: State, action: PayloadAction<removeIndexedValuePayload>) => {
  const obj = state.formDefinition;
  const path = pathFrom(action.payload.path);
  path.removeIndexed(obj, action.payload.key, action.payload.keyValue);
};


type removeValuePayload = {
  path: string
}
export const removeValue = createAction<removeValuePayload>('formBuilder/removeValue');
const handleRemoveValue = (state: State, action: PayloadAction<removeValuePayload>) => {
  const obj = state.formDefinition;
  const path = pathFrom(action.payload.path);
  path.remove(obj);
};


type addValuePayload = {
  path: string
  value: any
}
export const addValue = createAction<addValuePayload>('formBuilder/addValue');
const handleAddValue = (state: State, action: PayloadAction<addValuePayload>) => {
  const obj = state.formDefinition;
  const path = pathFrom(action.payload.path);
  path.add(obj, action.payload.value);
};

type setFormDefinitionPayload = {
  formDefinition: FormDefinition
}
/**
 * Adds a new field to the state at the given path
 * Path is jsonPointer
 * Eg.
 * /root/children/3/children/2
 */
export const setFormDefinition = createAction<setFormDefinitionPayload>('formBuilder/setFormDefinition');
const handleSetFormDefinition = (state: State, action: PayloadAction<setFormDefinitionPayload>) => {
  state.formDefinition = action.payload.formDefinition;
};


export type TranslationType = 'label' | 'tooltip'


/**
 * The reducer for the formBuilder
 */
export const formBuilderSlice = createSlice({
  name: KEY,
  initialState: INITIAL_STATE,
  reducers: {},
  extraReducers: builder => {
    builder.addCase(setFormDefinition, handleSetFormDefinition);
    builder.addCase(setValue, handleSetValue);
    builder.addCase(removeValue, handleRemoveValue);
    builder.addCase(addValue, handleAddValue);
    builder.addCase(setIndexedValue, handleSetIndexedValue);
    builder.addCase(removeIndexedValue, handleRemoveIndexedValue);
  }
});

export const reducer = formBuilderSlice.reducer;

