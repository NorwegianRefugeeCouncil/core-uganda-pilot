import { FormType, FormWithRecord } from 'core-api-client/src';
import { Recipient } from 'core-api-client/src/types/client/Recipient';
import { FieldDefinition } from 'core-api-client';

import { prettifyData } from '../prettifyData';

const makeField = (key: boolean, formIndex: number) => ({
  id: `${key ? 'keyF' : 'f'}ieldId${formIndex}`,
  name: `${key ? 'keyF' : 'f'}ieldName${formIndex}`,
  fieldType: {
    reference: { formId: 'formId', databaseId: 'databaseId' },
  },
  key,
  code: '',
  description: 'description',
  required: true,
});

const makeForm = (index: number, fields: FieldDefinition[]) => ({
  fields,
  id: `formId${index}`,
  formType: FormType.RecipientFormType,
  folderId: 'folderId',
  databaseId: 'databaseId',
  name: `formName${index}`,
  code: '',
});

const makeRecord = (index: number) => ({
  id: `recordId${index}`,
  formId: `formId${index}`,
  ownerId: undefined,
  values: [
    { value: 'keyValue', fieldId: `keyFieldId${index}` },
    { value: 'value', fieldId: `fieldId${index}` },
  ],
  databaseId: 'databaseId',
});

const keyField1 = makeField(true, 1);
const nonKeyField1 = makeField(false, 1);
const keyField2 = makeField(true, 2);
const nonKeyField2 = makeField(false, 2);
const keyField3 = makeField(true, 3);
const nonKeyField3 = makeField(false, 3);

const form1 = makeForm(1, [keyField1, nonKeyField1]);
const form2 = makeForm(2, [keyField2, nonKeyField2]);
const form3 = makeForm(3, [keyField3, nonKeyField3]);

const record1 = makeRecord(1);
const record2 = makeRecord(2);
const record3 = makeRecord(3);

describe('prettifyData', () => {
  it('should remove the key fields and merge the first two forms', () => {
    const originalData: FormWithRecord<Recipient>[] = [
      { form: form1, record: record1 },
      { form: form2, record: record2 },
      { form: form3, record: record3 },
    ];
    const result = prettifyData(originalData);
    expect(result).toEqual([
      {
        form: {
          fields: [nonKeyField1, nonKeyField2],
          id: 'formId2',
          formType: FormType.RecipientFormType,
          folderId: 'folderId',
          databaseId: 'databaseId',
          name: 'formName2',
          code: '',
        },
        record: {
          id: 'recordId2',
          formId: 'formId2',
          ownerId: undefined,
          values: [
            { value: 'keyValue', fieldId: 'keyFieldId1' },
            { value: 'value', fieldId: 'fieldId1' },
            { value: 'keyValue', fieldId: 'keyFieldId2' },
            { value: 'value', fieldId: 'fieldId2' },
          ],
          databaseId: 'databaseId',
        },
      },
      {
        form: {
          fields: [nonKeyField3],
          id: 'formId3',
          formType: FormType.RecipientFormType,
          folderId: 'folderId',
          databaseId: 'databaseId',
          name: 'formName3',
          code: '',
        },
        record: record3,
      },
    ]);
  });

  it('should skip merging if only 1 form', () => {
    const originalData: FormWithRecord<Recipient>[] = [
      {
        form: form1,
        record: record1,
      },
    ];
    const result = prettifyData(originalData);
    expect(result).toEqual([
      {
        form: {
          fields: [nonKeyField1],
          id: 'formId1',
          formType: FormType.RecipientFormType,
          folderId: 'folderId',
          databaseId: 'databaseId',
          name: 'formName1',
          code: '',
        },
        record: record1,
      },
    ]);
  });
});
