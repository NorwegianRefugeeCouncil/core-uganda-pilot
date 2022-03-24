import { createTableColumns } from '../createTableColumns';
import {
  baseForm,
  baseMultilineTextField,
  baseRecord1,
  baseTextField,
} from '../../../testUtils/baseObjects';

describe('createTableColumns', () => {
  const form = baseForm;
  form.fields = [baseTextField, baseMultilineTextField];

  it('should create the correct columns', () => {
    expect(
      createTableColumns([
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
    ).toEqual([{ Header: 'field 1', accessor: 'field1' }]);
  });
});
