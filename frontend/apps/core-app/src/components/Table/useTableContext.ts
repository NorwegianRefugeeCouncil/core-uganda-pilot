import { createContext, Dispatch, SetStateAction } from 'react';
import { TableInstance } from 'react-table';

import { TableProps } from './types';

export type TableContextType = {
  tableInstance: TableProps | null;
  setTableInstance: Dispatch<SetStateAction<TableInstance | null>>;
};

export const TableContext = createContext<TableContextType | null>(null);
