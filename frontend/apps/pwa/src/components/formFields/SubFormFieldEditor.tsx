import React, { FC } from 'react';

import { mapFieldDescription, mapSubRecords } from './helpers';
import { FieldEditorProps } from './types';

export const SubFormFieldEditor: FC<FieldEditorProps> = ({
  field,
  addSubRecord,
  selectSubRecord,
  subRecords,
}) => {
  return (
    <div className="mb-2">
      <div className="bg-primary border-2" />
      <span className="form-label opacity-75">{field.name}</span>
      {subRecords ? mapSubRecords(subRecords, selectSubRecord) : <></>}
      <button
        type="button"
        onClick={addSubRecord}
        className="btn btn-sm btn-outline-primary w-100"
      >
        Add record in {field.name}
      </button>
      {mapFieldDescription(field)}
    </div>
  );
};
