import React, { FC } from 'react';

import { FieldEditorProps } from './types';
import { FieldDescription } from './FieldDescription';
import { FieldLabel } from './FieldLabel';

export const TextFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
}) => {
  return (
    <div className="form-group mb-2">
      <FieldLabel fieldDefinition={field} />
      <input
        className="form-control bg-dark text-light border-secondary"
        type="text"
        id={field.id}
        value={value || ''}
        onChange={(event) => onChange(event.target.value)}
      />
      <FieldDescription text={field.description} />
    </div>
  );
};
