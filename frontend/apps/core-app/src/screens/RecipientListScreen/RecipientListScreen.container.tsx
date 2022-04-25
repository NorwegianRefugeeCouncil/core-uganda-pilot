import * as React from 'react';
import { NavigationProp, useNavigation } from '@react-navigation/native';
import { Button } from 'native-base';

import { RootParamList } from '../../navigation/types';
import { RecipientListTableContext } from '../../components/RecipientListTable/RecipientListTableContext';
import {
  RecipientListTableEntry,
  SortedFilteredTable,
} from '../../components/RecipientListTable/types';

import { RecipientListScreenComponent } from './RecipientListScreen.component';

export const RecipientListScreenContainer: React.FC = () => {
  const navigation = useNavigation<NavigationProp<RootParamList>>();

  const handleItemClick = (id: string) => {
    navigation.navigate('RecipientProfile', {
      id,
    });
  };

  const [tableInstance, setTableInstance] =
    React.useState<SortedFilteredTable<RecipientListTableEntry> | null>(null);

  return (
    <RecipientListTableContext.Provider
      value={{ tableInstance, setTableInstance }}
    >
      <Button
        variant="major"
        color="primary"
        onPress={() => navigation.navigate('RecipientRegistration', {})}
      >
        Register
      </Button>
      <RecipientListScreenComponent onItemClick={handleItemClick} />
    </RecipientListTableContext.Provider>
  );
};
