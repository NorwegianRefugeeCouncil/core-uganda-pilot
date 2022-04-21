import React from 'react';
import { useAsyncDebounce } from 'react-table';
import { FormControl, Input } from 'native-base';
import { NativeSyntheticEvent, TextInputChangeEventData } from 'react-native';

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
  const onChange = useAsyncDebounce(
    (v: NativeSyntheticEvent<TextInputChangeEventData>) => {
      setGlobalFilter(v.nativeEvent.text || undefined);
      setValue(v.nativeEvent.text);
    },
  );

  return (
    <FormControl>
      <FormControl.Label>Beneficiary Name</FormControl.Label>
      <Input
        type="text"
        placeholder="Search"
        value={value || ''}
        onChange={onChange}
      />
    </FormControl>
  );
};
