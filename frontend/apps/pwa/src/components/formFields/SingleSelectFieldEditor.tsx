import React, { FC } from 'react';
import { ErrorMessage } from '@hookform/error-message';

import { FieldEditorProps } from './types';
import { FieldDescription } from './FieldDescription';
import { FieldLabel } from './FieldLabel';
import { SelectOptionsList } from './SelectOptionsList';

export const SingleSelectFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
  register,
  errors,
}) => {
  const registerObject = register(`values.${field.id}`, {
    required: { value: field.required, message: 'This field is required' },
  });
  return (
    <div className="form-group mb-2">
      <FieldLabel fieldDefinition={field} />
      <select
        className={`form-control bg-dark text-light border-secondary ${
          errors?.values && errors?.values[field.id] ? 'is-invalid' : ''
        }`}
        id={field.id}
        value={value || ''}
        {...registerObject}
        onChange={(event) => {
          onChange(event.target.value);
          return registerObject.onChange(event);
        }}
        aria-describedby="errorMessages"
      >
        <SelectOptionsList field={field} />
      </select>

      <FieldDescription text={field.description} />
      <div className="invalid-feedback" id="errorMessages">
        <div>
          <ErrorMessage errors={errors} name={`values.${field.id}`} />
        </div>
      </div>
    </div>
  );
};
