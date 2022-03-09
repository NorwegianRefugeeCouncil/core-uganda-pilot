import * as React from 'react';
import { FormDefinition, Record } from 'core-api-client';

import { RecordEditorComponent } from './RecordEditor.component';

type Props = {
  form: FormDefinition;
  record: Record;
  onChange: () => void;
};

export const RecordEditorContainer: React.FC<Props> = ({
  form,
  record,
  onChange,
}) => {
  return (
    <RecordEditorComponent form={form} record={record} onChange={onChange} />
  );
};
