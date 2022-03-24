import * as React from 'react';
import { NavigationProp, useNavigation } from '@react-navigation/native';

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
      <RecipientListScreenComponent handleItemClick={handleItemClick} />
    </RecordTableContext.Provider>
  );
};
