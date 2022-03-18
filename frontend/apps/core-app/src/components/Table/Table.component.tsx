import React from 'react';
import { useSortBy, useTable } from 'react-table';
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

// Define a default UI for filtering
function GlobalFilter({
  preGlobalFilteredRows,
  globalFilter,
  setGlobalFilter,
}) {
  const count = preGlobalFilteredRows.length;
  const [value, setValue] = React.useState(globalFilter);
  const onChange = useAsyncDebounce((value) => {
    setGlobalFilter(value || undefined);
  }, 200);

  return (
    <span>
      Search:{' '}
      <input
        value={value || ''}
        onChange={(e) => {
          setValue(e.target.value);
          onChange(e.target.value);
        }}
        placeholder={`${count} records...`}
        style={{
          fontSize: '1.1rem',
          border: '0',
        }}
      />
    </span>
  );
}

// Define a default UI for filtering
function DefaultColumnFilter({
  column: { filterValue, preFilteredRows, setFilter },
}) {
  const count = preFilteredRows.length;

  return (
    <input
      value={filterValue || ''}
      onChange={(e) => {
        setFilter(e.target.value || undefined); // Set undefined to remove the filter entirely
      }}
      placeholder={`Search ${count} records...`}
    />
  );
}

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

  const { rows, headerGroups, prepareRow } = useTable(
    {
      data: memoizedData,
      columns: memoizedColumns,
    },
    useSortBy,
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
