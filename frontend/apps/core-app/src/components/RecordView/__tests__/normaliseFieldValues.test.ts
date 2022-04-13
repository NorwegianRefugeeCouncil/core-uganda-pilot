import {
  FieldDefinition,
  FormType,
  FieldValue,
  FieldKind,
} from 'core-api-client';

import { normaliseFieldValues } from '../normaliseFieldValues';

const makeForm = (fields: FieldDefinition[]) => ({
  id: 'form-id',
  code: 'form-code',
  databaseId: 'database-id',
  folderId: 'folder-id',
  name: 'form-name',
  formType: FormType.DefaultFormType,
  fields,
});

const makeRecord = (values: FieldValue[]) => ({
  id: 'record-id',
  databaseId: 'database-id',
  formId: 'form-id',
  ownerId: undefined,
  values,
});

it('should normalise a basic field', () => {
  const form = makeForm([
    {
      id: 'field-id',
      name: 'field-name',
      code: 'field-code',
      description: 'field-description',
      required: false,
      key: false,
      fieldType: {
        text: {},
      },
    },
  ]);

  const record = makeRecord([
    {
      fieldId: 'field-id',
      value: 'field-value',
    },
  ]);

  const values = normaliseFieldValues(form, record);

  expect(values).toEqual([
    {
      key: false,
      label: 'field-name',
      value: 'field-value',
      fieldType: FieldKind.Text,
      formattedValue: 'field-value',
    },
  ]);
});

it('should normalise a subform field', () => {
  const form = makeForm([
    {
      id: 'field-id',
      name: 'field-name',
      code: 'field-code',
      description: 'field-description',
      required: false,
      key: false,
      fieldType: {
        subForm: {
          fields: [
            {
              id: 'sub-field-id',
              name: 'sub-field-name',
              code: 'sub-field-code',
              description: 'sub-field-description',
              required: false,
              key: false,
              fieldType: {
                text: {},
              },
            },
          ],
        },
      },
    },
  ]);

  const record = makeRecord([
    {
      fieldId: 'field-id',
      value: [
        [{ fieldId: 'sub-field-id', value: 'sub-field-value-1' }],
        [{ fieldId: 'sub-field-id', value: 'sub-field-value-2' }],
      ],
    },
  ]);

  const values = normaliseFieldValues(form, record);

  expect(values).toEqual([
    {
      key: false,
      fieldType: FieldKind.SubForm,
      header: 'field-name',
      columns: [
        {
          Header: 'sub-field-name',
          accessor: 'sub-field-id',
        },
      ],
      data: [
        { 'sub-field-id': 'sub-field-value-1' },
        { 'sub-field-id': 'sub-field-value-2' },
      ],
    },
  ]);
});

it('should error if a subform has a nested subform', () => {
  const form = makeForm([
    {
      id: 'field-id',
      name: 'field-name',
      code: 'field-code',
      description: 'field-description',
      required: false,
      key: false,
      fieldType: {
        subForm: {
          fields: [
            {
              id: 'sub-field-id',
              name: 'sub-field-name',
              code: 'sub-field-code',
              description: 'sub-field-description',
              required: false,
              key: false,
              fieldType: {
                subForm: { fields: [] },
              },
            },
          ],
        },
      },
    },
  ]);

  const record = makeRecord([
    {
      fieldId: 'field-id',
      value: [[{ fieldId: 'sub-field-id', value: [] }]],
    },
  ]);

  expect(() => normaliseFieldValues(form, record)).toThrowError();
});
