import React, { useContext } from 'react';
import { FormWithRecord, Record } from 'core-api-client';
import { useGlobalFilter, useSortBy, useTable } from 'react-table';

import { RecordTableComponent } from './RecordTableComponent';
import { RecordTableContext } from './RecordTableContext';
import { createTableColumns } from './createTableColumns';
import { mapRecordsToRecordTableData } from './mapRecordsToRecordTableData';
import { RecordTableEntry, SortedFilteredTable } from './types';

type Props<T extends Record> = {
  data: FormWithRecord<T>[][];
  onItemClick: (id: string) => void;
};

export const RecordTableContainer: React.FC<Props<Record>> = ({
  data,
  onItemClick,
}) => {
  const context = useContext(RecordTableContext);

  const memoizedData = React.useMemo(
    () => mapRecordsToRecordTableData(data),
    [JSON.stringify(data)],
  );
  const memoizedColumns = React.useMemo(
    () => createTableColumns(data),
    [JSON.stringify(data)],
  );

  const table: SortedFilteredTable<RecordTableEntry> = useTable(
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

  return <RecordTableComponent onItemClick={onItemClick} />;
};
