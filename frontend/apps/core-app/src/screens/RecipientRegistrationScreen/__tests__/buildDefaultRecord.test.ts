import { FieldDefinition, FieldType, FormType } from 'core-api-client';

import { buildDefaultRecord } from '../buildDefaultRecord';

it('should set the correct ids', () => {
  const record = buildDefaultRecord({
    id: 'form-id',
    databaseId: 'database-id',
    code: 'form-code',
    folderId: 'folder-id',
    name: 'form-name',
    formType: FormType.DefaultFormType,
    fields: [],
  });
  expect(record.id).toBe('');
  expect(record.databaseId).toBe('database-id');
  expect(record.formId).toBe('form-id');
  expect(record.ownerId).toBeUndefined();
});

it('should set the correct field defaults', () => {
  const makeField = (fieldType: FieldType, i: number): FieldDefinition => ({
    id: `field-id-${i}`,
    code: '',
    name: `field-name-${i}`,
    description: '',
    required: false,
    key: false,
    fieldType,
  });

  const record = buildDefaultRecord({
    id: 'form-id',
    databaseId: 'database-id',
    code: 'form-code',
    folderId: 'folder-id',
    name: 'form-name',
    formType: FormType.DefaultFormType,
    fields: [
      makeField({ text: {} }, 1),
      makeField({ multilineText: {} }, 2),
      makeField(
        { reference: { databaseId: 'database-id', formId: 'other-form-id' } },
        3,
      ),
      makeField({ subForm: { fields: [] } }, 4),
      makeField({ date: {} }, 5),
      makeField({ month: {} }, 6),
      makeField({ week: {} }, 7),
      makeField({ quantity: {} }, 8),
      makeField({ singleSelect: { options: [] } }, 9),
      makeField({ multiSelect: { options: [] } }, 10),
      makeField({ checkbox: {} }, 11),
    ],
  });

  expect(record.values).toEqual([
    { fieldId: 'field-id-1', value: '' },
    { fieldId: 'field-id-2', value: '' },
    { fieldId: 'field-id-3', value: null },
    { fieldId: 'field-id-4', value: [] },
    { fieldId: 'field-id-5', value: null },
    { fieldId: 'field-id-6', value: null },
    { fieldId: 'field-id-7', value: null },
    { fieldId: 'field-id-8', value: '' },
    { fieldId: 'field-id-9', value: null },
    { fieldId: 'field-id-10', value: [] },
    { fieldId: 'field-id-11', value: 'false' },
  ]);
});
