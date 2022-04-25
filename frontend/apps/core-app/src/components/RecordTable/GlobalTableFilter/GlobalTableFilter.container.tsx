import React from 'react';
import { FilterValue, useAsyncDebounce } from 'react-table';

import { RecordTableEntry, SortedFilteredTable } from '../types';

import { GlobalTableFilterComponent } from './GlobalTableFilter.component';

type Props<T> = {
  table: SortedFilteredTable<T>;
  globalFilter?: FilterValue;
  setGlobalFilter?: (filterValue: FilterValue) => void;
};

export const GlobalTableFilterContainer: React.FC<Props<RecordTableEntry>> = ({
  table,
}) => {
  const {
    state: { globalFilter },
    setGlobalFilter,
  } = table;

  const [value, setValue] = React.useState(globalFilter);

  const handleChange = useAsyncDebounce((v) => {
    setGlobalFilter(v || undefined);
    setValue(v);
  }, 20);

  return (
    <GlobalTableFilterComponent handleChange={handleChange} value={value} />
  );
};
