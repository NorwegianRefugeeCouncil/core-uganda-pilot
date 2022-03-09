import * as React from 'react';
import { FormDefinition, Record } from 'core-api-client';

import { RecordViewComponent } from './RecordView.component';
import { normaliseFieldValues } from './normaliseFieldValues';

type Props = {
  form: FormDefinition;
  record: Record;
};

export const RecordViewContainer: React.FC<Props> = ({ form, record }) => {
  const fieldValues = normaliseFieldValues(form, record);

  return <RecordViewComponent fieldValues={fieldValues} />;
};
