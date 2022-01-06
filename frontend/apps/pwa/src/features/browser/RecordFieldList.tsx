import { FieldDefinition, Record, SelectOption } from 'core-api-client';
import React from 'react';

import { RecordField } from './RecordField';

type Props = {
  record: Record;
  field: FieldDefinition;
  subRecords: Record[] | undefined;
};

export const RecordFieldList: React.FC<Props> = ({
  record,
  field,
  subRecords,
}) => {
  let value = '';
  const fieldValue = record.values.find((v: any) => v.fieldId === field.id);
  if (fieldValue && typeof fieldValue.value === 'string') {
    value = fieldValue.value;
  }
  if (fieldValue && field.fieldType.multiSelect) {
    const selected = field.fieldType.multiSelect.options.filter(
      (o: SelectOption) => {
        if (fieldValue?.value == null) {
          return false;
        }
        return fieldValue.value.includes(o.id);
      },
    );

    value = selected.map((s) => s.name).join(', ');
  }
  return (
    <RecordField
      key={record.id}
      field={field}
      value={`${value}`}
      subRecords={subRecords}
    />
  );
};
