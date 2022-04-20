import * as React from 'react';
import { useTable, Column } from 'react-table';

import { SubFormTableComponent } from './SubFormTable.component';

type D = Record<string, string>;

type Props = {
  data: D[];
  columns: Column<D>[];
  onDelete?: (idx: number) => void;
};

export const SubFormTableContainer: React.FC<Props> = ({
  data,
  columns,
  onDelete,
}) => {
  const table = useTable({ data, columns });

  return <SubFormTableComponent table={table} onDelete={onDelete} />;
};
