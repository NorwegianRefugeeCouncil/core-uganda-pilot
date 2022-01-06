import React from 'react';
import { FieldDefinition } from 'core-api-client';

type Props = {
  field: FieldDefinition;
};

export const SelectOptionsList: React.FC<Props> = ({ field }) => {
  const { required, key } = field;

  const options = field.fieldType?.singleSelect?.options;

  if (!options) {
    return <></>;
  }
  return (
    <>
      <option aria-label="no value" disabled={required || key} value="" />
      {options.map((o) => (
        <option key={o.id} value={o.id}>
          {o.name}
        </option>
      ))}
    </>
  );
};
