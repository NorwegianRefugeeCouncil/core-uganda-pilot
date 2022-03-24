import { FormType, FormWithRecords } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';

export const formWithRecords: FormWithRecords<Recipient> = {
  form: {
    id: 'form1',
    name: 'name',
    databaseId: 'dbid',
    formType: FormType.RecipientFormType,
    folderId: '',
    fields: [
      {
        id: 'field1',
        name: 'fieldName',
        fieldType: { text: {} },
        key: false,
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
  records: [
    {
      id: 'id1',
      formId: 'formId',
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
      formId: 'formId',
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
  ],
};
