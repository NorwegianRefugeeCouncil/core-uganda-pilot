import React, { FC } from 'react';

import { FieldEditorProps } from './types';
import { FieldDescription } from './FieldDescription';
import { FieldLabel } from './FieldLabel';

export const WeekFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
}) => {
  const isValidWeek = (weekString: string) => {
    const weekRegex = /^(?:19|20|21)\d{2}-W[0-5]\d$/;
    return weekRegex.test(weekString) && +weekString.slice(6) <= 52;
  };

  const onChangeHandler = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (!isValidWeek(event.target.value)) return;
    onChange(event.target.value);
  };

  return (
    <div className="form-group mb-2">
      <FieldLabel fieldDefinition={field} />
      <input
        className="form-control bg-dark text-light border-secondary"
        type="week"
        name={field.name}
        maxLength={8}
        placeholder="2021-W52"
        id={field.id}
        value={value || ''}
        onChange={onChangeHandler}
      />
      <FieldDescription text={field.description} />
    </div>
  );
};
