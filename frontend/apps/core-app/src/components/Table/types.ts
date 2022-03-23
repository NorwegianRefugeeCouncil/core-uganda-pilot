import {
  ColumnInstance,
  TableInstance,
  TableState,
  UseGlobalFiltersInstanceProps,
  UseGlobalFiltersState,
  UseSortByColumnProps,
  UseSortByInstanceProps,
} from 'react-table';

export interface TableProps
  extends Omit<
      TableInstance,
      'flatRows' | 'rows' | 'rowsById' | 'columns' | 'state'
    >,
    UseGlobalFiltersInstanceProps<any>,
    UseSortByInstanceProps<any>,
    UseGlobalFiltersState<any> {
  columns: Array<ColumnInstance & UseSortByColumnProps<any>>;
  state: TableState & UseGlobalFiltersState<any>;
}
