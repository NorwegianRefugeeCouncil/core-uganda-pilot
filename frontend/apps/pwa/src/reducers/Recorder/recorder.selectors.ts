import { FormDefinition, Record } from 'core-api-client';

import { RootState } from '../../app/store';
import {
  FormInterface,
  selectFormOrSubFormById,
  selectRootForm,
} from '../form';

import { adapter } from './recorder.adapter';
import { FormValue, RecordMap } from './types';

const selectors = adapter.getSelectors();
const globalSelectors = adapter.getSelectors<RootState>(
  (state) => state.recorder,
);

export const selectCurrentRecord = (
  state: RootState,
): FormValue | undefined => {
  return globalSelectors.selectById(state, state.recorder.selectedRecordId);
};

export const selectSubRecords = (
  state: RootState,
  recordId: string,
): RecordMap => {
  const result: RecordMap = {};
  const allRecords = globalSelectors.selectAll(state);
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
  const allEntries = [...globalSelectors.selectAll(state)];
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

export const recorderSelectors = selectors;

export const recorderGlobalSelectors = {
  ...globalSelectors,
  selectCurrentRecord,
  selectCurrentForm,
  selectCurrentRecordForm,
  selectCurrentRootForm,
  selectPostRecords,
  selectSubRecords,
};
