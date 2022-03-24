import { createContext, Dispatch, SetStateAction } from 'react';

import { RecordTableEntry, SortedFilteredTable } from './types';

export type TableContext<T> = {
  tableInstance: SortedFilteredTable<T> | null;
  setTableInstance: Dispatch<SetStateAction<SortedFilteredTable<T> | null>>;
};

export const RecordTableContext =
  createContext<TableContext<RecordTableEntry> | null>(null);
