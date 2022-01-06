import {
  createAsyncThunk,
  createEntityAdapter,
  createSlice,
  PayloadAction,
} from '@reduxjs/toolkit';
import { v4 as uuidv4 } from 'uuid';
import { FormDefinition, Record } from 'core-api-client';

import { RootState } from '../app/store';
import client from '../app/client';

import { FormInterface, selectFormOrSubFormById, selectRootForm } from './form';
import { recordGlobalSelectors } from './records';

export interface FormValue extends Omit<Record, 'databaseId'> {
  // records the sub form field that the record belongs to, if any
  ownerFieldId?: string;
}

type RecordMap = { [key: string]: FormValue[] };

const adapter = createEntityAdapter<FormValue>({
  // Assume IDs are stored in a field other than `book.id`
  selectId: (folder) => folder.id,
  // Keep the "all IDs" array sorted based on book titles
  sortComparer: (a, b) => a.id.localeCompare(b.id),
});

export const resetForm = createAsyncThunk<
  { formValue: FormValue },
  { formId: string; ownerId: string | undefined }
>(
  'records/resetForm',
  ({ formId, ownerId }, { rejectWithValue, fulfillWithValue, getState }) => {
    const state = getState() as RootState;

    const form = selectFormOrSubFormById(state, formId);
    if (!form) {
      return rejectWithValue(
        `could not find form or sub form with id ${formId}`,
      );
    }

    const newRecord: FormValue = {
      id: uuidv4(),
      formId: form.id,
      values: [],
      ownerId: undefined,
    };

    if (ownerId) {
      const baseRecord = recordGlobalSelectors.selectById(state, ownerId);
      if (!baseRecord) {
        return rejectWithValue(`cannot find record with id ${ownerId}`);
      }
      newRecord.ownerId = ownerId;

      const ownerFormId = baseRecord.formId;
      const ownerForm = selectFormOrSubFormById(state, ownerFormId);
      if (!ownerForm) {
        return rejectWithValue(`cannot find form with id ${ownerId}`);
      }

      const ownerField = ownerForm.fields.find((f) => {
        if (!f.fieldType.subForm) {
          return false;
        }
        return f.id === formId;
      });
      if (!ownerField) {
        return rejectWithValue(`cannot find subform field with id ${formId}`);
      }

      newRecord.ownerFieldId = ownerField.id;
    }

    return { formValue: newRecord };
  },
);

export const recorder = createSlice({
  name: 'recorder',
  initialState: {
    ...adapter.getInitialState(),
    selectedRecordId: '',
    baseFormId: '',
    editingValues: {} as { [recordId: string]: { [key: string]: any } },
  },
  reducers: {
    setFieldValue(
      state,
      action: PayloadAction<{ recordId: string; fieldId: string; value: any }>,
    ) {
      const { recordId, fieldId, value } = action.payload;
      const record = state.entities[recordId];
      if (!record) {
        return;
      }

      const idx = record.values.findIndex((v) => v.fieldId === fieldId);
      if (idx === -1) {
        record.values = [...record.values, { fieldId, value }];
      } else {
        record.values[idx] = { ...record.values[idx], ...{ value } };
      }
    },
    clearFieldValue(
      state,
      action: PayloadAction<{ recordId: string; fieldId: string }>,
    ) {
      const { recordId, fieldId } = action.payload;
      const record = state.entities[recordId];
      if (!record) {
        return;
      }
      record.values = record.values.filter((v) => v.fieldId !== fieldId);
    },
    selectRecord(state, action: PayloadAction<{ recordId: string }>) {
      state.selectedRecordId = action.payload.recordId;
    },
    initRecord(state, action: PayloadAction<{ formId: string }>) {
      const newRecord: FormValue = {
        id: uuidv4(),
        formId: action.payload.formId,
        values: [],
        ownerId: undefined,
      };
      adapter.addOne(state, newRecord);
      state.selectedRecordId = newRecord.id;
    },
    addSubRecord(
      state,
      action: PayloadAction<{
        formId: string;
        ownerFieldId: string;
        ownerRecordId: string;
      }>,
    ) {
      const newRecord: FormValue = {
        id: uuidv4(),
        formId: action.payload.formId,
        values: [],
        ownerFieldId: action.payload.ownerFieldId,
        ownerId: action.payload.ownerRecordId,
      };
      adapter.addOne(state, newRecord);
      state.selectedRecordId = newRecord.id;
    },
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
    builder.addCase(resetForm.fulfilled, (state, action) => {
      adapter.removeAll(state);
      state.baseFormId = action.payload.formValue.formId;
      state.editingValues = {};
      adapter.addOne(state, action.payload.formValue);
      state.selectedRecordId = action.payload.formValue.id;
    });
  },
});

