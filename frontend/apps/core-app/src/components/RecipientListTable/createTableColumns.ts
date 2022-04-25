import { FormWithRecord, Record } from 'core-api-client';
import { Column } from 'react-table';

import { RecipientListTableEntry } from './types';

export const createTableColumns = (
  data: FormWithRecord<Record>[][],
): Column<RecipientListTableEntry>[] => {
  return data[0].reduce(
    (allColumns: Column<RecipientListTableEntry>[], formWithRecord) => {
      const columnsPerForm = formWithRecord.form.fields
        .filter((f) => !f.key)
        .map(({ name, id }) => ({
          Header: name,
          accessor: id,
        }));

      return [...allColumns, ...columnsPerForm];
    },
    [],
  );
};
