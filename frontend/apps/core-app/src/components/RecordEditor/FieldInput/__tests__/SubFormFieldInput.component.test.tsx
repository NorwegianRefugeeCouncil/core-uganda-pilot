import { Column } from 'react-table';
import { fireEvent, waitFor } from '@testing-library/react-native';
import { FormDefinition, FormType } from 'core-api-client';

import { render } from '../../../../testUtils/render';
import { withFormContext } from '../../../../testUtils/withFormContext';
import { SubFormFieldInput } from '../SubFormFieldInput.component';

// This normally returns a random id causing the snapshots to change on every run
jest.mock('@react-aria/ssr/dist/main', () => ({
  ...jest.requireActual('@react-aria/ssr/dist/main'),
  useSSRSafeId: () => 'react-aria-generated-id',
}));

jest.mock('../../../SubFormTable', () => {
  const { View, Text } = jest.requireActual('react-native');
  return {
    SubFormTable: ({
      data,
      columns,
    }: {
      data: Record<string, string>[];
      columns: Column<Record<string, string>>[];
    }) => (
      <View>
        <Text testID="subformtable-data">{JSON.stringify(data)}</Text>
        <Text testID="subformtable-columns">{JSON.stringify(columns)}</Text>
      </View>
    ),
  };
});

const field: FormDefinition['fields'][number] = {
  id: 'field-id',
  name: 'field-name',
  description: 'field-description',
  code: '',
  fieldType: {
    subForm: {
      fields: [
        {
          id: 'sub-field-id-0',
          name: 'sub-field-name-0',
          description: 'sub-field-description-0',
          code: '',
          fieldType: { text: {} },
          required: false,
          key: false,
        },
        {
          id: 'sub-field-id-1',
          name: 'sub-field-name-1',
          description: 'sub-field-description-1',
          code: '',
          fieldType: { text: {} },
          required: false,
          key: false,
        },
      ],
    },
  },
  required: false,
  key: false,
};

const form: FormDefinition = {
  id: 'form-id',
  name: 'form-name',
  code: '',
  formType: FormType.DefaultFormType,
  databaseId: 'database-id',
  folderId: 'folder-id',
  fields: [field],
};

const makeDefaultFormValues = (
  numberOfSubRecords: number,
): {
  'form-id': { 'field-id': Record<string, string>[] };
} => ({
  'form-id': {
    'field-id': new Array(numberOfSubRecords).fill(null).map((_, i) => ({
      'sub-field-id-0': `sub-value-0-${i}`,
      'sub-field-id-1': `sub-value-1-${i}`,
    })),
  },
});

describe('should match the snapshot', () => {
  it('no sub records', () => {
    const SubFromWithContext = withFormContext(
      SubFormFieldInput,
      makeDefaultFormValues(0),
    );

    const { toJSON } = render(<SubFromWithContext form={form} field={field} />);

    expect(toJSON()).toMatchSnapshot();
  });

  it('sub records', () => {
    const SubFromWithContext = withFormContext(
      SubFormFieldInput,
      makeDefaultFormValues(2),
    );

    const { toJSON } = render(<SubFromWithContext form={form} field={field} />);

    expect(toJSON()).toMatchSnapshot();
  });
});

it('should add a sub record', async () => {
  const onSubmitSpy = jest.fn();
  const SubFromWithContext = withFormContext(
    SubFormFieldInput,
    makeDefaultFormValues(0),
    onSubmitSpy,
  );

  const { getByTestId, getAllByTestId } = render(
    <SubFromWithContext form={form} field={field} />,
  );

  expect(getByTestId('sub-form-field-input-empty')).toBeTruthy();
  fireEvent.press(getByTestId('sub-form-field-input-open-modal-button'));
  const inputs = getAllByTestId('text-field-input');
  fireEvent.changeText(inputs[0], 'value 1');
  fireEvent.changeText(inputs[1], 'value 2');
  fireEvent.press(getByTestId('sub-form-field-input-modal-add-button'));

  await waitFor(() => {
    expect(getByTestId('subformtable-data')).toBeTruthy();
    expect(getByTestId('subformtable-columns')).toBeTruthy();
  });

  expect(
    JSON.parse(getByTestId('subformtable-data').children[0] as string),
  ).toEqual([
    {
      'sub-field-id-0': 'value 1',
      'sub-field-id-1': 'value 2',
    },
  ]);

  fireEvent.press(getByTestId('with-form-context-submit-button'));
  await waitFor(() => expect(onSubmitSpy).toBeCalledTimes(1));
  expect(onSubmitSpy).toBeCalledWith({
    'form-id': {
      'field-id': [
        {
          'sub-field-id-0': 'value 1',
          'sub-field-id-1': 'value 2',
        },
      ],
    },
  });
});

it('should not add a record when cancelling the modal', async () => {
  const onSubmitSpy = jest.fn();
  const SubFromWithContext = withFormContext(
    SubFormFieldInput,
    makeDefaultFormValues(0),
    onSubmitSpy,
  );

  const { getByTestId, getAllByTestId } = render(
    <SubFromWithContext form={form} field={field} />,
  );

  expect(getByTestId('sub-form-field-input-empty')).toBeTruthy();
  fireEvent.press(getByTestId('sub-form-field-input-open-modal-button'));
  const inputs = getAllByTestId('text-field-input');
  fireEvent.changeText(inputs[0], 'value 1');
  fireEvent.changeText(inputs[1], 'value 2');
  fireEvent.press(getByTestId('sub-form-field-input-modal-cancel-button'));
  fireEvent.press(getByTestId('with-form-context-submit-button'));
  await waitFor(() => expect(onSubmitSpy).toBeCalledTimes(1));
  expect(onSubmitSpy).toBeCalledWith({
    'form-id': {
      'field-id': [],
    },
  });
});
