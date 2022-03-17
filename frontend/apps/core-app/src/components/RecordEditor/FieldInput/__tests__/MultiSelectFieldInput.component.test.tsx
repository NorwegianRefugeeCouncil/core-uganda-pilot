import { fireEvent, waitFor } from '@testing-library/react-native';

import { MultiSelectFieldInput } from '../MultiSelectFieldInput.component';
import { render } from '../../../../testUtils/render';
import { withFormContext } from '../../../../testUtils/withFormContext';

const makeDefaultValues = (value: string[]) => ({
  'form-id': {
    'field-id': value,
  },
});

const field = {
  id: 'field-id',
  name: 'Field name',
  description: 'Field description',
  code: '',
  required: false,
  key: false,
  fieldType: {
    multiSelect: {
      options: [
        { id: 'option-id-1', name: 'Option 1' },
        { id: 'option-id-2', name: 'Option 2' },
        { id: 'option-id-3', name: 'Option 3' },
      ],
    },
  },
};

describe('should match the snapshot', () => {
  it('no selected value', () => {
    const MultiSelectWithContext = withFormContext(
      MultiSelectFieldInput,
      makeDefaultValues([]),
    );
    const { toJSON } = render(
      <MultiSelectWithContext formId="form-id" field={field} />,
    );
    expect(toJSON()).toMatchSnapshot();
  });

  it('a single selected value', () => {
    const MultiSelectWithContext = withFormContext(
      MultiSelectFieldInput,
      makeDefaultValues(['option-id-2']),
    );
    const { toJSON } = render(
      <MultiSelectWithContext formId="form-id" field={field} />,
    );
    expect(toJSON()).toMatchSnapshot();
  });

  it('multiple selected values', () => {
    const MultiSelectWithContext = withFormContext(
      MultiSelectFieldInput,
      makeDefaultValues(['option-id-1', 'option-id-3']),
    );
    const { toJSON } = render(
      <MultiSelectWithContext formId="form-id" field={field} />,
    );
    expect(toJSON()).toMatchSnapshot();
  });
});

it('should handle on change', async () => {
  const onSubmitSpy = jest.fn();
  const MultiSelectWithContext = withFormContext(
    MultiSelectFieldInput,
    makeDefaultValues([]),
    onSubmitSpy,
  );
  const { getByTestId } = render(
    <MultiSelectWithContext formId="form-id" field={field} />,
  );
  const submitButton = getByTestId('with-form-context-submit-button');
  const modalToggle = getByTestId(
    'multi-select-field-input-modal-toggle-button',
  );
  const value = getByTestId('multi-select-field-input-value');

  fireEvent.press(submitButton);
  await waitFor(() => expect(onSubmitSpy).toHaveBeenCalledTimes(1));
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': [],
    },
  });
  expect(value.props.value).toEqual('');

  fireEvent.press(modalToggle);
  fireEvent.press(getByTestId('multi-select-field-input-option-0'));
  fireEvent.press(getByTestId('multi-select-field-input-option-1'));
  expect(value.props.value).toEqual('');
  fireEvent.press(getByTestId('multi-select-field-input-modal-submit'));
  expect(value.props.value).toEqual('Option 1, Option 2');
  fireEvent.press(submitButton);
  await waitFor(() => expect(onSubmitSpy).toHaveBeenCalledTimes(2));
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': ['option-id-1', 'option-id-2'],
    },
  });
});

it('should not update values when cancelling the modal', async () => {
  const onSubmitSpy = jest.fn();
  const MultiSelectWithContext = withFormContext(
    MultiSelectFieldInput,
    makeDefaultValues([]),
    onSubmitSpy,
  );
  const { getByTestId } = render(
    <MultiSelectWithContext formId="form-id" field={field} />,
  );
  const submitButton = getByTestId('with-form-context-submit-button');
  const modalToggle = getByTestId(
    'multi-select-field-input-modal-toggle-button',
  );
  const value = getByTestId('multi-select-field-input-value');

  fireEvent.press(submitButton);
  await waitFor(() => expect(onSubmitSpy).toHaveBeenCalledTimes(1));
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': [],
    },
  });
  expect(value.props.value).toEqual('');

  fireEvent.press(modalToggle);
  fireEvent.press(getByTestId('multi-select-field-input-option-0'));
  fireEvent.press(getByTestId('multi-select-field-input-option-1'));
  expect(value.props.value).toEqual('');
  fireEvent.press(getByTestId('multi-select-field-input-modal-cancel'));
  expect(value.props.value).toEqual('');
  fireEvent.press(submitButton);
  await waitFor(() => expect(onSubmitSpy).toHaveBeenCalledTimes(2));
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': [],
    },
  });
});
