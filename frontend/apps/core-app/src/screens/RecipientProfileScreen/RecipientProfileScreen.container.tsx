import * as React from 'react';
import { RouteProp, useRoute } from '@react-navigation/native';
import { Recipient } from 'core-api-client/src/types/client/Recipient';
import { FormDefinition } from 'core-api-client';

import { RootParamList } from '../../navigation/types';
import { formsClient } from '../../clients/formsClient';

import { RecipientProfileScreenComponent } from './RecipientProfileScreen.component';

const DATABASE_ID = '5b5e0630-3985-4c04-a4ac-2ca958421ca7';
const FORM_ID = '989533fa-56e7-4771-8d68-c8c4a707c240';

export const RecipientProfileScreenContainer: React.FC = () => {
  const route = useRoute<RouteProp<RootParamList, 'RecipientProfile'>>();
  const [isLoading, setIsLoading] = React.useState<boolean>(true);
  const [form, setForm] = React.useState<FormDefinition | null>(null);
  const [recipient, setRecipient] = React.useState<Recipient | null>(null);
  const [error, setError] = React.useState<any>();

  React.useEffect(() => {
    (async () => {
      const formResponse = await formsClient.Form.get({ id: FORM_ID });
      if (formResponse.response) setForm(formResponse.response);

      const recipientGetResponse = await formsClient.Recipient.get({
        recordId: route.params.id,
        formId: FORM_ID,
        databaseId: DATABASE_ID,
      });
      if (recipientGetResponse.response) {
        setRecipient(recipientGetResponse.response);
      } else {
        setError(recipientGetResponse.error);
      }
      setIsLoading(false);
    })();
  }, [FORM_ID, DATABASE_ID, route.params.id]);

  return (
    <RecipientProfileScreenComponent
      recipient={recipient}
      form={form}
      isLoading={isLoading}
      error={error}
    />
  );
};
