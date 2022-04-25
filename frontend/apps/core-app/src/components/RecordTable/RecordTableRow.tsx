import { Box, HStack, Pressable, Text } from 'native-base';
import { Cell, Row } from 'react-table';
import React from 'react';

import { RecordTableEntry } from './types';

type Props<T extends Record<string, any>> = {
  row: Row<T>;
  onRowClick: (id: string) => void;
};

export const RecordTableRow: React.FC<Props<RecordTableEntry>> = ({
  row,
  onRowClick,
}) => {
  return (
    <Pressable onPress={() => onRowClick(row.id)}>
      <HStack>
        {row.cells.map((cell: Cell, i) => {
          return (
            <Box
              key={i}
              p={2}
              width={cell.column.width}
              borderBottomColor="neutral.300"
              borderBottomWidth="1"
              bg="white"
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
