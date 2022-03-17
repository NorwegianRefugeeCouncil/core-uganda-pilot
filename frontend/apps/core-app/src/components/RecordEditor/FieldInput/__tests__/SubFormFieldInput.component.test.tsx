import { fireEvent, waitFor } from '@testing-library/react-native';
import { FormDefinition, FormType, Record } from 'core-api-client';

import { render } from '../../../../testUtils/render';
import { withFormContext } from '../../../../testUtils/withFormContext';
import { SubFormFieldInput } from '../SubFormFieldInput.component';

// This normally returns a random id causing the snapshots to change on every run
jest.mock('@react-aria/ssr/dist/main', () => ({
  ...jest.requireActual('@react-aria/ssr/dist/main'),
  useSSRSafeId: () => 'react-aria-generated-id',
}));

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
  'form-id': { 'field-id': Record[] };
} => ({
  'form-id': {
    'field-id': new Array(numberOfSubRecords).fill(null).map((_, i) => ({
      id: `sub-record-id-${i}`,
      name: `sub-record-name-${i}`,
      formId: 'field-id',
      databaseId: 'database-id',
      ownerId: 'owner-id',
      code: '',
      values: [],
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
  await waitFor(() =>
    expect(getAllByTestId('sub-form-field-input-value').length).toBe(2),
  );
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
