import React, { useContext } from 'react';
import { Box, HStack, Text, VStack, FlatList } from 'native-base';

import { RecordTableContext } from './RecordTableContext';
import { RecordTableRow } from './RecordTableRow';
import { RecordTableHeaderCell } from './RecordTableHeaderCell';

type Props = {
  onItemClick: (id: string) => void;
};

export const RecordTableComponent: React.FC<Props> = ({ onItemClick }) => {
  const tableContext = useContext(RecordTableContext);
  if (!tableContext) return null;

  const { tableInstance } = tableContext;
  if (!tableInstance) return null;

  const { rows, columns, prepareRow, globalFilteredRows } = tableInstance;

  return (
    <Box maxWidth="100%" overflowX="scroll">
      <Text level="2">{globalFilteredRows.length} beneficiaries</Text>
      <HStack>
        {columns.map((column) => (
          <RecordTableHeaderCell column={column} key={column.id} />
        ))}
      </HStack>
      <VStack>
        <FlatList
          data={rows}
          renderItem={({ item }) => {
            prepareRow(item);
            return (
              <RecordTableRow
                key={item.id}
                row={item}
                onRowClick={onItemClick}
              />
            );
          }}
        />
      </VStack>
    </Box>
  );
};
