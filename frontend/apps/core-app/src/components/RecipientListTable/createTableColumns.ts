import { FormWithRecord } from 'core-api-client';
import { Column } from 'react-table';
import { Recipient } from 'core-api-client/src/types/client/Recipient';

import { RecipientListTableEntry } from './types';

export const createTableColumns = (
  recipientWithForm: FormWithRecord<Recipient>[] | null,
): Column<RecipientListTableEntry>[] => {
  if (!recipientWithForm) {
    return [];
  }

  const initColumn: Column<RecipientListTableEntry> = {
    Header: 'recordId',
    accessor: 'recordId',
    hidden: true,
  };

  return recipientWithForm.reduce(
    (allColumns: Column<RecipientListTableEntry>[], formWithRecord) => {
      const columns = formWithRecord.form.fields.map(({ name, id, key }) => ({
        Header: name,
        accessor: id,
        hidden: key,
      }));
      return allColumns.concat(columns);
    },
    [initColumn],
  );
};
