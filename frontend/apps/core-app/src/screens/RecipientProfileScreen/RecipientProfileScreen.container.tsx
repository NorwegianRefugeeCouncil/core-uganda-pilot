import * as React from 'react';
import { StackScreenProps } from '@react-navigation/stack';

import { formsClient } from '../../clients/formsClient';
import { useAPICall } from '../../hooks/useAPICall';
import { RecipientNavigatorParamList } from '../../navigation/recipients';

import { RecipientProfileScreenComponent } from './RecipientProfileScreen.component';

type Props = StackScreenProps<RecipientNavigatorParamList, 'recipientsProfile'>;

export const RecipientProfileScreenContainer: React.FC<Props> = ({ route }) => {
  const [_, state] = useAPICall(
    formsClient.Recipient.get,
    [
      {
        recordId: route.params.recordId,
        formId: route.params.formId,
        databaseId: route.params.databaseId,
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
