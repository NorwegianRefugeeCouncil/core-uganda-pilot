import React, { FC, useCallback, useEffect, useState } from 'react';
import { Record } from 'core-api-client';

import { useAppDispatch, useAppSelector } from '../app/hooks';
import { fetchRecords, recordGlobalSelectors, selectRecords } from '../reducers/records';
import { selectFormOrSubFormById, selectRootForm } from '../reducers/form';

export type RecordPickerProps = {
  disabled?: boolean;
  recordId: string | null;
  setRecordId: (recordId: string | null) => void;
  records: Record[];
  getDisplayStr: (record: Record) => string;
};

export const RecordPicker: FC<RecordPickerProps> = (props) => {
  const { disabled, recordId, setRecordId, records, getDisplayStr } = props;

  return (
    <div>
      <select
        disabled={disabled}
        onChange={(e) => setRecordId(e.target.value)}
        value={recordId || ''}
        className="form-select"
        aria-label="Select Record"
      >
        <option disabled value="">
          No Records
        </option>
        {records.map((r) => {
          return (
            <option key={r.id} value={r.id}>
              {getDisplayStr(r)}
            </option>
          );
        })}
      </select>
    </div>
  );
};

export type RecordPickerContainerProps = {
  recordId: string | null;
  setRecordId?: (recordId: string | null) => void;
  setRecord?: (record: Record | undefined) => void;
  ownerId?: string;
  formId?: string;
};

export const RecordPickerContainer: FC<RecordPickerContainerProps> = (props) => {
  const { ownerId, formId, recordId, setRecordId, setRecord } = props;

  const dispatch = useAppDispatch();

  const form = useAppSelector((state) => {
    return selectFormOrSubFormById(state, formId || '');
  });

  const rootForm = useAppSelector((state) => {
    return selectRootForm(state, form ? form.id : '');
  });

  const [pending, setPending] = useState(false);

  useEffect(() => {
    if (rootForm && form && !pending) {
      dispatch(fetchRecords({ databaseId: rootForm.databaseId, formId: form?.id }))
        .then(() => {
          setPending(true);
        })
        .catch(() => {
          setPending(false);
        });
    }
  }, [dispatch, form, pending, rootForm]);

  const records = useAppSelector((state) => {
    return selectRecords(state, { ownerId, formId });
  });

  const record = useAppSelector((state) => {
    return recordGlobalSelectors.selectById(state, recordId || '');
  });

  const callback = useCallback(
    (recId: string | null) => {
      if (setRecord) {
        setRecord(record);
      }
      if (setRecordId) {
        setRecordId(recId);
      }
    },
    [record, setRecord, setRecordId],
  );

  return (
    <RecordPicker
      disabled={false}
      recordId={recordId}
      setRecordId={callback}
      records={records}
      getDisplayStr={(r) => {
        const result = form?.fields.reduce((prev, next) => {
          const fieldVal = r.values.find((v) => v.fieldId === next.id)?.value;
          return fieldVal ? `${prev} ${fieldVal}` : prev;
        }, '');
        return result || '';
      }}
    />
  );
};
