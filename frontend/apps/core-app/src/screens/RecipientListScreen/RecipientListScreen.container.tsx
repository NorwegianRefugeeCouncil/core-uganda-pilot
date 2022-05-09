import * as React from 'react';
import { StackScreenProps } from '@react-navigation/stack';

import { RootNavigatorParamList } from '../../navigation/root';
import { useRecipientForms } from '../../contexts/RecipientForms';

import { RecipientListScreenComponent } from './RecipientListScreen.component';

type Props = StackScreenProps<RootNavigatorParamList, 'recipientsList'>;

export const RecipientListScreenContainer: React.FC<Props> = () => {
  const recipientForms = useRecipientForms();

  return (
    <RecipientListScreenComponent
      forms={recipientForms}
      isLoading={!recipientForms || !recipientForms.length}
    />
  );
};
