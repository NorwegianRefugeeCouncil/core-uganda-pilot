import React, { FC } from 'react';
import { useFormContext } from 'react-hook-form';

import { registeredValidation } from '../../features/former/validation';

import { FieldEditorProps } from './types';
import { SelectOptionsList } from './SelectOptionsList';

export const SingleSelectFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
  errors,
}) => {
  const { register } = useFormContext();

  const registerObject = register(
    `values.${field.id}`,
    registeredValidation.values(field),
  );

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
      aria-describedby={`errorMessages description-${field.id}`}
    >
      <SelectOptionsList field={field} />
    </select>
  );
};
