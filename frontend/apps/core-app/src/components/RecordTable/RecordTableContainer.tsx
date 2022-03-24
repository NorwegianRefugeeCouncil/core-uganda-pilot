import React, { useContext } from 'react';
import { Record } from 'core-api-client';
import {
  TableInstance,
  useGlobalFilter,
  useSortBy,
  useTable,
} from 'react-table';
import { FormWithRecords } from 'core-api-client/src';

import { RecordTableComponent } from './RecordTableComponent';
import { RecordTableContext } from './RecordTableContext';
import { createTableColumns } from './createTableColumns';
import { mapRecordsToRecordTableData } from './mapRecordsToRecordTableData';
import { RecordTableEntry, SortedFilteredTable } from './types';

type Props<T extends Record> = {
  data: FormWithRecords<T>;
  handleItemClick: (id: string) => void;
};

export const RecordTableContainer: React.FC<Props<Record>> = ({
  data,
  handleItemClick,
}) => {
  const context = useContext(RecordTableContext);

  const memoizedData = React.useMemo(
    () => mapRecordsToRecordTableData(data),
    [JSON.stringify(data.records)],
  );
  const memoizedColumns = React.useMemo(
    () => createTableColumns(data.form),
    [JSON.stringify(data.form)],
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

  return <RecordTableComponent handleItemClick={handleItemClick} />;
};
