import {
  FieldDefinition,
  FieldType,
  FormDefinition,
  FormType,
  Record,
} from 'core-api-client';

import * as ReactHookFormTransformer from '../ReactHookFormTransformer';

const makeForm = (i: number): FormDefinition => {
  const makeField = (fieldType: FieldType, j: number): FieldDefinition => ({
    id: `field-id-${i}-${j}`,
    code: '',
    name: `field-name-${i}-${j}`,
    description: '',
    required: false,
    key: false,
    fieldType,
  });

  return {
    id: `form-id-${i}`,
    databaseId: `database-id-${i}`,
    code: `form-code-${i}`,
    folderId: `folder-id-${i}`,
    name: `form-name-${i}`,
    formType: FormType.DefaultFormType,
    fields: [
      makeField({ text: {} }, 1),
      makeField({ multilineText: {} }, 2),
      makeField(
        {
          reference: {
            databaseId: `database-id-${i}`,
            formId: `other-form-id-${i}`,
          },
        },
        3,
      ),
    ],
  };
};

const makeRecord = (i: number): Record => ({
  id: `record-id-${i}`,
  databaseId: `database-id-${i}`,
  formId: `form-id-${i}`,
  ownerId: undefined,
  values: [
    { fieldId: `field-id-${i}-1`, value: `value-${i}-1` },
    { fieldId: `field-id-${i}-2`, value: `value-${i}-2` },
    { fieldId: `field-id-${i}-3`, value: `value-${i}-3` },
  ],
});

describe('toReactHookForm', () => {
  it('should convert a form and record to the correct format', () => {
    const form = makeForm(1);
    const record = makeRecord(1);
    expect(
      ReactHookFormTransformer.toReactHookForm([{ form, record }]),
    ).toEqual({
      [form.id]: {
        [form.fields[0].id]: record.values[0].value,
        [form.fields[1].id]: record.values[1].value,
        [form.fields[2].id]: record.values[2].value,
      },
    });
  });

  it('should convert multiple forms and records to the correct format', () => {
    const form1 = makeForm(1);
    const form2 = makeForm(2);
    const record1 = makeRecord(1);
    const record2 = makeRecord(2);
    expect(
      ReactHookFormTransformer.toReactHookForm([
        { form: form1, record: record1 },
        { form: form2, record: record2 },
      ]),
    ).toEqual({
      [form1.id]: {
        [form1.fields[0].id]: record1.values[0].value,
        [form1.fields[1].id]: record1.values[1].value,
        [form1.fields[2].id]: record1.values[2].value,
      },
      [form2.id]: {
        [form2.fields[0].id]: record2.values[0].value,
        [form2.fields[1].id]: record2.values[1].value,
        [form2.fields[2].id]: record2.values[2].value,
      },
    });
  });
});

describe('fromReactHookForm', () => {
  it('should convert a form and record to the correct format', () => {
    const form = makeForm(1);
    const record = makeRecord(1);
    expect(
      ReactHookFormTransformer.fromReactHookForm([{ form, record }], {
        [form.id]: {
          [form.fields[0].id]: record.values[0].value,
          [form.fields[1].id]: record.values[1].value,
          [form.fields[2].id]: record.values[2].value,
        },
      }),
    ).toEqual([{ form, record }]);
  });

  it('should convert multiple forms and records to the correct format', () => {
    const form1 = makeForm(1);
    const form2 = makeForm(2);
    const record1 = makeRecord(1);
    const record2 = makeRecord(2);
    expect(
      ReactHookFormTransformer.fromReactHookForm(
        [
          { form: form1, record: record1 },
          { form: form2, record: record2 },
        ],
        {
          [form1.id]: {
            [form1.fields[0].id]: record1.values[0].value,
            [form1.fields[1].id]: record1.values[1].value,
            [form1.fields[2].id]: record1.values[2].value,
          },
          [form2.id]: {
            [form2.fields[0].id]: record2.values[0].value,
            [form2.fields[1].id]: record2.values[1].value,
            [form2.fields[2].id]: record2.values[2].value,
          },
        },
      ),
    ).toEqual([
      { form: form1, record: record1 },
      { form: form2, record: record2 },
    ]);
  });
});
