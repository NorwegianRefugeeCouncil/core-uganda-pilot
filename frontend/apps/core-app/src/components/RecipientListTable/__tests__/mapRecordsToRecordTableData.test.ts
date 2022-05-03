import { FormType } from 'core-api-client';

import { mapRecordsToRecordTableData } from '../mapRecordsToRecordTableData';
import { makeForm, makeRecord, makeField } from '../../../testUtils/mockData';

describe('mapRecordsToTableData', () => {
  const fields = [
    makeField(1, false, false, { text: {} }),
    makeField(1, true, false, { multilineText: {} }),
  ];
  const form = makeForm(1, FormType.DefaultFormType, fields);

  const record1 = makeRecord(1, form);
  const record2 = makeRecord(2, form);

  it('should map records correctly', () => {
    expect(mapRecordsToRecordTableData([[{ record: record1, form }]])).toEqual([
      {
        'field-1': 'value-field-1',
      },
    ]);
  });

  it('should map multiple records correctly', () => {
    expect(
      mapRecordsToRecordTableData([
        [
          { record: record1, form },
          { record: record2, form },
        ],
      ]),
    ).toEqual([
      {
        'field-1': 'value-field-1',
      },
    ]);
  });

  it('should skip record values whose fields do not exist in the form', () => {
    expect(
      mapRecordsToRecordTableData([
        [
          {
            record: {
              ...record1,
              values: [
                ...record1.values,
                { fieldId: 'non-field', value: 'non-value' },
              ],
            },
            form,
          },
        ],
      ]),
    ).toEqual([
      {
        'field-1': 'value-field-1',
      },
    ]);
  });
});
