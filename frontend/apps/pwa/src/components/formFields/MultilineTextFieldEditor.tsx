import React, { FC } from 'react';
import { ErrorMessage } from '@hookform/error-message';

import { FieldEditorProps } from './types';
import { FieldDescription } from './FieldDescription';
import { FieldLabel } from './FieldLabel';

export const MultilineTextFieldEditor: FC<FieldEditorProps> = ({
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
      <textarea
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
