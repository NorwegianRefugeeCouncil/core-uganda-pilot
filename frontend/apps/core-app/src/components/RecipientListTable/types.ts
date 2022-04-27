import {
  ColumnInstance,
  TableInstance,
  TableState,
  UseGlobalFiltersInstanceProps,
  UseGlobalFiltersState,
  UseSortByColumnProps,
  UseSortByInstanceProps,
} from 'react-table';

export interface SortedFilteredTable<T extends Record<string, any>>
  extends Omit<
      TableInstance,
      'flatRows' | 'rows' | 'rowsById' | 'columns' | 'state'
    >,
    UseGlobalFiltersInstanceProps<T>,
    UseSortByInstanceProps<T>,
    UseGlobalFiltersState<T> {
  columns: Array<ColumnInstance & UseSortByColumnProps<T>>;
  state: TableState & UseGlobalFiltersState<T>;
}

export type RecipientListTableEntry = Record<string, any>;
