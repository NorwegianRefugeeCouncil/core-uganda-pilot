import { Box, HStack, Pressable, Text } from 'native-base';
import { Cell, Row } from 'react-table';
import React from 'react';

import { RecipientListTableEntry } from './types';

type Props<T extends Record<string, any>> = {
  row: Row<T>;
  onRowClick: (id: string) => void;
};

export const RecipientListTableRow: React.FC<
  Props<RecipientListTableEntry>
> = ({ row, onRowClick }) => {
  const handleRowClick = () => onRowClick(row.id);

  return (
    <Pressable
      onPress={handleRowClick}
      bg="white"
      _hover={{ backgroundColor: 'primary.100' }}
      testID={`recipient-list-table-row-${row.id}`}
    >
      <HStack>
        {row.cells.map((cell: Cell, i) => {
          return (
            <Box
              key={i}
              p={2}
              width={cell.column.width}
              borderBottomColor="neutral.200"
              borderBottomWidth="1"
              flexGrow={1}
            >
              {cell.render(Text, {
                variant: 'body',
                level: '2',
                children: cell.value,
              })}
            </Box>
          );
        })}
      </HStack>
    </Pressable>
  );
};
