import React, { useContext } from 'react';
import { Box, HStack, VStack, FlatList, ScrollView } from 'native-base';
import { Row } from 'react-table';

import { RecipientListTableContext } from './RecipientListTableContext';
import { RecipientListTableRow } from './RecipientListTableRow';
import { RecipientListTableHeaderCell } from './RecipientListTableHeaderCell';

type Props = {
  onItemClick: (id: string) => void;
};

export const RecipientListTableComponent: React.FC<Props> = ({
  onItemClick,
}) => {
  const tableContext = useContext(RecipientListTableContext);
  if (!tableContext) return null;

  const { tableInstance } = tableContext;
  if (!tableInstance) return null;

  const { rows, columns, prepareRow } = tableInstance;

  const renderRow = ({ item }: { item: Row }) => {
    prepareRow(item);
    return (
      <RecipientListTableRow
        key={item.id}
        row={item}
        onRowClick={onItemClick}
      />
    );
  };

  return (
    <ScrollView
      horizontal
      contentContainerStyle={{
        flexGrow: 1,
      }}
    >
      <Box width="100%">
        <HStack>
          {columns.map((column) => (
            <RecipientListTableHeaderCell column={column} key={column.id} />
          ))}
        </HStack>
        <VStack>
          <FlatList data={rows} renderItem={renderRow} />
        </VStack>
      </Box>
    </ScrollView>
  );
};
