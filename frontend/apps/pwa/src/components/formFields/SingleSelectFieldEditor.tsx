import React, { FC } from 'react';

import { FieldEditorProps } from './types';
import { FieldDescription } from './FieldDescription';
import { FieldLabel } from './FieldLabel';
import { SelectOptionsList } from './SelectOptionsList';

export const SingleSelectFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
}) => {
  return (
    <div className="form-group mb-2">
      <FieldLabel fieldDefinition={field} />
      <select
        className="form-control bg-dark text-light border-secondary"
        id={field.id}
        value={value || ''}
        onChange={(event) => onChange(event.target.value)}
      >
        <SelectOptionsList field={field} />
      </select>
      <FieldDescription text={field.description} />
    </div>
  );
};
