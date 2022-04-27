import * as React from 'react';
import { FormWithRecord } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';
import { useNavigation } from '@react-navigation/native';

import { formsClient } from '../../clients/formsClient';
import configuration from '../../config';
import { useAPICall } from '../../hooks/useAPICall';
import { routes } from '../../constants/routes';

import { RecipientRegistrationScreenComponent } from './RecipientRegistrationScreen.component';

export const RecipientRegistrationScreenContainer: React.FC = () => {
  const navigation = useNavigation();

  const [mode, setMode] = React.useState<'edit' | 'review'>('edit');
  const [data, setData] = React.useState<FormWithRecord<Recipient>[]>([]);

  const [_, getRecipientFormsState] = useAPICall(
    formsClient.Form.getAncestors,
    [configuration.recipient.registrationForm.formId],
    true,
  );

  const [saveRecipient, saveRecipientState] = useAPICall(
    formsClient.Recipient.create,
    [data],
    false,
  );

  React.useEffect(() => {
    const unsubscribe = navigation.addListener('blur', () => {
      setMode('edit');
      setData([]);
    });

    return unsubscribe;
  }, [navigation]);

  React.useEffect(() => {
    setData(
      (getRecipientFormsState.data || []).map((form) => ({
        form,
        record: formsClient.Record.buildDefaultRecord(form),
      })),
    );
  }, [JSON.stringify(getRecipientFormsState.data)]);

  React.useEffect(() => {
    if (saveRecipientState.data) {
      navigation.navigate(routes.recipientsProfile.name, {
        id: saveRecipientState.data[saveRecipientState.data.length - 1].record
          .id,
      });
    }
  }, [JSON.stringify(saveRecipientState.data)]);

  const handleSubmit = (d: FormWithRecord<Recipient>[]) => {
    if (mode === 'edit') {
      setData(d);
      setMode('review');
    }
    if (mode === 'review') {
      saveRecipient();
    }
  };

  const handleCancel = () => {
    if (mode === 'edit') {
      navigation.navigate(routes.recipientsList.name);
    }
    if (mode === 'review') {
      setMode('edit');
    }
  };

  return (
    <RecipientRegistrationScreenComponent
      mode={mode}
      data={data}
      onSubmit={handleSubmit}
      onCancel={handleCancel}
      error={getRecipientFormsState.error}
      loading={getRecipientFormsState.loading}
    />
  );
};
