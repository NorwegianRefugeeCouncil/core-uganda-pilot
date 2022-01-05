import React, { FC } from 'react';

import { FieldEditorProps } from './types';
import {
  mapFieldDescription,
  mapFieldLabel,
  mapSelectOptions,
} from './helpers';

export const SingleSelectFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  setValue,
}) => {
  return (
    <div className="form-group mb-2">
      {mapFieldLabel(field)}
      <select
        className="form-control bg-dark text-light border-secondary"
        id={field.id}
        value={value || ''}
        onChange={(event) => setValue(event.target.value)}
      >
        {mapSelectOptions(
          field.required,
          field.key,
          field.fieldType?.singleSelect?.options,
        )}
      </select>
      {mapFieldDescription(field)}
    </div>
  );
};
