import * as React from 'react';
import { FormDefinition } from 'core-api-client';

import { RecordEditorComponent } from './RecordEditor.component';

type Props = {
  form: FormDefinition;
};

export const RecordEditorContainer: React.FC<Props> = ({ form }) => {
  return <RecordEditorComponent form={form} />;
};
