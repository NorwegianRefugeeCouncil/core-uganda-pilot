import React, { ChangeEvent, FC } from 'react';

import { FieldEditorProps } from './types';
import { FieldDescription } from './FieldDescription';
import { FieldLabel } from './FieldLabel';

export const MonthFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
}) => {
  const expectedLength = 7;

  const isValid = (s: string) => {
    const valid = /^(?:19|20|21)\d{2}-[01]\d$/;
    const m = +s.slice(5);
    return valid.test(s) && m > 0 && m <= 12;
  };

  const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
    const v = event.target.value;
    if (!isValid(v)) return;
    onChange(v);
  };

  if (Array.isArray(value)) {
    return <></>;
  }
  return (
    <div className="form-group mb-2">
      <FieldLabel fieldDefinition={field} />
      <input
        className="form-control bg-dark text-light border-secondary"
        type="month"
        maxLength={expectedLength}
        id={field.id}
        value={value || ''}
        name={field.name}
        pattern="[0-9]{4}-[0-9]{2}"
        placeholder="YYYY-MM"
        onChange={handleChange}
      />
      <FieldDescription text={field.description} />
    </div>
  );
};
