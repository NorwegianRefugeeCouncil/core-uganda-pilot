import * as React from 'react';
import { TableInstance, Row, Cell } from 'react-table';
import { Box, HStack, VStack, Text, Button, ScrollView } from 'native-base';
import { Icon } from 'core-design-system';

type Props = {
  table: TableInstance<any>;
  onDelete?: (idx: number) => void;
};

export const SubFormTableComponent: React.FC<Props> = ({ table, onDelete }) => {
  const handleDelete = (idx: number) => () => onDelete?.(idx);

  return (
    <ScrollView width="100%" horizontal>
      <VStack width="100%">
        <HStack width="100%">
          {table.columns.map((col) => (
            <Box
              key={col.id}
              width={col.width}
              p="2"
              borderBottomColor="neutral.300"
              borderBottomWidth="1"
              // flexGrow={1}
            >
              {col.render(Text, {
                variant: 'body',
                fontWeight: '700',
                children: col.Header,
              })}
            </Box>
          ))}
        </HStack>
        <VStack>
          {table.rows.map((row: Row, i) => {
            table.prepareRow(row);
            return (
              <HStack key={row.id} width="100%">
                {row.cells.map((cell: Cell, j) => (
                  <Box
                    key={`${i}-${j}`}
                    p={2}
                    // flexGrow={1}
                    width={cell.column.width}
                    borderBottomColor="neutral.300"
                    borderBottomWidth={i === table.rows.length - 1 ? 0 : 1}
                  >
                    {cell.render(Text, {
                      level: '2',
                      children: cell.value,
                    })}
                  </Box>
                ))}
                {onDelete && (
                  <Box
                    p={2}
                    // flexGrow={1}
                    borderBottomColor="neutral.300"
                    borderBottomWidth={i === table.rows.length - 1 ? 0 : 1}
                  >
                    <Button
                      onPress={handleDelete(i)}
                      colorScheme="secondary"
                      variant="naked"
                      startIcon={<Icon name="delete" />}
                    />
                  </Box>
                )}
              </HStack>
            );
          })}
        </VStack>
      </VStack>
    </ScrollView>
  );
};
