import React from 'react';
import {
  Box,
  HStack,
  VStack,
  FlatList,
  ScrollView,
  Text,
  Badge,
  Skeleton,
} from 'native-base';
import { Row } from 'react-table';

import { RecipientListTableRow } from './RecipientListTableRow';
import { RecipientListTableHeaderCell } from './RecipientListTableHeaderCell';
import { RecipientListTableEntry, SortedFilteredTable } from './types';

type Props = {
  onItemClick: (id: string) => void;
  title: string;
  table: SortedFilteredTable<RecipientListTableEntry>;
  loading: boolean;
  error: string | null;
};

export const RecipientListTableComponent: React.FC<Props> = ({
  onItemClick,
  title,
  table,
  error,
  loading,
}) => {
  const { rows, columns, prepareRow } = table;

  const renderRow = ({ item }: { item: Row }) => {
    prepareRow(item);
    return (
      <RecipientListTableRow
        key={item.id}
        row={item}
        handleRowClick={() => onItemClick(item.values.recordId)}
      />
    );
  };

  return (
    <Box mb="63px">
      <Box>
        <Box
          flexDirection="row"
          alignItems="center"
          justifyContent="flex-start"
        >
          <Text variant="heading" level="2">
            {title}
          </Text>
          <Badge
            bg="secondary.500"
            height="5"
            width="7"
            borderRadius="4px"
            mx="12px"
          >
            <Text variant="heading" level="5" color="white">
              {rows.length}
            </Text>
          </Badge>
        </Box>
      </Box>
      {!loading && !error && (
        <ScrollView
          horizontal
          contentContainerStyle={{
            flexGrow: 1,
          }}
        >
          <Box width="100%">
            <HStack
              bg="secondary.100"
              borderBottomColor="neutral.400"
              borderBottomWidth="1"
            >
              {columns
                .filter((c) => !c.hidden)
                .map((column) => (
                  <RecipientListTableHeaderCell
                    column={column}
                    key={column.id}
                  />
                ))}
            </HStack>
            <VStack>
              <FlatList data={rows} renderItem={renderRow} />
            </VStack>
          </Box>
        </ScrollView>
      )}
      {loading && <Skeleton h="40" p="4" />}
      {error && <Text>{error}</Text>}
    </Box>
  );
};
