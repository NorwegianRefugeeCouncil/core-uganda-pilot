import * as React from 'react';
import { NavigationProp, useNavigation } from '@react-navigation/native';
import { TableInstance } from 'react-table';

import { RootParamList } from '../../navigation/types';
import { TableContext } from '../../components/Table/useTableContext';

import { RecipientListScreenComponent } from './RecipientListScreen.component';

export const RecipientListScreenContainer: React.FC = () => {
  const navigation = useNavigation<NavigationProp<RootParamList>>();

  const handleItemClick = (id: string) => {
    navigation.navigate('RecipientProfile', {
      id,
    });
  };

  const [tableInstance, setTableInstance] =
    React.useState<TableInstance | null>(null);

  return (
    <TableContext.Provider value={{ tableInstance, setTableInstance }}>
      <RecipientListScreenComponent handleItemClick={handleItemClick} />
    </TableContext.Provider>
  );
};
