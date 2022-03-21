import { fireEvent, waitFor } from '@testing-library/react-native';

import { render } from '../../../../testUtils/render';
import { withFormContext } from '../../../../testUtils/withFormContext';
import { CheckboxFieldInput } from '../CheckboxFieldInput.component';

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
  fieldType: { checkbox: {} },
};

describe('should match the snapshot', () => {
  it('checked by default', () => {
    const CheckboxWithContext = withFormContext(
      CheckboxFieldInput,
      makeDefaultValues('true'),
    );
    const { toJSON } = render(
      <CheckboxWithContext formId="form-id" field={field} />,
    );
    expect(toJSON()).toMatchSnapshot();
  });

  it('unchecked by default', () => {
    const CheckboxWithContext = withFormContext(
      CheckboxFieldInput,
      makeDefaultValues('false'),
    );
    const { toJSON } = render(
      <CheckboxWithContext formId="form-id" field={field} />,
    );
    expect(toJSON()).toMatchSnapshot();
  });

  it('invalid value', () => {
    const CheckboxWithContext = withFormContext(
      CheckboxFieldInput,
      makeDefaultValues('wrong-value'),
    );
    const { toJSON } = render(
      <CheckboxWithContext formId="form-id" field={field} />,
    );
    expect(toJSON()).toMatchSnapshot();
  });
});

it('should handle on change', async () => {
  const onSubmitSpy = jest.fn();
  const CheckboxWithContext = withFormContext(
    CheckboxFieldInput,
    makeDefaultValues('false'),
    onSubmitSpy,
  );
  const { getByTestId } = render(
    <CheckboxWithContext formId="form-id" field={field} />,
  );
  const submitButton = getByTestId('with-form-context-submit-button');
  const checkbox = getByTestId('checkbox-field-input');

  fireEvent.press(submitButton);
  await waitFor(() => expect(onSubmitSpy).toHaveBeenCalledTimes(1));
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': 'false',
    },
  });

  fireEvent.press(checkbox);
  fireEvent.press(submitButton);
  await waitFor(() => expect(onSubmitSpy).toHaveBeenCalledTimes(2));
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': 'true',
    },
  });

  fireEvent.press(checkbox);
  fireEvent.press(submitButton);
  await waitFor(() => expect(onSubmitSpy).toHaveBeenCalledTimes(3));
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': 'false',
    },
  });
});
