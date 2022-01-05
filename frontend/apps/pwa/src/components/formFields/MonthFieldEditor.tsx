import React, { FC } from 'react';

import { mapFieldDescription, mapFieldLabel } from './helpers';
import { FieldEditorProps } from './types';

export const MonthFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  setValue,
}) => {
  const expectedLength = 7;

  function isValid(s: string) {
    const valid = /^(?:19|20|21)\d{2}-[01]\d$/;
    const m = +s.slice(5);
    return valid.test(s) && m > 0 && m <= 12;
  }

  if (Array.isArray(value)) {
    return <></>;
  }
  return (
    <div className="form-group mb-2">
      {mapFieldLabel(field)}
      <input
        className="form-control bg-dark text-light border-secondary"
        type="month"
        maxLength={expectedLength}
        id={field.id}
        value={value || ''}
        name={field.name}
        pattern="[0-9]{4}-[0-9]{2}"
        placeholder="YYYY-MM"
        onChange={(event) => {
          const v = event.target.value;
          if (!isValid(v)) return;
          setValue(v);
        }}
      />
      {mapFieldDescription(field)}
    </div>
  );
};
