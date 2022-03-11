import * as React from 'react';
import { FormDefinition } from 'core-api-client';

import { RecordEditorComponent } from './RecordEditor.component';

type Props = {
  form: FormDefinition;
  direction?: 'row' | 'column';
};

export const RecordEditorContainer: React.FC<Props> = ({
  form,
  direction = 'column',
}) => {
  return <RecordEditorComponent form={form} direction={direction} />;
};
