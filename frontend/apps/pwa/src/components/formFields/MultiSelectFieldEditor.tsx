import React, { FC } from 'react';

import { FieldEditorProps } from './types';
import {
  mapFieldDescription,
  mapFieldLabel,
  mapSelectOptions,
} from './helpers';

export const MultiSelectFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  setValue,
}) => {
  return (
    <div className="form-group mb-2">
      {mapFieldLabel(field)}
      <select
        className="form-control bg-dark text-light border-secondary"
        id={field.id}
        value={value || []}
        multiple
        onChange={(event) => {
          const { options } = event.target;
          const selected = Object.entries(options).filter((o) => o[1].selected);
          setValue(selected.map((s) => s[1].value));
        }}
      >
        {mapSelectOptions(
          field.required,
          field.key,
          field.fieldType?.multiSelect?.options,
        )}
      </select>
      {mapFieldDescription(field)}
    </div>
  );
};
