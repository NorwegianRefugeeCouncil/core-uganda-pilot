import React, { FC } from 'react';

import { RecordPickerContainer } from '../RecordPicker';

import { FieldEditorProps } from './types';
import { FieldDescription } from './FieldDescription';
import { FieldLabel } from './FieldLabel';

export const ReferenceFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
}) => {
  if (Array.isArray(value)) {
    return <></>;
  }
  return (
    <div className="form-group mb-2">
      <FieldLabel fieldDefinition={field} />
      <RecordPickerContainer
        formId={field.fieldType.reference?.formId}
        recordId={value}
        setRecordId={onChange}
      />
      <FieldDescription text={field.description} />
    </div>
  );
};
