import * as React from 'react';
import { FormWithRecord } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';
import { StackScreenProps } from '@react-navigation/stack';

import { formsClient } from '../../clients/formsClient';
import configuration from '../../config';
import { linkingConfig } from '../../navigation/linking.config';
import { useAPICall } from '../../hooks/useAPICall';

import { buildDefaultRecord } from './buildDefaultRecord';
import { RecipientRegistrationScreenComponent } from './RecipientRegistrationScreen.component';

// TODO better type
type Props = StackScreenProps<any, any>;

export const RecipientRegistrationScreenContainer: React.FC<Props> = ({
  navigation,
}) => {
  const [mode, setMode] = React.useState<'register' | 'review'>('register');
  const [data, setData] = React.useState<FormWithRecord<Recipient>[]>([]);

  const [_, state] = useAPICall(
    formsClient.Form.getAncestors,
    [configuration.recipient.registrationForm.formId],
    true,
  );

  React.useEffect(() => {
    setData(
      (state.data || []).map((form) => ({
        form,
        record: buildDefaultRecord(form),
      })),
    );
  }, [JSON.stringify(state.data)]);

  const handleSubmit = (d: FormWithRecord<Recipient>[]) => {
    if (mode === 'register') {
      setData(d);
      setMode('review');
    }
    if (mode === 'review') {
      console.log('SAVING DATA', d);
    }
  };

  const handleCancel = () => {
    if (mode === 'register') {
      navigation.navigate(linkingConfig.config.screens.Recipients);
    }
    if (mode === 'review') {
      setMode('register');
    }
  };

  return (
    <RecipientRegistrationScreenComponent
      mode={mode}
      data={data}
      onSubmit={handleSubmit}
      onCancel={handleCancel}
      error={state.error}
      loading={state.loading}
    />
  );
};
