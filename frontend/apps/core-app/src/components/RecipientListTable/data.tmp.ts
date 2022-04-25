import { FormDefinition, FormType, FormWithRecord } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';

export const forms: FormDefinition[] = [
  {
    id: 'form1',
    name: 'name',
    databaseId: 'dbid',
    formType: FormType.RecipientFormType,
    folderId: '',
    fields: [
      {
        id: 'field1',
        name: 'fieldName',
        fieldType: { reference: { formId: 'form0', databaseId: 'dbid' } },
        key: true,
        code: '',
        required: false,
        description: '',
      },
      {
        id: 'field2',
        name: 'fieldName2',
        fieldType: { text: {} },
        key: false,
        code: '',
        required: false,
        description: '',
      },
    ],
    code: '',
  },
  {
    id: 'form0',
    name: 'name',
    databaseId: 'dbid',
    formType: FormType.RecipientFormType,
    folderId: '',
    fields: [
      {
        id: 'field01',
        name: 'fieldName01',
        fieldType: { text: {} },
        key: false,
        code: '',
        required: false,
        description: '',
      },
      {
        id: 'field02',
        name: 'fieldName02',
        fieldType: { text: {} },
        key: false,
        code: '',
        required: false,
        description: '',
      },
    ],
    code: '',
  },
];

export const recordsForm0: Recipient[] = [
  {
    id: 'id01',
    formId: 'form0',
    ownerId: undefined,
    databaseId: 'dbid',
    values: [
      {
        value: 'value01',
        fieldId: 'field01',
      },
      { value: 'value02', fieldId: 'field02' },
    ],
  },
  {
    id: 'id02',
    formId: 'form0',
    ownerId: undefined,
    databaseId: 'dbid',
    values: [
      {
        value: 'value03',
        fieldId: 'field01',
      },
      { value: 'value04', fieldId: 'field02' },
    ],
  },
];
export const recordsForm1 = [
  {
    id: 'id1',
    formId: 'form1',
    ownerId: undefined,
    databaseId: 'dbid',
    values: [
      {
        value: 'value1',
        fieldId: 'field1',
      },
      { value: 'value2', fieldId: 'field2' },
    ],
  },
  {
    id: 'id2',
    formId: 'form1',
    ownerId: undefined,
    databaseId: 'dbid',
    values: [
      {
        value: 'value3',
        fieldId: 'field1',
      },
      { value: 'value4', fieldId: 'field2' },
    ],
  },
];

export const data: FormWithRecord<Recipient>[][] = [
  [
    { form: forms[0], record: recordsForm1[0] },
    { form: forms[1], record: recordsForm0[0] },
  ],
  [
    { form: forms[0], record: recordsForm1[1] },
    { form: forms[1], record: recordsForm0[1] },
  ],
];
