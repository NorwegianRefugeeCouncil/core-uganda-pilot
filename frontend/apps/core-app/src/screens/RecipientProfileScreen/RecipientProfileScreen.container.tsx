import * as React from 'react';
import { RouteProp, useRoute } from '@react-navigation/native';

import { RootParamList } from '../../navigation/types';
import { formsClient } from '../../clients/formsClient';
import config from '../../config';
import { useAPICall } from '../../hooks/useAPICall';

import { RecipientProfileScreenComponent } from './RecipientProfileScreen.component';

export const RecipientProfileScreenContainer: React.FC = () => {
  const route = useRoute<RouteProp<RootParamList, 'RecipientProfile'>>();

  const [_, state] = useAPICall(
    formsClient.Recipient.get,
    [
      {
        recordId: route.params.id,
        formId: config.recipient.registrationForm.formId,
        databaseId: config.recipient.registrationForm.databaseId,
      },
    ],
    true,
  );

  return (
    <RecipientProfileScreenComponent
      data={state.data || []}
      isLoading={state.loading || !state.data || state.data?.length === 0}
      error={state.error}
    />
  );
};
