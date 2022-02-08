import React, { ChangeEvent, FC } from 'react';
import { ErrorMessage } from '@hookform/error-message';

import { FieldEditorProps } from './types';
import { FieldDescription } from './FieldDescription';
import { FieldLabel } from './FieldLabel';
import { SelectOptionsList } from './SelectOptionsList';

export const MultiSelectFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
  register,
  errors,
}) => {

  const registerObject = register(`values.${field.id}`, {
    required: { value: field.required, message: 'This field is required' },
  });

  const handleChange = (event: ChangeEvent<HTMLSelectElement>) => {
    const { options } = event.target;
    const selected = Object.entries(options).filter((o) => o[1].selected);
    onChange(selected.map((s) => s[1].value));
  };

  return (
    <div className="form-group mb-2">
      <FieldLabel fieldDefinition={field} />
      <select
        className={`form-control bg-dark text-light border-secondary ${
          errors?.values && errors?.values[field.id] ? 'is-invalid' : ''
        }`}
        id={field.id}
        value={value || []}
        multiple
        {...registerObject}
        onChange={(event) => {
          handleChange(event);
          return registerObject.onChange(event);
        }}
        aria-describedby="errorMessages"
      >
        <SelectOptionsList field={field} isMultiSelect />
      </select>
      <FieldDescription text={field.description} />
      <div className="invalid-feedback" id="errorMessages">
        <ErrorMessage errors={errors} name={`values.${field.id}`} />
      </div>
    </div>
  );
};
