import { mapRecordsToRecordTableData } from '../mapRecordsToRecordTableData';
import {
  baseForm,
  baseMultilineTextField,
  baseRecord1,
  baseTextField,
} from '../../../testUtils/baseObjects';

describe('mapRecordsToTableData', () => {
  const form = baseForm;
  form.fields = [baseTextField, baseMultilineTextField];

  it('should map records correctly', () => {
    expect(
      mapRecordsToRecordTableData([[{ record: baseRecord1, form }]]),
    ).toEqual([
      {
        field1: 'text',
      },
    ]);
  });

  it('should skip record values whose fields do not exist in the form', () => {
    expect(
      mapRecordsToRecordTableData([
        [
          {
            record: {
              ...baseRecord1,
              values: [
                ...baseRecord1.values,
                { fieldId: 'non-field', value: 'non-value' },
              ],
            },
            form,
          },
        ],
      ]),
    ).toEqual([
      {
        field1: 'text',
      },
    ]);
  });
});
