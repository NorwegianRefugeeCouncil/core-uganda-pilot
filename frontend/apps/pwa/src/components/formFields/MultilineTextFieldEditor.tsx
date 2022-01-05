import React, { FC } from 'react';

import { FieldEditorProps } from './types';
import { FieldDescription } from './FieldDescription';
import { FieldLabel } from './FieldLabel';

export const MultilineTextFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
}) => {
  return (
    <div className="form-group mb-2">
      <FieldLabel fieldDefinition={field} />
      <textarea
        className="form-control bg-dark text-light border-secondary"
        id={field.id}
        value={value || ''}
        onChange={(event) => onChange(event.target.value)}
      />
      <FieldDescription text={field.description} />
    </div>
  );
};
