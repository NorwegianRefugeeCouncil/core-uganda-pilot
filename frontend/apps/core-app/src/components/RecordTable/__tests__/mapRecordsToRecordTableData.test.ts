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

  it('should return empty array in case of no records', () => {
    expect(mapRecordsToRecordTableData({ records: [], form })).toEqual([]);
  });

  it('should map records correctly', () => {
    expect(mapRecordsToRecordTableData({ records: [baseRecord1], form })).toEqual([
      {
        field1: 'text',
        field2: 'multi line text',
      },
    ]);
  });

  it('should skip record values whose fields do not exist in the form', () => {
    expect(
      mapRecordsToRecordTableData({
        records: [
          {
            ...baseRecord1,
            values: [
              ...baseRecord1.values,
              { fieldId: 'non-field', value: 'non-value' },
            ],
          },
        ],
        form,
      }),
    ).toEqual([
      {
        field1: 'text',
        field2: 'multi line text',
      },
    ]);
  });
});
