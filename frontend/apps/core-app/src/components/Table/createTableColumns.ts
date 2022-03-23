import { FormDefinition } from 'core-api-client';
import { Column } from 'react-table';

export const createTableColumns = (form: FormDefinition): Column[] =>
  form.fields.map((field) => {
    return {
      Header: field.name,
      accessor: field.id,
    };
  });
