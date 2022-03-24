import { createTableColumns } from '../createTableColumns';
import {
  baseForm,
  baseMultilineTextField,
  baseTextField,
} from '../../../testUtils/baseObjects';

describe('createTableColumns', () => {
  const form = baseForm;
  form.fields = [baseTextField, baseMultilineTextField];

  it('should create the correct columns', () => {
    expect(createTableColumns(form)).toEqual([
      { Header: 'field 1', accessor: 'field1' },
      { Header: 'field 2', accessor: 'field2' },
    ]);
  });
});
