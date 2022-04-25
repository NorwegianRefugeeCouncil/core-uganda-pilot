import { HStack, IconButton, Pressable, Text } from 'native-base';
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
  const { onClick: handleOnClick } = column.getSortByToggleProps();
  const sortIcon = column.isSorted
    ? column.isSortedDesc
      ? 'more'
      : 'arrowUp'
    : 'arrowDown';

  return (
    <HStack
      width={column.width}
      p="2"
      borderBottomColor="neutral.400"
      borderBottomWidth="1"
      alignItems="center"
      justifyContent="space-between"
      bg="neutral.100"
      flexGrow={1}
    >
      {column.render(Text, {
        // variant: 'body',
        fontWeight: '700',
        fontSize: '18px',
        lineHeight: '17px',
        children: column.Header,
      })}
      {handleOnClick && (
        <IconButton
          onPress={(e) => handleOnClick(e)}
          colorScheme="secondary"
          variant="ghost"
          size="sm"
          icon={<Icon size={5} name={sortIcon} />}
        />
      )}
    </HStack>
  );
};
