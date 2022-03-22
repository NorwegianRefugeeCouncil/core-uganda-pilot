import React from 'react';
import {
  // useAsyncDebounce,
  useFilters,
  useGlobalFilter,
  useSortBy,
  useTable,
  // Column
} from 'react-table';
import { Box, Button, HStack, Pressable, Text, VStack } from 'native-base';
import { FieldValue, FormDefinition } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';
import { Icon } from 'core-design-system';

type Props = {
  records: Recipient[];
  handleItemClick: (id: string) => void;
  form: FormDefinition;
  searchTerm: string;
};

// Let the table remove the filter if the string is empty
export const TableComponent: React.FC<Props> = ({
  records,
  handleItemClick,
  form,
  searchTerm,
}) => {
  const memoizedData = React.useMemo(
    () =>
      records.map((record) => {
        return record.values.reduce((acc, value: FieldValue) => {
          const field = form.fields.find((f) => {
            return f.id === value.fieldId;
          });
          if (field) return { ...acc, [field?.id]: value.value };
          return acc;
        }, {});
      }),
    [records],
  );
  const memoizedColumns = React.useMemo(
    () =>
      form.fields.map((field) => {
        return {
          Header: field.name,
          accessor: field.id,
        };
      }),
    [form],
  );

  const filterTypes = React.useMemo(
    () => ({
      // Add a new fuzzyTextFilterFn filter type.
      fuzzyText: fuzzyTextFilterFn,
      // Or, override the default text filter to use
      // "startWith"
      text: (rows, id, filterValue) => {
        return rows.filter((row) => {
          const rowValue = row.values[id];
          return rowValue !== undefined
            ? String(rowValue)
                .toLowerCase()
                .startsWith(String(filterValue).toLowerCase())
            : true;
        });
      },
    }),
    [],
  );
  const defaultColumn = React.useMemo(
    () => ({
      // Let's set up our default Filter UI
      Filter: DefaultColumnFilter,
    }),
    [],
  );

  const { rows, headerGroups, prepareRow } = useTable(
    {
      data: memoizedData,
      columns: memoizedColumns,
      defaultColumn,
      filterTypes,
    },
    useSortBy,
    useFilters,
    useGlobalFilter,
  );

  return (
    <Box p={2}>
      {/* <Box> */}
      {headerGroups.map((headerGroup, i) => (
        <HStack key={i} space={2}>
          {headerGroup.headers.map((column) => {
            console.log('COL', column.width);
            const columnSortProps = column.getSortByToggleProps();
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
                  <Icon
                    size="3"
                    viewBox="10 10 20 20"
                    {...columnSortProps}
                    name={
                      column.isSorted
                        ? column.isSortedDesc
                          ? 'more'
                          : 'next'
                        : 'plus'
                    }
                  />
                </HStack>
              </Box>
            );
          })}
        </HStack>
      ))}
      {/* </Box> */}
      <VStack space={2}>
        {rows.map((row) => {
          prepareRow(row);

          return (
            <Pressable key={row.id} onPress={() => handleItemClick(row.id)}>
              <HStack space={2}>
                {row.cells.map((cell, i) => {
                  console.log('CELL', cell);
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
