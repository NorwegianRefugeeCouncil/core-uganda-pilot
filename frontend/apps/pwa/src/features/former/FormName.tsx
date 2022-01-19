import React, { FC } from 'react';

import { FormerProps } from './types';

export const FormName: FC<
  Pick<FormerProps, 'formName' | 'setFormName' | 'errors'>
> = ({ formName = '', setFormName, errors }) => {
  return (
    <div className="form-group mb-2">
      <label className="form-label" htmlFor="formName">
        Form Name
      </label>
      <input
        className={`form-control ${errors ? 'is-invalid' : ''}`}
        id="formName"
        type="text"
        value={formName || ''}
        onChange={(event) => setFormName(event.target.value)}
        aria-describedby="formNameFeedback"
      />
      {errors && (
        <div className="invalid-feedback is-invalid" id="formNameFeedback">
          {Object.values(errors)?.map((error) => (
            <div key={error}>{error}</div>
          ))}
        </div>
      )}
    </div>
  );
};
