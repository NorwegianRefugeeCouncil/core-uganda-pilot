import * as React from 'react';
import { NavigationProp, useNavigation } from '@react-navigation/native';
import { Button } from 'native-base';

import { RootParamList } from '../../navigation/types';
import { RecordTableContext } from '../../components/RecordTable/RecordTableContext';
import {
  RecordTableEntry,
  SortedFilteredTable,
} from '../../components/RecordTable/types';

import { RecipientListScreenComponent } from './RecipientListScreen.component';

export const RecipientListScreenContainer: React.FC = () => {
  const navigation = useNavigation<NavigationProp<RootParamList>>();

  const handleItemClick = (id: string) => {
    navigation.navigate('RecipientProfile', {
      id,
    });
  };

  const [tableInstance, setTableInstance] =
    React.useState<SortedFilteredTable<RecordTableEntry> | null>(null);

  return (
    <RecordTableContext.Provider value={{ tableInstance, setTableInstance }}>
      <Button
        variant="major"
        color="primary"
        onPress={() => navigation.navigate('RecipientRegistration', {})}
      >
        Register
      </Button>
      <RecipientListScreenComponent onItemClick={handleItemClick} />
    </RecordTableContext.Provider>
  );
};
