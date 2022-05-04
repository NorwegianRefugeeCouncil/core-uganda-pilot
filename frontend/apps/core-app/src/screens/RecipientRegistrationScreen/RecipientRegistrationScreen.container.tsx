import * as React from 'react';
import { FormWithRecord } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';
import { StackScreenProps } from '@react-navigation/stack';

import { formsClient } from '../../clients/formsClient';
import { useAPICall } from '../../hooks/useAPICall';
import { routes } from '../../constants/routes';
import { RootNavigatorParamList } from '../../navigation/root';

import { RecipientRegistrationScreenComponent } from './RecipientRegistrationScreen.component';

type Props = StackScreenProps<RootNavigatorParamList, 'recipientsRegistration'>;

export const RecipientRegistrationScreenContainer: React.FC<Props> = ({
  route,
  navigation,
}) => {
  const [mode, setMode] = React.useState<'edit' | 'review'>('edit');
  const [data, setData] = React.useState<FormWithRecord<Recipient>[]>([]);

  const [_, getRecipientFormsState] = useAPICall(
    formsClient.Form.getAncestors,
    [route.params.formId],
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
        recordId:
          saveRecipientState.data[saveRecipientState.data.length - 1].record.id,
        formId: route.params.formId,
        databaseId: route.params.databaseId,
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
