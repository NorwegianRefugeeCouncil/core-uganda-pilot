import * as React from 'react';

import { FieldEditorProps } from './types';
import { FieldDescription } from './FieldDescription';

export const CheckboxFieldEditor: React.FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
}) => {
  // It doesn't make sense for this to be null/undefined/empty
  React.useEffect(() => {
    if (value === null || value === undefined || value === '') {
      onChange('false');
    }
  }, [value]);

  return (
    <div className="form-group mb-2">
      <div className="form-check">
        <input
          className="form-check-input"
          type="checkbox"
          checked={value === 'true'}
          onChange={(event) =>
            onChange(event.target.checked ? 'true' : 'false')
          }
          id={field.id}
        />
        <label className="form-check-label opacity-75" htmlFor={field.id}>
          {field.name}
        </label>
      </div>
      <FieldDescription text={field.description} />
    </div>
  );
};
