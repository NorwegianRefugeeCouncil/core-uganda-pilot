import { FieldType } from 'core-api-client';

import { formatFieldValue } from '../formatFieldValue';

const makeField = (fieldType: FieldType) => ({
  id: 'field-id',
  code: 'field-code',
  name: 'field-name',
  description: 'field-description',
  required: false,
  key: false,
  fieldType,
});

const runDefaultTests = (
  fieldType: FieldType,
  expectedValues = ['test', 'one, two, three', '-'],
  testValues = ['test', ['one', 'two', 'three'], null],
) => {
  it('should render a string value', () => {
    const formattedValue = formatFieldValue(
      testValues[0],
      makeField(fieldType),
    );
    expect(formattedValue).toEqual(expectedValues[0]);
  });

  it('should render a string[] value', () => {
    const formattedValue = formatFieldValue(
      testValues[1],
      makeField(fieldType),
    );
    expect(formattedValue).toEqual(expectedValues[1]);
  });

  it('should render a null value', () => {
    const formattedValue = formatFieldValue(
      testValues[2],
      makeField(fieldType),
    );
    expect(formattedValue).toEqual(expectedValues[2]);
  });
};

describe('FieldKind - Text', () => runDefaultTests({ text: {} }));

describe('FieldKind - MultilineText', () =>
  runDefaultTests({ multilineText: {} }));

describe('FieldKind - Quantity', () => runDefaultTests({ multilineText: {} }));

describe('FieldKind - Reference', () =>
  runDefaultTests(
    {
      reference: {
        databaseId: 'database-id',
        formId: 'form-id',
      },
    },
    ['test', '-', '-'],
  ));

describe('FieldKind - SubForm', () =>
  runDefaultTests({ subForm: { fields: [] } }, [
    'Error: Subform fields should not reach this point',
    'Error: Subform fields should not reach this point',
    'Error: Subform fields should not reach this point',
  ]));

describe('FieldKind - Date', () =>
  runDefaultTests(
    {
      date: {},
    },
    ['test', '-', '-'],
  ));

describe('FieldKind - Week', () =>
  runDefaultTests(
    {
      week: {},
    },
    ['test', '-', '-'],
  ));

describe('FieldKind - Month', () =>
  runDefaultTests(
    {
      month: {},
    },
    ['test', '-', '-'],
  ));

describe('FieldKind - SingleSelect', () => {
  const field = {
    singleSelect: {
      options: [
        { id: 'opt-1', name: 'Option 1' },
        { id: 'opt-2', name: 'Option 2' },
      ],
    },
  };

  runDefaultTests(
    field,
    ['Option 1', '-', '-'],
    ['opt-1', ['opt-1', 'opt-2'], null],
  );

  it('should handle an incorrect value', () => {
    const formattedValue = formatFieldValue(
      'incorrect-value',
      makeField(field),
    );
    expect(formattedValue).toEqual('-');
  });
});

describe('FieldKind - MultiSelect', () => {
  const field = {
    multiSelect: {
      options: [
        { id: 'opt-1', name: 'Option 1' },
        { id: 'opt-2', name: 'Option 2' },
      ],
    },
  };

  runDefaultTests(
    field,
    ['Option 1', 'Option 1, Option 2', '-'],
    ['opt-1', ['opt-1', 'opt-2'], null],
  );

  it('should handle an incorrect value', () => {
    const formattedValue = formatFieldValue(
      ['opt-1', 'incorrect-value'],
      makeField(field),
    );
    expect(formattedValue).toEqual('Option 1, -');
  });
});

describe('FieldKind - Checkbox', () =>
  runDefaultTests(
    { checkbox: {} },
    ['Yes', 'No', 'No'],
    ['true', ['true', 'false'], null],
  ));
