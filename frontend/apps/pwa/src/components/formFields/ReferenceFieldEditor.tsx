import React, { FC } from 'react';

import { RecordPickerContainer } from '../RecordPicker';
import { FieldDefinitionNC } from '../../reducers/Former/types';

import { FieldEditorProps } from './types';
import { FieldDescription } from './FieldDescription';
import { FieldLabel } from './FieldLabel';

export const ReferenceFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
  register,
  errors,
}) => {
  if (Array.isArray(value)) {
    return <></>;
  }
  return (
    <div className="form-group mb-2">
      <FieldLabel fieldDefinition={field} />
      <RecordPickerContainer
        formId={field.fieldType.reference?.formId}
        field={field as FieldDefinitionNC}
        recordId={value}
        setRecordId={onChange}
        register={register}
        errors={errors}
      />
      <FieldDescription text={field.description} />
    </div>
  );
};
