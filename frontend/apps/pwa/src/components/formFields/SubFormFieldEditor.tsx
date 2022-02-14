import React, { FC } from 'react';

import { FieldEditorProps } from './types';
import { SubRecordList } from './SubRecordList';

export const SubFormFieldEditor: FC<FieldEditorProps> = ({
  field,
  onAddSubRecord,
  onSelectSubRecord,
  subRecords,
}) => {
  return (
    <div className="mb-2">
      <div className="bg-primary border-2" />
      <span className="form-label opacity-75">{field.name}</span>
      {subRecords ? (
        <SubRecordList records={subRecords} select={onSelectSubRecord} />
      ) : (
        <></>
      )}
      <button
        type="button"
        onClick={onAddSubRecord}
        className="btn btn-sm btn-outline-primary w-100"
      >
        Add record in {field.name}
      </button>
    </div>
  );
};
