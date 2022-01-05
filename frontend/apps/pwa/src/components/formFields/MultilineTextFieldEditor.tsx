import React, { FC } from 'react';

import { FieldEditorProps } from './types';
import { mapFieldDescription, mapFieldLabel } from './helpers';

export const MultilineTextFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  setValue,
}) => {
  return (
    <div className="form-group mb-2">
      {mapFieldLabel(field)}
      <textarea
        className="form-control bg-dark text-light border-secondary"
        id={field.id}
        value={value || ''}
        onChange={(event) => setValue(event.target.value)}
      />
      {mapFieldDescription(field)}
    </div>
  );
};
