import React, { FC } from 'react';
import _ from 'lodash';
import { ErrorMessage } from '@hookform/error-message';

import { FormerProps } from './types';
import validation from './validation';

export const FormName: FC<
  Pick<FormerProps, 'formName' | 'setFormName' | 'errors' | 'register'>
> = ({ formName = '', setFormName, errors, register }) => {
  const registerObject = register('name', validation.name);

  return (
    <div className="form-group mb-2">
      <label className="form-label" htmlFor="name">
        Form Name
        <input
          className={`form-control ${errors?.name ? 'is-invalid' : ''}`}
          id="name"
          type="text"
          value={formName || ''}
          aria-describedby="nameFeedback"
          {...registerObject}
          onChange={(event) => {
            setFormName(event.target.value);
            registerObject.onChange(event);
          }}
        />
      </label>

      <div className="invalid-feedback" id="nameFeedback">
        <ErrorMessage
          errors={errors}
          name="name"
          // render={({ message }) => <p>{message}</p>}
        />
      </div>
    </div>
  );
};
