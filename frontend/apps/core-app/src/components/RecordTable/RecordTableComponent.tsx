import React, { useContext } from 'react';
import { Box, HStack, Text, VStack } from 'native-base';
import { Row } from 'react-table';

import { RecordTableContext } from './RecordTableContext';
import { RecordTableRow } from './RecordTableRow';
import { RecordTableHeader } from './RecordTableHeader';

type Props = {
  handleItemClick: (id: string) => void;
};

export const RecordTableComponent: React.FC<Props> = ({ handleItemClick }) => {
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
          <RecordTableHeader column={column} key={column.id} />
        ))}
      </HStack>
      <VStack>
        {rows.map((row: Row) => {
          prepareRow(row);
          return (
            <RecordTableRow
              key={row.id}
              row={row}
              handleRowClick={handleItemClick}
            />
          );
        })}
      </VStack>
    </Box>
  );
};
