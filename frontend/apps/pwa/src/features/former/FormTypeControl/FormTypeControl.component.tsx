import { FormType } from 'core-api-client';
import * as React from 'react';

type Props = {
  formType: FormType;
  onFormTypeChange: (type: FormType) => void;
};

export const FormTypeControlComponent: React.FC<Props> = ({
  formType,
  onFormTypeChange,
}) => {
  const handleOnChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    onFormTypeChange(event.target.value as FormType);
  };

  return (
    <div className="form-group mb-2">
      <label className="form-label" htmlFor="form-type-control">
        Form type
      </label>
      <select
        className="form-control bg-dark text-light border-secondary"
        id="form-type-control"
        value={formType}
        onChange={handleOnChange}
      >
        <option value={FormType.DefaultFormType}>
          {FormType.DefaultFormType}
        </option>
        <option value={FormType.RecipientFormType}>
          {FormType.RecipientFormType}
        </option>
      </select>
    </div>
  );
};
