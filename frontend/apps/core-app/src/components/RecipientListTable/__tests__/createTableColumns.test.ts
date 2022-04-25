import { FormType } from 'core-api-client';

import { createTableColumns } from '../createTableColumns';
import { makeField, makeForm, makeRecord } from '../../../testUtils/mockData';

describe('createTableColumns', () => {
  const fields = [
    makeField(1, false, false, { text: {} }),
    makeField(1, true, false, { multilineText: {} }),
  ];
  const form = makeForm(1, FormType.DefaultFormType, fields);

  const record = makeRecord(1, form);

  it('should create the correct columns', () => {
    expect(
      createTableColumns([
        [
          {
            record: {
              ...record,
              values: [
                ...record.values,
                { fieldId: 'non-field', value: 'non-value' },
              ],
            },
            form,
          },
        ],
      ]),
    ).toEqual([{ Header: 'field 1', accessor: 'field1' }]);
  });
});
