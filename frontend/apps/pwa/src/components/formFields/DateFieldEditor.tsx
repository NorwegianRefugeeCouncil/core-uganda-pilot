import React, { FC } from 'react';

import { FieldEditorProps } from './types';
import { mapFieldDescription, mapFieldLabel } from './helpers';

export const DateFieldEditor: FC<FieldEditorProps> = (props) => {
  const { field, value, setValue } = props;
  return (
    <div className="form-group mb-2">
      {mapFieldLabel(field)}
      <input
        className="form-control bg-dark text-light border-secondary"
        type="date"
        id={field.id}
        value={value || ''}
        onChange={(event) => setValue(event.target.value)}
      />
      {mapFieldDescription(field)}
    </div>
  );
};
