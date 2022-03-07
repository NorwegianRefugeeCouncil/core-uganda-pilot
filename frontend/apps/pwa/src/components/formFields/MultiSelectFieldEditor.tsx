import React, { ChangeEvent, FC } from 'react';
import { useFormContext } from 'react-hook-form';

import { registeredValidation } from '../../features/former/validation';
import { NonSubFormFieldValue } from '../../types/Field';

import { FieldEditorProps } from './types';
import { SelectOptionsList } from './SelectOptionsList';

export const MultiSelectFieldEditor: FC<FieldEditorProps> = ({
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

  const handleChange = (event: ChangeEvent<HTMLSelectElement>) => {
    const { options } = event.target;
    const selected = Object.entries(options).filter((o) => o[1].selected);
    onChange(selected.map((s) => s[1].value));
  };

  return (
    <select
      className={`form-control bg-dark text-light border-secondary ${
        errors?.values && errors?.values[field.id] ? 'is-invalid' : ''
      }`}
      id={field.id}
      value={(value as NonSubFormFieldValue['value']) || []}
      multiple
      {...registerObject}
      onChange={(event) => {
        handleChange(event);
        return registerObject.onChange(event);
      }}
      aria-describedby={`errorMessages description-${field.id}`}
    >
      <SelectOptionsList field={field} isMultiSelect />
    </select>
  );
};
