import { FormDefinition, FormWithRecord, Record } from 'core-api-client';
import { Column } from 'react-table';

import { RecipientListTableEntry } from './types';

export const createTableColumns = (
  form: FormDefinition,
): Column<RecipientListTableEntry>[] => {
  console.log('createTableColumns', form);

  return form.fields
    .filter((f) => !f.key)
    .map(({ name, id }) => ({
      Header: name,
      accessor: id,
    }));
};
