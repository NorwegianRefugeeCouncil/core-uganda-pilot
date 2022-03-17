import { act, fireEvent, waitFor } from '@testing-library/react-native';
import { FormDefinition, Record } from 'core-api-client';

import { formsClient } from '../../../../clients/formsClient';
import { render } from '../../../../testUtils/render';
import { withFormContext } from '../../../../testUtils/withFormContext';
import { ReferenceFieldInput } from '../ReferenceFieldInput.component';

const makeDefaultValues = (value: string | null) => ({
  'form-id': {
    'field-id': value,
  },
});

const makeRecord = (i: number): Record => ({
  id: `reference-record-id-${i}`,
  formId: 'reference-form-id',
  databaseId: 'reference-database-id',
  ownerId: undefined,
  values: [],
});

const field: FormDefinition['fields'][number] = {
  id: 'field-id',
  name: 'field-name',
  description: 'field-description',
  code: '',
  required: false,
  key: false,
  fieldType: {
    reference: {
      formId: 'reference-form-id',
      databaseId: 'reference-database-id',
    },
  },
};

const recordListResponse = {
  response: {
    items: [makeRecord(1), makeRecord(2), makeRecord(3)],
  },
  request: {
    formId: 'reference-form-id',
    databaseId: 'reference-database-id',
  },
  status: 'ok',
  statusCode: 200,
  success: true,
  error: undefined,
};

afterEach(() => {
  jest.clearAllMocks();
});

it('should match the snapshot', async () => {
  const ReferenceWithContext = withFormContext(
    ReferenceFieldInput,
    makeDefaultValues(null),
  );
  const recordListSpy = jest.spyOn(formsClient.Record, 'list');
  recordListSpy.mockImplementationOnce(() =>
    Promise.resolve(recordListResponse),
  );

  const { toJSON } = render(
    <ReferenceWithContext formId="form-id" field={field} />,
  );

  await waitFor(() => expect(recordListSpy).toHaveBeenCalledTimes(1));

  expect(toJSON()).toMatchSnapshot();
});

it('should fetch the reference records', async () => {
  const onSubmitSpy = jest.fn();
  const ReferenceWithContext = withFormContext(
    ReferenceFieldInput,
    makeDefaultValues(null),
    onSubmitSpy,
  );

  const recordListSpy = jest.spyOn(formsClient.Record, 'list');
  recordListSpy.mockImplementationOnce(() =>
    Promise.resolve(recordListResponse),
  );

  const { getByTestId } = render(
    <ReferenceWithContext formId="form-id" field={field} />,
  );

  await waitFor(() => expect(recordListSpy).toHaveBeenCalledTimes(1));
  expect(recordListSpy).toHaveBeenCalledWith({
    formId: 'reference-form-id',
    databaseId: 'reference-database-id',
  });
  expect(getByTestId('reference-field-input-option-0')).toBeTruthy();
  expect(getByTestId('reference-field-input-option-1')).toBeTruthy();
  expect(getByTestId('reference-field-input-option-2')).toBeTruthy();
});

it('should handle on change', async () => {
  const onSubmitSpy = jest.fn();
  const ReferenceWithContext = withFormContext(
    ReferenceFieldInput,
    makeDefaultValues(null),
    onSubmitSpy,
  );

  const recordListSpy = jest.spyOn(formsClient.Record, 'list');
  recordListSpy.mockImplementationOnce(() =>
    Promise.resolve(recordListResponse),
  );

  const { getByTestId } = render(
    <ReferenceWithContext formId="form-id" field={field} />,
  );

  await waitFor(() =>
    expect(getByTestId('reference-field-input-option-0')).toBeTruthy(),
  );

  const submitButton = getByTestId('with-form-context-submit-button');

  await act(() => fireEvent.press(submitButton));
  expect(onSubmitSpy).toHaveBeenCalledTimes(1);
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': null,
    },
  });

  const select = getByTestId('reference-field-input');
  fireEvent.press(select);
  fireEvent.press(getByTestId('reference-field-input-option-0'));
  expect(select.props.value).toEqual('reference-record-id-1');
  await act(() => fireEvent.press(submitButton));
  expect(onSubmitSpy).toHaveBeenCalledTimes(2);
  expect(onSubmitSpy).toHaveBeenCalledWith({
    'form-id': {
      'field-id': 'reference-record-id-1',
    },
  });
});
