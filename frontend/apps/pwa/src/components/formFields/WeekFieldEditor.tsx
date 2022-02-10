import React, { FC } from 'react';
import { useFormContext } from 'react-hook-form';

import { FieldEditorProps } from './types';

export const WeekFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
  errors,
}) => {
  const { register } = useFormContext();

  const registerObject = register(`values.${field.id}`, {
    required: { value: field.required, message: 'This field is required' },
    pattern: {
      value: /^(?:19|20|21)\d{2}-W[0-5]\d$/,
      message: 'wrong pattern',
    },
    maxLength: { value: 8, message: 'Value is too long' },
    valueAsDate: true,
  });

  return (
    <input
      className={`form-control bg-dark text-light border-secondary ${
        errors?.values && errors?.values[field.id] ? 'is-invalid' : ''
      }`}
      type="week"
      placeholder="2021-W52"
      id={field.id}
      value={value || ''}
      {...registerObject}
      onChange={(event) => {
        onChange(event.target.value);
        return registerObject.onChange(event);
      }}
      aria-describedby={`errorMessages description-${field.id}`}
    />
  );
};
