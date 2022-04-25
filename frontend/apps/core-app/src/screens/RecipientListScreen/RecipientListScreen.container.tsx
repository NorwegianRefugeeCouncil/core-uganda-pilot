import * as React from 'react';
import { StackScreenProps } from '@react-navigation/stack';

import { RecipientListTableContext } from '../../components/RecipientListTable/RecipientListTableContext';
import {
  RecipientListTableEntry,
  SortedFilteredTable,
} from '../../components/RecipientListTable/types';
import { RootNavigatorParamList } from '../../navigation/root';
import { routes } from '../../constants/routes';
import { useRecipientForms } from '../../contexts/RecipientForms';
import { formsClient } from '../../clients/formsClient';
import { useAPICall } from '../../hooks/useAPICall';

import { RecipientListScreenComponent } from './RecipientListScreen.component';

type Props = StackScreenProps<RootNavigatorParamList, 'recipientsList'>;

export const RecipientListScreenContainer: React.FC<Props> = ({
  navigation,
  route,
}) => {
  const recipientForms = useRecipientForms();

  const handleItemClick = (id: string) => {
    navigation.navigate(routes.recipientsProfile.name, {
      recordId: id,
      formId: recipientForms[0]?.id,
      databaseId: recipientForms[0]?.databaseId,
    });
  };

  const [tableInstance, setTableInstance] =
    React.useState<SortedFilteredTable<RecipientListTableEntry> | null>(null);

  const [_, formState] = useAPICall(
    formsClient.Form.getAncestors,
    [route.params.formId],
    true,
  );

  const [_, recipientsState] = useAPICall(
    formsClient.Recipient.list,
    [
      {
        formId: route.params.formId,
        databaseId: route.params.databaseId,
      },
    ],
    true,
  );

  console.log('recipientsState.data', recipientsState.data, formState);

  return (
    <RecipientListTableContext.Provider
      value={{ tableInstance, setTableInstance }}
    >
      <RecipientListScreenComponent
        onItemClick={handleItemClick}
        data={recipientsState.data}
        form={formState.data}
        isLoading={formState.loading || recipientsState.loading}
        error={undefined}
      />
    </RecipientListTableContext.Provider>
  );
};
