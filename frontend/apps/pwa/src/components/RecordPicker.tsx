import React, { FC, useCallback, useEffect, useState } from 'react';
import { FieldDefinition, Record } from 'core-api-client';
import { FieldErrors, useFormContext } from 'react-hook-form';
import { ErrorMessage } from '@hookform/error-message';

import { useAppDispatch, useAppSelector } from '../app/hooks';
import {
  fetchRecords,
  recordGlobalSelectors,
  selectRecords,
} from '../reducers/records';
import { selectFormOrSubFormById, selectRootForm } from '../reducers/form';
import { FieldDefinitionNC } from '../reducers/Former/types';
import { registeredValidation } from '../features/former/validation';

export type RecordPickerProps = {
  disabled?: boolean;
  recordId: string | null;
  field: FieldDefinition;
  setRecordId: (recordId: string | null) => void;
  records: Record[];
  getDisplayStr: (record: Record) => string;
  errors: FieldErrors;
};

export const RecordPicker: FC<RecordPickerProps> = ({
  disabled,
  recordId,
  setRecordId,
  field,
  records,
  getDisplayStr,
  errors,
}) => {
  const { register } = useFormContext();

  const registerObject = register(
    `values.${field.id}`,
    registeredValidation.values(field),
  );

  return (
    <div>
      <select
        disabled={disabled}
        value={recordId || ''}
        className={`form-select ${
          errors?.values && errors?.values[field.id] ? 'is-invalid' : ''
        }`}
        aria-label="Select Record"
        {...registerObject}
        onChange={(event) => {
          setRecordId(event.target.value);
          return registerObject.onChange(event);
        }}
        aria-describedby="errorMessages"
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

      <div className="invalid-feedback" id="errorMessages">
        <div>
          <ErrorMessage errors={errors} name={`values.${field.id}`} />
        </div>
      </div>
    </div>
  );
};

export type RecordPickerContainerProps = {
  recordId: string | null;
  setRecordId?: (recordId: string | null) => void;
  setRecord?: (record: Record | undefined) => void;
  ownerId?: string;
  formId?: string;
  field: FieldDefinition;
  errors: FieldErrors;
};

export const RecordPickerContainer: FC<RecordPickerContainerProps> = ({
  ownerId,
  formId,
  field,
  recordId,
  setRecordId,
  setRecord,
  errors,
}) => {
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
      dispatch(
        fetchRecords({ databaseId: rootForm.databaseId, formId: form?.id }),
      )
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
      field={field}
      setRecordId={callback}
      records={records}
      getDisplayStr={(r) => {
        const result = form?.fields.reduce((prev, next) => {
          const fieldVal = r.values.find((v) => v.fieldId === next.id)?.value;
          return fieldVal ? `${prev} ${fieldVal}` : prev;
        }, '');
        return result || '';
      }}
      errors={errors}
    />
  );
};
