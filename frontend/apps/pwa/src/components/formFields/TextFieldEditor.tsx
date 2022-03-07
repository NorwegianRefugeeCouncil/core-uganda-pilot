import React, { FC } from 'react';
import { useFormContext } from 'react-hook-form';

import { registeredValidation } from '../../features/former/validation';
import { NonSubFormFieldValue } from '../../types/Field';

import { FieldEditorProps } from './types';

export const TextFieldEditor: FC<FieldEditorProps> = ({
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
    <input
      className={`form-control bg-dark text-light border-secondary ${
        errors?.values && errors?.values[field.id] ? 'is-invalid' : ''
      }`}
      type="text"
      id={field.id}
      value={(value as NonSubFormFieldValue['value']) || ''}
      {...registerObject}
      onChange={(event) => {
        onChange(event.target.value);
        return registerObject.onChange(event);
      }}
      aria-describedby={`errorMessages description-${field.id}`}
    />
  );
};
