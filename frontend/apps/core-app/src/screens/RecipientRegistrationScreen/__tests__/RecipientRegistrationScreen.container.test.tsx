import { fireEvent, waitFor } from '@testing-library/react-native';
import { FormDefinition, FormType, FormWithRecord } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';

import { render } from '../../../testUtils/render';
import { buildDefaultRecord } from '../../../utils/buildDefaultRecord';
import { RecipientRegistrationScreenContainer } from '../RecipientRegistrationScreen.container';
import * as hooks from '../../../hooks/useAPICall';
import configuration from '../../../config';
import { formsClient } from '../../../clients/formsClient';
import { linkingConfig } from '../../../navigation/linking.config';

const makeForm = (i: number): FormDefinition => ({
  id: `form-id-${i}`,
  code: 'form-code',
  databaseId: 'database-id',
  folderId: 'folder-id',
  name: `form-name-${i}`,
  formType: FormType.DefaultFormType,
  fields: [
    {
      id: `field-id-${i}`,
      name: `field-name-${i}`,
      code: '',
      description: '',
      required: false,
      key: false,
      fieldType: { text: {} },
    },
  ],
});

const makeFormWithRecord = (i: number): FormWithRecord<Recipient> => {
  const form = makeForm(i);
  return {
    form,
    record: buildDefaultRecord(form),
  };
};

const forms = [makeForm(0), makeForm(1)];
const data = [makeFormWithRecord(0), makeFormWithRecord(1)];
const mockSubmitData = [makeFormWithRecord(0)];

const mockNavigate = jest.fn();

jest.mock('@react-navigation/native', () => {
  const actualNav = jest.requireActual('@react-navigation/native');
  return {
    ...actualNav,
    useNavigation: () => ({
      navigate: mockNavigate,
    }),
  };
});

jest.mock('../RecipientRegistrationScreen.component', () => {
  const { View, Text, Button } = jest.requireActual('react-native');
  return {
    RecipientRegistrationScreenComponent: ({
      mode,
      data,
      onSubmit,
      onCancel,
      error,
      loading,
    }: {
      mode: 'edit' | 'review';
      data: FormWithRecord<Recipient>[];
      onSubmit: (data: FormWithRecord<Recipient>[]) => void;
      onCancel: () => void;
      error: string | null;
      loading: boolean;
    }) => (
      <View>
        <Text>{mode}</Text>
        <Text testID="mock-data">{JSON.stringify(data)}</Text>
        <Text>{error || 'no-error'}</Text>
        <Text>loading - {loading.toString()}</Text>
        <Button onPress={() => onSubmit(mockSubmitData)} title="Submit" />
        <Button onPress={onCancel} title="Cancel" />
      </View>
    ),
  };
});

afterEach(() => {
  jest.clearAllMocks();
});

it('should call the api', () => {
  jest
    .spyOn(hooks, 'useAPICall')
    .mockImplementation(() => [
      () => Promise.resolve(),
      { data: forms, loading: false, error: null },
    ]);

  render(<RecipientRegistrationScreenContainer />);

  expect(hooks.useAPICall).toHaveBeenCalledTimes(2);
  expect(hooks.useAPICall).toHaveBeenCalledWith(
    formsClient.Form.getAncestors,
    [configuration.recipient.registrationForm.formId],
    true,
  );
});

describe('mode', () => {
  describe('edit', () => {
    it('should render', async () => {
      jest
        .spyOn(hooks, 'useAPICall')
        .mockImplementation(() => [
          () => Promise.resolve(),
          { data: forms, loading: false, error: null },
        ]);

      const { getByText, getByTestId } = render(
        <RecipientRegistrationScreenContainer />,
      );

      await waitFor(() => {
        expect(getByText('edit')).toBeTruthy();
        expect(
          JSON.parse(getByTestId('mock-data').children[0].toString()),
        ).toEqual(data);
        expect(getByText('no-error')).toBeTruthy();
        expect(getByText('loading - false')).toBeTruthy();
      });
    });

    it('should call onSubmit', async () => {
      jest
        .spyOn(hooks, 'useAPICall')
        .mockImplementation(() => [
          () => Promise.resolve(),
          { data: forms, loading: false, error: null },
        ]);

      const { getByText, getByTestId } = render(
        <RecipientRegistrationScreenContainer />,
      );

      fireEvent.press(getByText('Submit'));

      expect(getByText('review')).toBeTruthy();

      await waitFor(() =>
        expect(
          JSON.parse(getByTestId('mock-data').children[0].toString()),
        ).toEqual(mockSubmitData),
      );
    });

    it('should call onCancel', () => {
      jest
        .spyOn(hooks, 'useAPICall')
        .mockImplementation(() => [
          () => Promise.resolve(),
          { data: forms, loading: false, error: null },
        ]);

      const { getByText } = render(<RecipientRegistrationScreenContainer />);

      fireEvent.press(getByText('Cancel'));

      expect(mockNavigate).toHaveBeenCalledWith({
        key: linkingConfig.config.screens.Recipients,
      });
    });
  });
});

describe.only('review', () => {
  it('should render', async () => {
    jest
      .spyOn(hooks, 'useAPICall')
      .mockImplementation(() => [
        () => Promise.resolve(),
        { data: forms, loading: false, error: null },
      ]);

    const { getByText, getByTestId } = render(
      <RecipientRegistrationScreenContainer />,
    );

    fireEvent.press(getByText('Submit'));

    await waitFor(() => {
      expect(getByText('review')).toBeTruthy();
      expect(
        JSON.parse(getByTestId('mock-data').children[0].toString()),
      ).toEqual(mockSubmitData);
      expect(getByText('no-error')).toBeTruthy();
      expect(getByText('loading - false')).toBeTruthy();
    });
  });

  // it('should call onSubmit', () => {
  //   // Not implemented
  //   expect(true).toBeFalsy();
  // });

  it('should call onCancel', () => {
    jest
      .spyOn(hooks, 'useAPICall')
      .mockImplementation(() => [
        () => Promise.resolve(),
        { data: forms, loading: false, error: null },
      ]);

    const { getByText } = render(<RecipientRegistrationScreenContainer />);

    expect(getByText('edit')).toBeTruthy();

    fireEvent.press(getByText('Submit'));

    expect(getByText('review')).toBeTruthy();

    fireEvent.press(getByText('Cancel'));

    expect(getByText('edit')).toBeTruthy();
  });
});

it('should handle loading', () => {
  jest
    .spyOn(hooks, 'useAPICall')
    .mockImplementation(() => [
      () => Promise.resolve(),
      { data: forms, loading: true, error: null },
    ]);

  const { getByText } = render(<RecipientRegistrationScreenContainer />);

  expect(getByText('loading - true')).toBeTruthy();
});

it('should handle error', () => {
  jest
    .spyOn(hooks, 'useAPICall')
    .mockImplementation(() => [
      () => Promise.resolve(),
      { data: forms, loading: false, error: 'error' },
    ]);

  const { getByText } = render(<RecipientRegistrationScreenContainer />);

  expect(getByText('error')).toBeTruthy();
});
