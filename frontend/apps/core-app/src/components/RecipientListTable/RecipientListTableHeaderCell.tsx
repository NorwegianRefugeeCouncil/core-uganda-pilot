import { HStack, IconButton, Text } from 'native-base';
import { Icon } from 'core-design-system';
import React from 'react';
import { ColumnInstance, UseSortByColumnProps } from 'react-table';

import { RecipientListTableEntry } from './types';

type Props<T extends Record<string, any>> = {
  column: ColumnInstance & UseSortByColumnProps<T>;
};

export const RecipientListTableHeaderCell: React.FC<
  Props<RecipientListTableEntry>
> = ({ column }) => {
  const sortIcon = column.isSorted
    ? column.isSortedDesc
      ? 'more'
      : 'arrowUp'
    : 'arrowDown';

  const handleSortToggle = () => column.toggleSortBy(!column.isSortedDesc);

  return (
    <HStack
      width={column.width}
      p="2"
      alignItems="center"
      justifyContent="space-between"
      flexGrow={1}
    >
      {column.render(Text, {
        fontWeight: 'bold',
        fontSize: 'xs',
        lineHeight: '4xs',
        children: column.Header,
      })}
      <IconButton
        onPress={handleSortToggle}
        colorScheme="secondary"
        variant="ghost"
        size="sm"
        icon={<Icon size={5} name={sortIcon} />}
      />
    </HStack>
  );
};
