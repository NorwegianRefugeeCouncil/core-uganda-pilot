import { FormDefinition } from 'core-api-client';
import { Column } from 'react-table';

import { RecipientListTableEntry } from './types';

export const createTableColumns = (
  form: FormDefinition,
): Column<RecipientListTableEntry>[] =>
  form.fields
    .filter((f) => !f.key)
    .map(({ name, id }) => ({
      Header: name,
      accessor: id,
    }));
