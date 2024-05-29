import { FormType } from 'core-api-client';

import { mapRecordsToRecipientTableData } from '../mapRecordsToRecipientTableData';
import { makeForm, makeRecord, makeField } from '../../../testUtils/mockData';

describe('mapRecordsToTableData', () => {
  const field1 = makeField(1, false, false, { text: {} });
  const field2 = makeField(2, true, false, { multilineText: {} });
  const field3 = makeField(3, false, false, { multilineText: {} });
  const field4 = makeField(4, false, false, { multilineText: {} });
  const form1 = makeForm(1, FormType.RecipientFormType, [field1, field2]);
  const form2 = makeForm(2, FormType.RecipientFormType, [field3, field4]);
  const record1 = makeRecord(1, form1);
  const record2 = makeRecord(2, form2);
  const record3 = makeRecord(3, form1);
  const record4 = makeRecord(4, form2);

  it('should map records correctly', () => {
    const data = [
      [
        { form: form1, record: record1 },
        { form: form2, record: record2 },
      ],
    ];
    expect(mapRecordsToRecipientTableData(data)).toEqual([
      {
        'field-1': 'value-field-1',
        'field-2': 'value-field-2',
        'field-3': 'value-field-3',
        'field-4': 'value-field-4',
        recordId: 'record-2',
      },
    ]);
  });

  it('should map multiple records correctly', () => {
    const data = [
      [
        { form: form1, record: record1 },
        { form: form2, record: record2 },
      ],
      [
        { form: form1, record: record3 },
        { form: form2, record: record4 },
      ],
    ];
    expect(mapRecordsToRecipientTableData(data)).toEqual([
      {
        'field-1': 'value-field-1',
        'field-2': 'value-field-2',
        'field-3': 'value-field-3',
        'field-4': 'value-field-4',
        recordId: 'record-2',
      },
      {
        'field-1': 'value-field-1',
        'field-2': 'value-field-2',
        'field-3': 'value-field-3',
        'field-4': 'value-field-4',
        recordId: 'record-4',
      },
    ]);
  });

  it('should skip record values whose fields do not exist in the form', () => {
    expect(
      mapRecordsToRecipientTableData([
        [
          {
            record: {
              ...record1,
              values: [
                ...record1.values,
                { fieldId: 'non-field', value: 'non-value' },
              ],
            },
            form: form1,
          },
        ],
      ]),
    ).toEqual([
      {
        'field-1': 'value-field-1',
        'field-2': 'value-field-2',
        recordId: 'record-1',
      },
    ]);
  });
});
