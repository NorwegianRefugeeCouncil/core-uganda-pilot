import { mapRecordsToTableData } from '../mapRecordsToTableData';
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
    expect(mapRecordsToTableData([], form)).toEqual([]);
  });

  it('should map records correctly', () => {
    expect(mapRecordsToTableData([baseRecord1], form)).toEqual([
      {
        field1: 'text',
        field2: 'multi line text',
      },
    ]);
  });

  it('should skip record values whose fields do not exist in the form', () => {
    expect(
      mapRecordsToTableData(
        [
          {
            ...baseRecord1,
            values: [
              ...baseRecord1.values,
              { fieldId: 'non-field', value: 'non-value' },
            ],
          },
        ],
        form,
      ),
    ).toEqual([
      {
        field1: 'text',
        field2: 'multi line text',
      },
    ]);
  });
});
