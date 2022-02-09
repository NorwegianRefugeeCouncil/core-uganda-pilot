import React, { FC } from 'react';

import { FieldEditorProps } from './types';
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
  );
};
