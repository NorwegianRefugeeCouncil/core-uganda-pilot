import { FieldDefinition } from 'core-api-client';
import React from 'react';

type Props = {
  fieldDefinition: FieldDefinition;
};

export const FieldLabel: React.FC<Props> = ({ fieldDefinition }) => {
  return (
    <label className="form-label opacity-75" htmlFor={fieldDefinition.id}>
      {fieldDefinition.name}
    </label>
  );
};
