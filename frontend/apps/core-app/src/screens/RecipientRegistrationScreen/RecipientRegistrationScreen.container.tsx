import * as React from 'react';
import { FormDefinition, Record } from 'core-api-client';

import { buildDefaultRecord } from './buildDefaultRecord';
import { RecipientRegistrationScreenComponent } from './RecipientRegistrationScreen.component';

const getForms = async (): Promise<FormDefinition[]> => Promise.resolve([]);

export const RecipientRegistrationScreenContainer: React.FC = () => {
  const [mode, setMode] = React.useState<'register' | 'review'>('register');
  const [forms, setForms] = React.useState<FormDefinition[]>([]);
  const [records, setRecords] = React.useState<Record[]>([]);

  React.useEffect(() => {
    (async () => {
      const formsResponse = await getForms();
      setForms(formsResponse);
    })();
  }, []);

  React.useEffect(() => {
    setRecords(forms.map(buildDefaultRecord));
    setMode('register');
  }, [JSON.stringify(forms)]);

  if (mode === 'register') {
    return (
      <RecipientRegistrationScreenComponent
        forms={forms}
        records={records}
        onSubmit={() => {}}
        onCancel={() => {}}
      />
    );
  }

  if (mode === 'review') {
    return null;
  }

  return null;
};
