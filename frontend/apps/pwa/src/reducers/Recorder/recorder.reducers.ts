import {
  ActionReducerMapBuilder,
  createAsyncThunk,
  PayloadAction,
} from '@reduxjs/toolkit';
import { v4 as uuidv4 } from 'uuid';
import { Record } from 'core-api-client';

import client from '../../app/client';
import { RootState } from '../../app/store';
import { selectFormOrSubFormById } from '../form';
import { recordGlobalSelectors } from '../records';

import { FormValue, RecorderState } from './types';
import { adapter } from './recorder.adapter';

export const reducers = {
  setFieldValue(
    state: RecorderState,
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
    state: RecorderState,
    action: PayloadAction<{ recordId: string; fieldId: string }>,
  ) {
    const { recordId, fieldId } = action.payload;
    const record = state.entities[recordId];
    if (!record) {
      return;
    }
    record.values = record.values.filter((v) => v.fieldId !== fieldId);
  },
  selectRecord(
    state: RecorderState,
    action: PayloadAction<{ recordId: string }>,
  ) {
    state.selectedRecordId = action.payload.recordId;
  },
  initRecord(state: RecorderState, action: PayloadAction<{ formId: string }>) {
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
    state: RecorderState,
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
};

export const resetForm = createAsyncThunk<
  { formValue: FormValue },
  { formId: string; ownerId: string | undefined }
>('records/resetForm', ({ formId, ownerId }, { rejectWithValue, getState }) => {
  const state = getState() as RootState;

  const form = selectFormOrSubFormById(state, formId);
  if (!form) {
    return rejectWithValue(`could not find form or sub form with id ${formId}`);
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
});

export const postRecord = createAsyncThunk<Record[], Record[]>(
  'records/post',
  async (arg, thunkAPI) => {
    const result: Record[] = [];
    for (const record of arg) {
      try {
        const response = await client.Record.create({ object: record });
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

export const extraReducers = (builder: ActionReducerMapBuilder<any>) => {
  builder.addCase(resetForm.fulfilled, (state, action) => {
    adapter.removeAll(state);
    state.baseFormId = action.payload.formValue.formId;
    state.editingValues = {};
    adapter.addOne(state, action.payload.formValue);
    state.selectedRecordId = action.payload.formValue.id;
  });
  builder.addCase(postRecord.rejected, (state, payload) => {
    state.saveError = payload.payload as any;
  });
};
