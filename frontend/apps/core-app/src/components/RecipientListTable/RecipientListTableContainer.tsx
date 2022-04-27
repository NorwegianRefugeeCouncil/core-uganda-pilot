import React from 'react';
import { FormWithRecord, Record } from 'core-api-client';
import { useGlobalFilter, useSortBy, useTable } from 'react-table';

import { RecipientListTableComponent } from './RecipientListTableComponent';
import { RecipientListTableContext } from './RecipientListTableContext';
import { createTableColumns } from './createTableColumns';
import { mapRecordsToRecordTableData } from './mapRecordsToRecordTableData';
import { RecipientListTableEntry, SortedFilteredTable } from './types';

type Props<T extends Record> = {
  data: FormWithRecord<T>[][];
  onItemClick: (id: string) => void;
};

export const RecipientListTableContainer: React.FC<Props<Record>> = ({
  data,
  onItemClick,
}) => {
  const context = React.useContext(RecipientListTableContext);

  const memoizedData = React.useMemo(
    () => mapRecordsToRecordTableData(data),
    [JSON.stringify(data)],
  );
  const memoizedColumns = React.useMemo(
    () => createTableColumns(data),
    [JSON.stringify(data)],
  );

  const table: SortedFilteredTable<RecipientListTableEntry> = useTable(
    {
      data: memoizedData,
      columns: memoizedColumns,
    },
    useGlobalFilter,
    useSortBy,
  );

  React.useEffect(() => {
    if (!context) return;

    context.setTableInstance(table);
  }, [context]);

  return <RecipientListTableComponent onItemClick={onItemClick} />;
};
