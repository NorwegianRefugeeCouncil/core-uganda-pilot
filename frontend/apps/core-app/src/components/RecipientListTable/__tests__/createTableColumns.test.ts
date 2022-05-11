import { FormType } from 'core-api-client';

import { createTableColumns } from '../createTableColumns';
import { makeField, makeForm, makeRecord } from '../../../testUtils/mockData';

describe('createTableColumns', () => {
  const field1 = makeField(1, false, false, { text: {} });
  const field2 = makeField(2, true, false, { multilineText: {} });
  const field3 = makeField(3, false, false, { multilineText: {} });
  const field4 = makeField(4, false, false, { multilineText: {} });
  const form1 = makeForm(1, FormType.DefaultFormType, [field1, field2]);
  const form2 = makeForm(2, FormType.DefaultFormType, [field3, field4]);
  const record1 = makeRecord(1, form1);
  const record2 = makeRecord(2, form2);

  const data = [
    { form: form1, record: record1 },
    { form: form2, record: record2 },
  ];

  it('should return an empty array if data is null', () => {
    expect(createTableColumns(null)).toEqual([]);
  });

  it('should create the correct columns', () => {
    expect(createTableColumns(data)).toEqual([
      { Header: 'recordId', accessor: 'recordId', hidden: true },
      { Header: field1.name, accessor: field1.id, hidden: false },
      { Header: field2.name, accessor: field2.id, hidden: true },
      { Header: field3.name, accessor: field3.id, hidden: false },
      { Header: field4.name, accessor: field4.id, hidden: false },
    ]);
  });
});
