import * as React from 'react';
import { FormDefinition } from 'core-api-client';
import { useFormContext } from 'react-hook-form';

import { RecordEditorComponent } from './RecordEditor.component';

type Props = {
  form: FormDefinition;
  hideKeyFields?: boolean;
};

export const RecordEditorContainer: React.FC<Props> = ({
  form,
  hideKeyFields,
}) => {
  const formContext = useFormContext();

  if (!formContext) throw new Error('Form context is not available');

  return <RecordEditorComponent form={form} hideKeyFields={hideKeyFields} />;
};
