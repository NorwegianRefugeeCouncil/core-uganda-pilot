import React, { ChangeEvent, FC } from 'react';

import { FieldEditorProps } from './types';
import { FieldDescription } from './FieldDescription';
import { FieldLabel } from './FieldLabel';
import { SelectOptionsList } from './SelectOptionsList';

export const MultiSelectFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
}) => {
  const handleChange = (event: ChangeEvent<HTMLSelectElement>) => {
    const { options } = event.target;
    const selected = Object.entries(options).filter((o) => o[1].selected);
    onChange(selected.map((s) => s[1].value));
  };

  return (
    <div className="form-group mb-2">
      <FieldLabel fieldDefinition={field} />
      <select
        className="form-control bg-dark text-light border-secondary"
        id={field.id}
        value={value || []}
        multiple
        onChange={handleChange}
      >
        <SelectOptionsList field={field} isMultiSelect />
      </select>
      <FieldDescription text={field.description} />
    </div>
  );
};
