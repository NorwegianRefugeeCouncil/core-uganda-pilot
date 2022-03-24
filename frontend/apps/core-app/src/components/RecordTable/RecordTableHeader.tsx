import { HStack, Pressable, Text } from 'native-base';
import { Icon } from 'core-design-system';
import React from 'react';
import { ColumnInstance, UseSortByColumnProps } from 'react-table';

import { RecordTableEntry } from './types';

type Props<T extends Record<string, any>> = {
  column: ColumnInstance & UseSortByColumnProps<T>;
};

export const RecordTableHeader: React.FC<Props<RecordTableEntry>> = ({
  column,
}) => {
  const { onClick } = column.getSortByToggleProps();
  const sortIcon = column.isSorted
    ? column.isSortedDesc
      ? 'more'
      : 'next'
    : 'plus';

  return (
    <HStack
      width={column.width}
      p="2"
      borderBottomColor="neutral.300"
      borderBottomWidth="1"
      alignItems="center"
      justifyContent="space-between"
      bg="neutral.200"
    >
      {column.render(Text, {
        variant: 'body',
        fontWeight: '700',
        children: column.Header,
      })}
      {onClick && (
        <Pressable onPress={(e) => onClick(e)} p="1">
          <Icon size="3" viewBox="10 10 20 20" name={sortIcon} />
        </Pressable>
      )}
    </HStack>
  );
};
