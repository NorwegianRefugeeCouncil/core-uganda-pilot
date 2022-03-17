import { fireEvent, waitFor } from '@testing-library/react-native';

import { render } from '../../../../testUtils/render';
import { withFormContext } from '../../../../testUtils/withFormContext';
import { QuantityFieldInput } from '../QuantityFieldInput.component';

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
  fieldType: { quantity: {} },
};

it('should match the snapshot', () => {
  const QuantityWithContext = withFormContext(
    QuantityFieldInput,
    makeDefaultValues('an input value'),
  );
  const { toJSON } = render(
    <QuantityWithContext formId="form-id" field={field} />,
  );
  expect(toJSON()).toMatchSnapshot();
});

it('should handle on change', async () => {
  const onSubmitSpy = jest.fn();
  const QuantityWithContext = withFormContext(
    QuantityFieldInput,
    makeDefaultValues(''),
    onSubmitSpy,
  );
  const { getByTestId } = render(
    <QuantityWithContext formId="form-id" field={field} />,
  );
  const submitButton = getByTestId('with-form-context-submit-button');
  const input = getByTestId('quantity-field-input');

  fireEvent.press(submitButton);
  await waitFor(() => expect(onSubmitSpy).toHaveBeenCalledTimes(1));
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': '',
    },
  });

  fireEvent.changeText(input, '123');
  fireEvent.press(submitButton);
  await waitFor(() => expect(onSubmitSpy).toHaveBeenCalledTimes(2));
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': '123',
    },
  });
});

it('should not accept non-numeric characters', async () => {
  const onSubmitSpy = jest.fn();
  const QuantityWithContext = withFormContext(
    QuantityFieldInput,
    makeDefaultValues(''),
    onSubmitSpy,
  );
  const { getByTestId } = render(
    <QuantityWithContext formId="form-id" field={field} />,
  );
  const submitButton = getByTestId('with-form-context-submit-button');
  const input = getByTestId('quantity-field-input');

  fireEvent.press(submitButton);
  await waitFor(() => expect(onSubmitSpy).toHaveBeenCalledTimes(1));
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': '',
    },
  });

  fireEvent.changeText(input, 'a new value');
  fireEvent.press(submitButton);
  await waitFor(() => expect(onSubmitSpy).toHaveBeenCalledTimes(2));
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': '',
    },
  });
});

it('should accept decimal numbers', async () => {
  const onSubmitSpy = jest.fn();
  const QuantityWithContext = withFormContext(
    QuantityFieldInput,
    makeDefaultValues(''),
    onSubmitSpy,
  );
  const { getByTestId } = render(
    <QuantityWithContext formId="form-id" field={field} />,
  );
  const submitButton = getByTestId('with-form-context-submit-button');
  const input = getByTestId('quantity-field-input');

  fireEvent.press(submitButton);
  await waitFor(() => expect(onSubmitSpy).toHaveBeenCalledTimes(1));
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': '',
    },
  });

  fireEvent.changeText(input, '123.456');
  fireEvent.press(submitButton);
  await waitFor(() => expect(onSubmitSpy).toHaveBeenCalledTimes(2));
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': '123.456',
    },
  });
});
