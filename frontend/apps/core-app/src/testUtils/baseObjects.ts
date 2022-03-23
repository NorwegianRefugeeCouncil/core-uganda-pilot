import { FormDefinition, FormType, Record } from 'core-api-client';

export const baseForm: FormDefinition = {
  id: 'form1',
  code: '',
  name: 'form 1',
  databaseId: 'databaseId1',
  folderId: 'folderId1',
  formType: FormType.DefaultFormType,
  fields: [],
};

export const baseTextField = {
  id: 'field1',
  name: 'field 1',
  description: 'description 1',
  code: '',
  required: false,
  key: false,
  fieldType: {
    text: {},
  },
};

export const baseMultilineTextField = {
  id: 'field2',
  name: 'field 2',
  description: 'description 2',
  code: '',
  required: true,
  key: true,
  fieldType: {
    multilineText: {},
  },
};

export const baseRecord1: Record = {
  id: 'record1',
  values: [
    {
      value: 'text',
      fieldId: 'field1',
    },
    {
      value: 'multi line text',
      fieldId: 'field2',
    },
  ],
  formId: 'form1',
  databaseId: 'databaseId1',
  ownerId: undefined,
};
