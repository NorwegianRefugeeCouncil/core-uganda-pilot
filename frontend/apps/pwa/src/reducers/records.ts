import {
  createAsyncThunk,
  createEntityAdapter,
  createSlice,
} from '@reduxjs/toolkit';
import {
  FieldDefinition,
  FormLookup,
  Record,
  RecordList,
} from 'core-api-client';

import { RootState } from '../app/store';
import client from '../app/client';

import {
  formGlobalSelectors,
  selectFieldForSubForm,
  selectFormOrSubFormById,
  selectRootForm,
  selectSubFormForField,
} from './form';

const adapter = createEntityAdapter<Record>({
  // Assume IDs are stored in a field other than `book.id`
  selectId: (db) => db.id,
  // Keep the "all IDs" array sorted based on book titles
  sortComparer: (a, b) => a.id.localeCompare(b.id),
});

export const fetchRecords = createAsyncThunk<RecordList, FormLookup>(
  'records/fetch',
  async (args, thunkAPI) => {
    try {
      return await client.Record.list(args);
    } catch (err) {
      return thunkAPI.rejectWithValue(err);
    }
  },
);

export const records = createSlice({
  name: 'records',
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
    builder.addCase(fetchRecords.pending, (state, action) => {
      state.fetchSuccess = false;
      state.fetchPending = true;
    });
    builder.addCase(fetchRecords.rejected, (state, action) => {
      state.fetchSuccess = false;
      state.fetchPending = false;
      state.fetchError = action.payload;
    });
    builder.addCase(fetchRecords.fulfilled, (state, action) => {
      state.fetchSuccess = true;
      state.fetchPending = false;
      state.fetchError = undefined;
      if (action.payload.items) {
        adapter.addMany(state, action.payload.items);
      }
    });
  },
});

export const recordActions = records.actions;
export const recordSelectors = adapter.getSelectors();
export const recordGlobalSelectors = adapter.getSelectors<RootState>(
  (state) => state.records,
);
export default records.reducer;

export const selectRecordsForForm = (
  state: RootState,
  formId: string,
  ownerRecordId?: string,
) => {
  return recordGlobalSelectors.selectAll(state).filter((r) => {
    return ownerRecordId
      ? r.formId === formId && r.ownerId === ownerRecordId
      : r.formId === formId;
  });
};

export const selectSubFormCount = (
  recordId: string,
  fieldId: string,
): ((rootState: RootState) => number) => {
  return (rootState) => {
    const record = recordGlobalSelectors.selectById(rootState, recordId);
    if (!record) {
      return 0;
    }
    const subForm = selectSubFormForField(fieldId)(rootState);
    if (!subForm) {
      return 0;
    }
    const subFormId = fieldId;
    const allRecords = recordGlobalSelectors.selectAll(rootState);
    return allRecords.reduce((prec, current) => {
      if (current.ownerId === recordId && current.formId === subFormId) {
        return prec++;
      }
      return prec;
    }, 0);
  };
};

export type SubRecordResult = {
  byFormId: { [formId: string]: Record[] };
  byFieldId: { [fieldId: string]: Record[] };
};

export const selectSubRecords: (
  state: RootState,
  recordId: string,
) => SubRecordResult = (state, recordId) => {
  const result: SubRecordResult = {
    byFieldId: {},
    byFormId: {},
  };

  const record = recordGlobalSelectors.selectById(state, recordId);
  if (!record) {
    return result;
  }

  const form = selectFormOrSubFormById(state, record.formId);
  if (!form) {
    return result;
  }

  const rootForm = selectRootForm(state, form.id);
  if (!rootForm) {
    return result;
  }

  const fieldMap: { [formId: string]: FieldDefinition } = {};

  const subFormIds = new Set<string>();
  for (const field of form.fields) {
    if (field.fieldType.subForm) {
      const fieldForSubForm = selectFieldForSubForm(rootForm, field.id);
      if (fieldForSubForm) {
        fieldMap[field.id] = fieldForSubForm;
      }
      subFormIds.add(field.id);
    }
  }

  const allRecords = recordGlobalSelectors.selectAll(state);
  for (const candidateRecord of allRecords) {
    if (
      candidateRecord.ownerId === recordId &&
      subFormIds.has(candidateRecord.formId)
    ) {
      const fieldForSubform = fieldMap[candidateRecord.formId];
      if (fieldForSubform) {
        if (!result.byFieldId.hasOwnProperty(fieldForSubform.id)) {
          result.byFieldId[fieldForSubform.id] = [];
        }
        result.byFieldId[fieldForSubform.id].push(candidateRecord);
        if (!result.byFormId.hasOwnProperty(candidateRecord.formId)) {
          result.byFormId[candidateRecord.formId] = [];
        }
        result.byFormId[candidateRecord.formId].push(candidateRecord);
      }
    }
  }

  return result;
};

export const selectRecordsSubFormCounts: (formId?: string) => (
  rootState: RootState,
) => {
  [recordId: string]: { [fieldId: string]: number };
} = (formId) => {
  return (rootState) => {
    if (!formId) {
      return {};
    }

    const result: { [recordId: string]: { [fieldId: string]: number } } = {};

    const form = formGlobalSelectors.selectById(rootState, formId);
    if (!form) {
      return {};
    }

    if (!form.fields) {
      return {};
    }

    // maps which form ids correspond to which field ids [formId] -> [fieldId]
    const formIdFieldIdMap: { [formId: string]: string } = {};
    for (const formField of form?.fields) {
      if (!formField?.fieldType?.subForm) {
        continue;
      }
      formIdFieldIdMap[formField.id] = formField.id;
    }

    // records for the given form
    const formRecords = selectRecordsForForm(rootState, formId);
    if (formRecords.length === 0) {
      return {};
    }

    // map from recordId -> record
    const formRecordsMap: { [key: string]: Record } = {};
    for (const record of formRecords) {
      formRecordsMap[record.id] = record;
    }

    // all records
    const allRecords = recordGlobalSelectors.selectAll(rootState);

    for (const record of allRecords) {
      // record does not have a owner record, does not qualify as subform
      if (!record.ownerId) {
        continue;
      }
      // owner record is not part of the current form
      if (!formRecordsMap.hasOwnProperty(record.ownerId)) {
        continue;
      }
      // field is not part of the current form
      if (!formIdFieldIdMap.hasOwnProperty(record.formId)) {
        continue;
      }
      // the field id for that sub record
      const recordFieldId = formIdFieldIdMap[record.formId];

      // construct result
      if (!result.hasOwnProperty(record.ownerId)) {
        result[record.ownerId] = {};
      }
      if (!result[record.ownerId].hasOwnProperty(recordFieldId)) {
        result[record.ownerId][recordFieldId] = 0;
      }

      // increase
      result[record.ownerId][recordFieldId]++;
    }

    return result;
  };
};

export function selectRecords(
  state: RootState,
  options: { formId?: string; ownerId?: string },
) {
  const allRecords = recordGlobalSelectors.selectAll(state);
  const result: Record[] = [];
  for (const record of allRecords) {
    if (options.formId && record.formId !== options.formId) {
      continue;
    }
    if (options.ownerId && record.ownerId !== options.ownerId) {
      continue;
    }
    result.push(record);
  }
  return result;
}
