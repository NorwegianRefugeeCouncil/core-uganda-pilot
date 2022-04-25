import React from 'react';
import { FilterValue, useAsyncDebounce } from 'react-table';

import { RecipientListTableEntry, SortedFilteredTable } from '../types';

import { RecipientListTableFilterComponent } from './RecipientListTableFilter.component';

type Props<T> = {
  table: SortedFilteredTable<T>;
  globalFilter?: FilterValue;
  setGlobalFilter?: (filterValue: FilterValue) => void;
};

export const RecipientListTableFilterContainer: React.FC<
  Props<RecipientListTableEntry>
> = ({ table }) => {
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
    <RecipientListTableFilterComponent onChange={handleChange} value={value} />
  );
};
