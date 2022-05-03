import React from 'react';
import { FilterValue, useAsyncDebounce } from 'react-table';

import { RecipientListTableFilterComponent } from './RecipientListTableFilter.component';

type Props = {
  globalFilter: FilterValue;
  setGlobalFilter: (filterValue: FilterValue) => void;
};

export const RecipientListTableFilterContainer: React.FC<Props> = ({
  globalFilter,
  setGlobalFilter,
}) => {
  const [value, setValue] = React.useState(globalFilter);

  const handleChange = useAsyncDebounce((v) => {
    setGlobalFilter(v || undefined);
    setValue(v);
  }, 20);

  return (
    <RecipientListTableFilterComponent onChange={handleChange} value={value} />
  );
};
