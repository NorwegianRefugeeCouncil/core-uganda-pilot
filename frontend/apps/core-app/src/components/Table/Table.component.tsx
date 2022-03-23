import React, { useContext } from 'react';
import { Box, HStack, Text, VStack } from 'native-base';
import { Row } from 'react-table';

import { TableContext } from './TableContext';
import { TableRow } from './TableRow';
import { TableHeader } from './TableHeader';

type Props = {
  handleItemClick: (id: string) => void;
};

export const TableComponent: React.FC<Props> = ({ handleItemClick }) => {
  const tableContext = useContext(TableContext);
  if (!tableContext) return null;

  const { tableInstance } = tableContext;
  if (!tableInstance) return null;

  const { rows, columns, prepareRow, globalFilteredRows } = tableInstance;

  return (
    <Box>
      <Text level="2">{globalFilteredRows.length} beneficiaries</Text>
      <HStack>
        {columns.map((column) => (
          <TableHeader column={column} key={column.id} />
        ))}
      </HStack>
      <VStack>
        {rows.map((row: Row) => {
          prepareRow(row);
          return (
            <TableRow
              key={row.id}
              row={row}
              handleRowClick={() => handleItemClick(row.id)}
            />
          );
        })}
      </VStack>
    </Box>
  );
};
