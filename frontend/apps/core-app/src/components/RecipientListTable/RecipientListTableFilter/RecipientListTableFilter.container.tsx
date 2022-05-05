import React from 'react';
import { FilterValue, useAsyncDebounce } from 'react-table';

import { RecipientListTableFilterComponent } from './RecipientListTableFilter.component';

type Props = {
  filter: FilterValue;
  setFilter: (filterValue: FilterValue) => void;
};

export const RecipientListTableFilterContainer: React.FC<Props> = ({
  filter,
  setFilter,
}) => {
  const [value, setValue] = React.useState(filter);

  const handleChange = useAsyncDebounce((v) => {
    setFilter(v || undefined);
    setValue(v);
  }, 20);

  return (
    <RecipientListTableFilterComponent onChange={handleChange} value={value} />
  );
};
