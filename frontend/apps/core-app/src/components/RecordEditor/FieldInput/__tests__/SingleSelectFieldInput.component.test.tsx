import { fireEvent, waitFor } from '@testing-library/react-native';

import { SingleSelectFieldInput } from '../SingleSelectFieldInput.component';
import { render } from '../../../../testUtils/render';
import { withFormContext } from '../../../../testUtils/withFormContext';

const makeDefaultValues = (value: string | null) => ({
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
    singleSelect: {
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
    const SingleSelectWithContext = withFormContext(
      SingleSelectFieldInput,
      makeDefaultValues(null),
    );
    const { toJSON } = render(
      <SingleSelectWithContext formId="form-id" field={field} />,
    );
    expect(toJSON()).toMatchSnapshot();
  });

  it('selected value', () => {
    const SingleSelectWithContext = withFormContext(
      SingleSelectFieldInput,
      makeDefaultValues('option-id-2'),
    );
    const { toJSON } = render(
      <SingleSelectWithContext formId="form-id" field={field} />,
    );
    expect(toJSON()).toMatchSnapshot();
  });
});

it('should handle on change', async () => {
  const onSubmitSpy = jest.fn();
  const SingleSelectWithContext = withFormContext(
    SingleSelectFieldInput,
    makeDefaultValues(null),
    onSubmitSpy,
  );
  const { getByTestId } = render(
    <SingleSelectWithContext formId="form-id" field={field} />,
  );
  const submitButton = getByTestId('with-form-context-submit-button');

  fireEvent.press(submitButton);
  await waitFor(() => expect(onSubmitSpy).toHaveBeenCalledTimes(1));
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': null,
    },
  });

  fireEvent.press(getByTestId('single-select-field-input'));
  fireEvent.press(getByTestId('single-select-field-input-option-0'));
  fireEvent.press(submitButton);
  await waitFor(() => expect(onSubmitSpy).toHaveBeenCalledTimes(2));
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': 'option-id-1',
    },
  });
});
