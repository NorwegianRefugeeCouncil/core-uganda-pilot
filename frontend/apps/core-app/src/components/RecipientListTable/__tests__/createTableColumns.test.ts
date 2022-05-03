import { FormType } from 'core-api-client';

import { createTableColumns } from '../createTableColumns';
import { makeField, makeForm } from '../../../testUtils/mockData';

describe('createTableColumns', () => {
  const fields = [
    makeField(1, false, false, { text: {} }),
    makeField(1, true, false, { multilineText: {} }),
  ];
  const form = makeForm(1, FormType.DefaultFormType, fields);

  it('should create the correct columns', () => {
    expect(createTableColumns(form)).toEqual([
      { Header: 'field-name-1', accessor: 'field-1' },
    ]);
  });
});
