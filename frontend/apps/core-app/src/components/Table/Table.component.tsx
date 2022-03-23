import React, { useContext } from 'react';
import { Box, Button, HStack, Pressable, Text, VStack } from 'native-base';
import { Icon } from 'core-design-system';
import { Cell, ColumnInstance, Row, UseSortByColumnProps } from 'react-table';

import { TableContext } from './useTableContext';

type Props = {
  handleItemClick: (id: string) => void;
};

// Let the table remove the filter if the string is empty
export const TableComponent: React.FC<Props> = ({ handleItemClick }) => {
  const tableContext = useContext(TableContext);
  if (!tableContext) return null;

  const { tableInstance } = tableContext;
  if (!tableInstance) return null;

  const { rows, headers, prepareRow } = tableInstance;

  return (
    <Box p={2}>
      <HStack space={2}>
        {headers.map((column: ColumnInstance & UseSortByColumnProps<any>) => {
          const { onClick } = column.getSortByToggleProps();
          const sortIcon = column.isSorted
            ? column.isSortedDesc
              ? 'more'
              : 'next'
            : 'plus';
          return (
            <Box key={column.id} width={column.width}>
              <HStack
                p="2"
                borderColor="primary.500"
                borderWidth="1"
                width="100%"
                alignItems="center"
              >
                {column.render('Header')}
                {onClick && (
                  <Pressable onPress={(e) => onClick(e)}>
                    <Icon size="3" viewBox="10 10 20 20" name={sortIcon} />
                  </Pressable>
                )}
              </HStack>
            </Box>
          );
        })}
      </HStack>
      <VStack space={2}>
        {rows.map((row: Row) => {
          prepareRow(row);
          return (
            <Pressable key={row.id} onPress={() => handleItemClick(row.id)}>
              <HStack space={2}>
                {row.cells.map((cell: Cell, i) => {
                  return (
                    <Box key={i} p={2} width={cell.column.width}>
                      {i === 0 &&
                        cell.render(Button, {
                          variant: 'minor',
                          children: cell.value,
                        })}
                      {i !== 0 &&
                        cell.render(Text, {
                          variant: 'body',
                          color: 'primary.500',
                          children: cell.value,
                        })}
                    </Box>
                  );
                })}
              </HStack>
            </Pressable>
          );
        })}
      </VStack>
    </Box>
  );
};
