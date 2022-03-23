import React, { useContext } from 'react';
import { Recipient } from 'core-api-client/src/types/client/Recipient';
import { FormDefinition } from 'core-api-client';
import { useGlobalFilter, useSortBy, useTable } from 'react-table';

import { TableComponent } from './Table.component';
import { TableContext } from './useTableContext';
import { fuzzyTextFilterFn } from './fuzzyFilter';
import { createTableColumns } from './createTableColumns';
import { mapRecordsToTableData } from './mapRecordsToTableData';
import { TableProps } from './types';

type Props = {
  form: FormDefinition;
  records: Recipient[];
  handleItemClick: (id: string) => void;
};

export const TableContainer: React.FC<Props> = ({
  form,
  records,
  handleItemClick,
}) => {
  const tableContext = useContext(TableContext);

  const memoizedData = React.useMemo(
    () => mapRecordsToTableData(records, form),
    [records],
  );
  const memoizedColumns = React.useMemo(() => createTableColumns(form), [form]);

  const filterTypes = React.useMemo(
    () => ({
      fuzzyText: fuzzyTextFilterFn,
    }),
    [],
  );

  const table: TableProps = useTable(
    {
      data: memoizedData,
      columns: memoizedColumns,
      filterTypes,
    },
    useGlobalFilter,
    useSortBy,
  );

  React.useEffect(() => {
    if (!tableContext) return;

    tableContext.setTableInstance(table);
  }, [tableContext]);

  return <TableComponent handleItemClick={handleItemClick} />;
};
