import React, { FC } from 'react';

import { FieldEditorProps } from './types';
import { FieldDescription } from './FieldDescription';
import { FieldLabel } from './FieldLabel';

export const DateFieldEditor: FC<FieldEditorProps> = (props) => {
  const { field, value, onChange } = props;
  return (
    <div className="form-group mb-2">
      <FieldLabel fieldDefinition={field} />
      <input
        className="form-control bg-dark text-light border-secondary"
        type="date"
        id={field.id}
        value={value || ''}
        onChange={(event) => onChange(event.target.value)}
      />
      <FieldDescription text={field.description} />
    </div>
  );
};
