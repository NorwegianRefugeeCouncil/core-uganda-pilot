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

import { RecipientListScreenComponent } from './RecipientListScreen.component';

type Props = StackScreenProps<RootNavigatorParamList, 'recipientsList'>;

export const RecipientListScreenContainer: React.FC<Props> = ({
  navigation,
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

  return (
    <RecipientListTableContext.Provider
      value={{ tableInstance, setTableInstance }}
    >
      <RecipientListScreenComponent onItemClick={handleItemClick} />
    </RecipientListTableContext.Provider>
  );
};
