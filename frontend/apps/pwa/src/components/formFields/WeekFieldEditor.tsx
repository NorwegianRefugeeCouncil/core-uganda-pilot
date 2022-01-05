import React, { FC } from 'react';

import { FieldEditorProps } from './types';
import { mapFieldDescription, mapFieldLabel } from './helpers';

export const WeekFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  setValue,
}) => {
  function isValidWeek(weekString: string) {
    const weekRegex = /^(?:19|20|21)\d{2}-W[0-5]\d$/;
    return weekRegex.test(weekString) && +weekString.slice(6) <= 52;
  }

  function onChangeHandler(event: React.ChangeEvent<HTMLInputElement>) {
    if (!isValidWeek(event.target.value)) return;
    setValue(event.target.value);
  }

  return (
    <div className="form-group mb-2">
      {mapFieldLabel(field)}
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
      {mapFieldDescription(field)}
    </div>
  );
};
