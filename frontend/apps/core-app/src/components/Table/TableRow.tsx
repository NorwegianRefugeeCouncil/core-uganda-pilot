import { Box, HStack, Pressable, Text } from 'native-base';
import { Cell, Row } from 'react-table';
import React from 'react';

type Props = {
  row: Row;
  handleRowClick: () => void;
};

export const TableRow: React.FC<Props> = ({ row, handleRowClick }) => {
  return (
    <Pressable onPress={() => handleRowClick()}>
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
              {i === 0 &&
                cell.render(Text, {
                  level: '2',
                  children: cell.value,
                })}
              {i !== 0 &&
                cell.render(Text, {
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
