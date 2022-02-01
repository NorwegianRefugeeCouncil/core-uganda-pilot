import React, { FC } from 'react';
import { ErrorMessage } from '@hookform/error-message';

import { FormerProps } from './types';
import { registeredValidation } from './validation';

export const FormName: FC<
  Pick<FormerProps, 'formName' | 'setFormName' | 'errors' | 'register'>
> = ({ formName = '', setFormName, errors, register }) => {
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
        aria-describedby="nameFeedback"
        {...registerObject}
        onChange={(event) => {
          setFormName(event.target.value);
          return registerObject.onChange(event);
        }}
      />

      <div className="invalid-feedback" id="nameFeedback">
        <ErrorMessage errors={errors} name="name" />
      </div>
    </div>
  );
};
