import React from 'react';
import { useAsyncDebounce } from 'react-table';
import { FormControl } from 'native-base';

import { Input } from '../Web/Input';

import { SortedFilteredTable } from './types';

type Props = {
  table: SortedFilteredTable;
};

export const GlobalTableFilter: React.FC<Props> = ({ table }) => {
  const {
    state: { globalFilter },
    setGlobalFilter,
  } = table;

  const [value, setValue] = React.useState(globalFilter);
  const onChange = useAsyncDebounce((v) => {
    setGlobalFilter(v || undefined);
  }, 200);

  return (
    <FormControl>
      <FormControl.Label>Beneficiary Name</FormControl.Label>
      <Input
        type="text"
        placeholder="Search"
        value={value || ''}
        onChange={(e) => {
          onChange(e);
          setValue(e);
        }}
      />
    </FormControl>
  );
};