export const recorderActions = recorder.actions;
export const recorderSelectors = adapter.getSelectors();
export const recorderGlobalSelectors = adapter.getSelectors<RootState>(
  (state) => state.recorder,
);
export default recorder.reducer;

export const selectCurrentRecord = (
  state: RootState,
): FormValue | undefined => {
  return recorderGlobalSelectors.selectById(
    state,
    state.recorder.selectedRecordId,
  );
};

export const selectSubRecords = (
  state: RootState,
  recordId: string,
): RecordMap => {
  const result: RecordMap = {};
  const allRecords = recorderGlobalSelectors.selectAll(state);
  for (const record of allRecords) {
    if (record.ownerId === recordId && record.ownerFieldId) {
      if (!result.hasOwnProperty(record.ownerFieldId)) {
        result[record.ownerFieldId] = [];
      }
      result[record.ownerFieldId].push(record);
    }
  }
  return result;
};

export const selectCurrentForm = (
  state: RootState,
): FormInterface | undefined => {
  const selectedRecord = selectCurrentRecord(state);
  if (!selectedRecord) {
    return;
  }
  return selectFormOrSubFormById(state, selectedRecord.formId);
};

export const selectCurrentRecordForm = (
  state: RootState,
): FormInterface | undefined => {
  const currentRecord = selectCurrentRecord(state);
  if (!currentRecord) {
    return undefined;
  }
  return selectFormOrSubFormById(state, currentRecord.formId);
};

export const selectCurrentRootForm = (
  state: RootState,
): FormDefinition | undefined => {
  const currentRecord = selectCurrentRecord(state);
  if (!currentRecord) {
    return undefined;
  }
  return selectRootForm(state, currentRecord.formId);
};

export const selectPostRecords = (state: RootState): Record[] => {
  const result: Record[] = [];
  const allEntries = [...recorderGlobalSelectors.selectAll(state)];
  const handledRecords: { [key: string]: boolean } = {};
  const { baseFormId } = state.recorder;
  const baseForm = selectFormOrSubFormById(state, baseFormId);
  if (!baseForm) {
    return [];
  }
  const rootForm = selectRootForm(state, baseForm.id);
  if (!rootForm) {
    return [];
  }
  const { databaseId } = rootForm;

  for (
    let i = allEntries.length - 1;
    allEntries.length > 0;
    i === 0 ? (i = allEntries.length - 1) : i--
  ) {
    const entry = allEntries[i];

    if (baseFormId !== rootForm.id) {
      if (
        entry.ownerId &&
        entry.formId !== baseFormId &&
        !handledRecords[entry.ownerId]
      ) {
        continue;
      }
    } else if (entry.ownerId && !handledRecords[entry.ownerId]) {
      continue;
    }

    const record: Record = {
      formId: entry.formId,
      id: entry.id,
      databaseId,
      values: entry.values,
      ownerId: entry.ownerId,
    };
    result.push(record);
    handledRecords[record.id] = true;
    allEntries.splice(i, 1);
  }
  return result;
};

export const postRecord = createAsyncThunk<Record[], Record[]>(
  'records/post',
  async (arg, thunkAPI) => {
    const result: Record[] = [];
    for (const record of arg) {
      try {
        const response = await client.createRecord({ object: record });
        if (!response.success) {
          return thunkAPI.rejectWithValue(response.error);
        }
        if (!response.response) {
          return thunkAPI.rejectWithValue('no record in response');
        }
        for (let i = 1; i < arg.length; i++) {
          const otherRecord = arg[i];
          if (otherRecord.ownerId === record.id) {
            otherRecord.ownerId = response.response.id;
          }
        }
        result.push(response.response);
      } catch (err) {
        return thunkAPI.rejectWithValue(err);
      }
    }
    return result;
  },
);
