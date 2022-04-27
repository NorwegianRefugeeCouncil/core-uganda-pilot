import { createContext, Dispatch, SetStateAction } from 'react';

import { RecipientListTableEntry, SortedFilteredTable } from './types';

export type TableContext<T> = {
  tableInstance: SortedFilteredTable<T> | null;
  setTableInstance: Dispatch<SetStateAction<SortedFilteredTable<T> | null>>;
};

export const RecipientListTableContext =
  createContext<TableContext<RecipientListTableEntry> | null>(null);
