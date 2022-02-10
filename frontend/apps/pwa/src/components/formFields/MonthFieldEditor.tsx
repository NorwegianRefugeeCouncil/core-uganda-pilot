import React, { FC } from 'react';
import { useFormContext } from 'react-hook-form';

import { FieldEditorProps } from './types';

export const MonthFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
  errors,
}) => {
  const { register } = useFormContext();

  const registerObject = register(`values.${field.id}`, {
    required: { value: field.required, message: 'This field is required' },
    pattern: {
      value: /^(?:19|20|21)\d{2}-[01]\d$/,
      message: 'wrong pattern',
    },
    valueAsDate: true,
  });

  if (Array.isArray(value)) {
    return <></>;
  }

  return (
    <input
      className={`form-control bg-dark text-light border-secondary ${
        errors?.values && errors?.values[field.id] ? 'is-invalid' : ''
      }`}
      type="month"
      id={field.id}
      value={value || ''}
      placeholder="YYYY-MM"
      {...registerObject}
      onChange={(event) => {
        onChange(event.target.value);
        return registerObject.onChange(event);
      }}
      aria-describedby={`errorMessages description-${field.id}`}
    />
  );
};
