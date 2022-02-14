import React, { FC } from 'react';
import { ErrorMessage } from '@hookform/error-message';
import { useFormContext } from 'react-hook-form';

import { FormerProps } from './types';
import { registeredValidation } from './validation';

export const FormName: FC<
  Pick<FormerProps, 'formName' | 'setFormName' | 'errors'>
> = ({ formName = '', setFormName, errors }) => {
  const { register } = useFormContext();

  const registerObject = register('name', registeredValidation.name);

  return (
    <div className="form-group mb-2">
      <label className="form-label" htmlFor="name">
        Form Name
      </label>
      <input
        className={`form-control ${errors?.name ? 'is-invalid' : ''}`}
        id="name"
        type="text"
        value={formName || ''}
        aria-describedby="errorMessages-formName"
        {...registerObject}
        onChange={(event) => {
          setFormName(event.target.value);
          return registerObject.onChange(event);
        }}
      />

      <div className="invalid-feedback" id="errorMessages-formName">
        <ErrorMessage errors={errors} name="name" />
      </div>
    </div>
  );
};
