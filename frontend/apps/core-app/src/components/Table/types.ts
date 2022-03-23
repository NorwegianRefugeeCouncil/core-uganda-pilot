import {
  TableInstance,
  UseGlobalFiltersInstanceProps, UseSortByColumnProps,
  UseSortByInstanceProps
} from "react-table";

export type TableProps = TableInstance extends
    UseGlobalFiltersInstanceProps<any> &
    UseSortByInstanceProps<any> &
    UseSortByColumnProps<any>
