import React, { FC } from 'react';

import { RecordPickerContainer } from '../RecordPicker';

import { FieldEditorProps } from './types';
import { mapFieldDescription, mapFieldLabel } from './helpers';

export const ReferenceFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  setValue,
}) => {
  if (Array.isArray(value)) {
    return <></>;
  }
  return (
    <div className="form-group mb-2">
      {mapFieldLabel(field)}
      <RecordPickerContainer
        formId={field.fieldType.reference?.formId}
        recordId={value}
        setRecordId={setValue}
      />
      {mapFieldDescription(field)}
    </div>
  );
};
