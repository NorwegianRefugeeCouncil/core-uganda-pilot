import React from 'react';
import { useAsyncDebounce } from 'react-table';
import { FormControl, Input } from 'native-base';

import { TableProps } from './types';

type Props = {
  table: TableProps;
};

export const GlobalFilter: React.FC<Props> = ({ table }) => {
  const {
    preGlobalFilteredRows,
    state: { globalFilter },
    setGlobalFilter,
  } = table;

  const [value, setValue] = React.useState(globalFilter);
  const onChange = useAsyncDebounce((v) => {
    setGlobalFilter(v || undefined);
  }, 200);
  const count = preGlobalFilteredRows.length;

  return (
    <FormControl>
      <FormControl.Label>Search</FormControl.Label>
      <Input
        placeholder="Search"
        value={value || ''}
        onChange={(e) => {
          console.log('FILTER', e.target.value);
          onChange(e.target.value);
          setValue(e.target.value);
        }}
      />
      <FormControl.HelperText>{`${count} results`}</FormControl.HelperText>
    </FormControl>
  );
};
