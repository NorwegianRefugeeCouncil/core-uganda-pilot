import React, { FC } from 'react';

import { RecordPickerContainer } from '../RecordPicker';

import { FieldEditorProps } from './types';

export const ReferenceFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
  errors,
}) => {
  if (Array.isArray(value)) {
    return <></>;
  }
  return (
    <RecordPickerContainer
      formId={field.fieldType.reference?.formId}
      field={field}
      recordId={value}
      setRecordId={onChange}
      errors={errors}
    />
  );
};
