import * as React from 'react';
import { useFormContext } from 'react-hook-form';

import { FieldEditorProps } from './types';

export const CheckboxFieldEditor: React.FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
}) => {
  const { register } = useFormContext();

  const registerObject = register && register(`values.${field.id}`);
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
          {...registerObject}
          onChange={(event) => {
            onChange(event.target.checked ? 'true' : 'false');
            return registerObject.onChange(event);
          }}
          id={field.id}
          aria-describedby={`description-${field.id}`}
        />
        <label className="form-check-label opacity-75" htmlFor={field.id}>
          {field.name}
        </label>
      </div>
    </div>
  );
};
