import React, { FC, useCallback } from 'react';
import { Database } from 'core-api-client';
import { FieldErrors, useFormContext } from 'react-hook-form';
import { ErrorMessage } from '@hookform/error-message';

import { useDatabases } from '../app/hooks';
import { registeredValidation } from '../features/former/validation';
import { ValidationForm } from '../reducers/Former/types';

type DatabasePickerProps = {
  databaseId: string | undefined;
  databases: Database[];
  setDatabaseId: (databaseId: string) => void;
  errors?: FieldErrors<ValidationForm>;
};

export const DatabasePicker: FC<DatabasePickerProps> = ({
  databases,
  databaseId,
  setDatabaseId,
  errors,
}) => {
  const { register } = useFormContext();

  const registerSelectedFieldReference = register(
    'selectedField.fieldType.reference.databaseId',
    registeredValidation.selectedField.fieldType.reference.databaseId,
  );

  return (
    <div>
      <select
        placeholder="Select Database"
        {...registerSelectedFieldReference}
        onChange={(e) => {
          setDatabaseId(e.target.value);
          return registerSelectedFieldReference.onChange(e);
        }}
        value={databaseId || ''}
        className={`form-select ${
          errors?.selectedField?.fieldType?.reference?.databaseId
            ? 'is-invalid'
            : ''
        }`}
        aria-label="Select Database"
        aria-describedby="errorMessages"
      >
        <option disabled value="">
          Select Database
        </option>
        {databases.map((d) => {
          return (
            <option key={d.id} value={d.id}>
              {d.name}
            </option>
          );
        })}
      </select>

      <div className="invalid-feedback" id="errorMessages">
        <ErrorMessage
          errors={errors}
          name="selectedField.fieldType.reference.databaseId"
        />
      </div>
    </div>
  );
};

type DatabasePickerContainerProps = {
  databaseId: string | undefined;
  setDatabaseId?: (databaseId: string) => void;
  setDatabase?: (database: Database | undefined) => void;
  errors?: FieldErrors<ValidationForm>;
};

const DatabasePickerContainer: FC<DatabasePickerContainerProps> = ({
  databaseId,
  setDatabaseId,
  setDatabase,
  errors,
}) => {
  const databases = useDatabases();

  const setDbCallback = useCallback(
    (dbId: string) => {
      if (setDatabaseId) {
        setDatabaseId(dbId);
      }
      const database = databases.find((d) => d.id === databaseId);
      if (setDatabase) {
        setDatabase(database);
      }
    },
    [databases, setDatabase, setDatabaseId],
  );

  return (
    <DatabasePicker
      databaseId={databaseId}
      setDatabaseId={setDbCallback}
      databases={databases}
      errors={errors}
    />
  );
};

export default DatabasePickerContainer;
