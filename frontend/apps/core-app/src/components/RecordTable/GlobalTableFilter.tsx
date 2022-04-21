import React from 'react';
import { useAsyncDebounce } from 'react-table';
import { FormControl, Input } from 'native-base';

// import { Input } from '../Web/Input';

import { RecordTableEntry, SortedFilteredTable } from './types';

type Props<T> = {
  table: SortedFilteredTable<T>;
};

export const GlobalTableFilter: React.FC<Props<RecordTableEntry>> = ({
  table,
}) => {
  const {
    state: { globalFilter },
    setGlobalFilter,
  } = table;

  const [value, setValue] = React.useState(globalFilter);
  const onChange = useAsyncDebounce((v) => {
    setGlobalFilter(v || undefined);
    setValue(v);
  }, 200);

  return (
    <FormControl>
      <FormControl.Label>Beneficiary Name</FormControl.Label>
      <Input
        testID="jkdsfbf"
        type="text"
        // role="search"
        placeholder="Search"
        value={value || ''}
        onChange={onChange}
      />
    </FormControl>
  );
};
