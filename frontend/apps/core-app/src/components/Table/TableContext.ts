import { createContext, Dispatch, SetStateAction } from 'react';

import { TableProps } from './types';

export type TableContextType = {
  tableInstance: TableProps | null;
  setTableInstance: Dispatch<SetStateAction<TableProps | null>>;
};

export const TableContext = createContext<TableContextType | null>(null);
