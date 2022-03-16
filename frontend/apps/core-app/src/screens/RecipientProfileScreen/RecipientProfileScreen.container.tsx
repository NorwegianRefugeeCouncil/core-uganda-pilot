import * as React from 'react';
import { RouteProp, useRoute } from '@react-navigation/native';
import { Recipient } from 'core-api-client/src/types/client/Recipient';
import { FormWithRecord } from 'core-api-client';

import { RootParamList } from '../../navigation/types';
import { formsClient } from '../../clients/formsClient';
import config from '../../config';

import { prettifyData } from './prettifyData';
import { RecipientProfileScreenComponent } from './RecipientProfileScreen.component';

export const RecipientProfileScreenContainer: React.FC = () => {
  const route = useRoute<RouteProp<RootParamList, 'RecipientProfile'>>();
  const [isLoading, setIsLoading] = React.useState<boolean>(true);
  const [data, setData] = React.useState<FormWithRecord<Recipient>[]>([]);
  const [prettifiedData, setPrettifiedData] = React.useState<
    FormWithRecord<Recipient>[]
  >([]);
  const [error, setError] = React.useState<string>();

  React.useEffect(() => {
    (async () => {
      try {
        const recipientData = await formsClient.Recipient.get({
          recordId: route.params.id,
          formId: config.recipient.registrationForm.formId,
          databaseId: config.recipient.registrationForm.databaseId,
        });
        setData(recipientData);
      } catch (err) {
        setError(JSON.stringify(err));
      }
      setIsLoading(false);
    })();
  }, [
    config.recipient.registrationForm.formId,
    config.recipient.registrationForm.databaseId,
    route.params.id,
  ]);

  React.useEffect(() => {
    if (data.length) {
      setPrettifiedData(prettifyData(data));
    }
  }, [data]);

  return (
    <RecipientProfileScreenComponent
      data={prettifiedData}
      isLoading={isLoading}
      error={error}
    />
  );
};
