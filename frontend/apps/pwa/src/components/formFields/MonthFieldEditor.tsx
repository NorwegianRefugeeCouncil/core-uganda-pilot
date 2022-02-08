import React, { ChangeEvent, FC } from 'react';
import { ErrorMessage } from '@hookform/error-message';

import { FieldEditorProps } from './types';
import { FieldDescription } from './FieldDescription';
import { FieldLabel } from './FieldLabel';

export const MonthFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
  register,
  errors,
}) => {
  const expectedLength = 7;
  const registerObject = register(`values.${field.id}`, {
    required: { value: field.required, message: 'This field is required' },
    pattern: {
      value: /^(?:19|20|21)\d{2}-[01]\d$/,
      message: 'wrong pattern',
    },
    maxLength: { value: expectedLength, message: 'Value is too long' },
    valueAsDate: { value: true, message: 'not a date' },
  });

  if (Array.isArray(value)) {
    return <></>;
  }
  return (
    <div className="form-group mb-2">
      <FieldLabel fieldDefinition={field} />
      <input
        className={`form-control bg-dark text-light border-secondary ${
          errors?.values && errors?.values[field.id] ? 'is-invalid' : ''
        }`}
        type="month"
        id={field.id}
        value={value || ''}
        name={field.name}
        placeholder="YYYY-MM"
        {...registerObject}
        onChange={(event) => {
          onChange(event.target.value);
          return registerObject.onChange(event);
        }}
        aria-describedby="errorMessages"
      />
      <FieldDescription text={field.description} />
      <div className="invalid-feedback" id="errorMessages">
        <div>
          <ErrorMessage errors={errors} name={`values.${field.id}`} />
        </div>
      </div>
    </div>
  );
};
