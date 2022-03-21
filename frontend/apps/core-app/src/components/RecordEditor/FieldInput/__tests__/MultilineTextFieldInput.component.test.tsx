import { fireEvent, waitFor } from '@testing-library/react-native';

import { render } from '../../../../testUtils/render';
import { withFormContext } from '../../../../testUtils/withFormContext';
import { MultilineTextFieldInput } from '../MultilineTextFieldInput.component';

const makeDefaultValues = (value: string) => ({
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
  fieldType: { multilineText: {} },
};

it('should match the snapshot', () => {
  const TextWithContext = withFormContext(
    MultilineTextFieldInput,
    makeDefaultValues('an input value'),
  );
  const { toJSON } = render(<TextWithContext formId="form-id" field={field} />);
  expect(toJSON()).toMatchSnapshot();
});

it('should handle on change', async () => {
  const onSubmitSpy = jest.fn();
  const TextWithContext = withFormContext(
    MultilineTextFieldInput,
    makeDefaultValues(''),
    onSubmitSpy,
  );
  const { getByTestId } = render(
    <TextWithContext formId="form-id" field={field} />,
  );
  const submitButton = getByTestId('with-form-context-submit-button');
  const text = getByTestId('multiline-text-field-input');

  fireEvent.press(submitButton);
  await waitFor(() => expect(onSubmitSpy).toHaveBeenCalledTimes(1));
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': '',
    },
  });

  fireEvent.changeText(text, 'a new value');
  fireEvent.press(submitButton);
  await waitFor(() => expect(onSubmitSpy).toHaveBeenCalledTimes(2));
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': 'a new value',
    },
  });
});
