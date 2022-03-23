import { FilterValue, Row } from 'react-table';

export function fuzzyTextFilterFn(
  rows: Row,
  id: string,
  filterValue: FilterValue,
) {
  // TODO use fusejs
  console.log('FUZZY', id, filterValue);
  return rows;
}
// fuzzyTextFilterFn.autoRemove = (val) => !val;
