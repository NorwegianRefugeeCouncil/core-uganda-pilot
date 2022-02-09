import React, { FC } from 'react';

import { FieldEditorProps } from './types';

export const TextFieldEditor: FC<FieldEditorProps> = ({
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
    <input
      className={`form-control bg-dark text-light border-secondary ${
        errors?.values && errors?.values[field.id] ? 'is-invalid' : ''
      }`}
      type="text"
      id={field.id}
      value={value || ''}
      {...registerObject}
      onChange={(event) => {
        onChange(event.target.value);
        return registerObject.onChange(event);
      }}
      aria-describedby="errorMessages"
    />
  );
};
