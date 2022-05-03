import { fireEvent, waitFor } from '@testing-library/react-native';
import { FormWithRecord } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';

import { render } from '../../../testUtils/render';
import * as ReactHookFormTransformer from '../../../utils/ReactHookFormTransformer';
import { RecipientRegistrationScreenComponent } from '../RecipientRegistrationScreen.component';
import { makeFormWithRecord } from '../../../testUtils/mockData';

jest.mock('../../../components/Recipient/RecipientEditor', () => {
  const { View, Text } = jest.requireActual('react-native');
  return {
    RecipientEditor: ({ data }: { data: FormWithRecord<Recipient>[] }) => (
      <View>
        <Text>{JSON.stringify(data)}</Text>
      </View>
    ),
  };
});

jest.mock('../../../components/Recipient/RecipientViewer', () => {
  const { View, Text } = jest.requireActual('react-native');
  return {
    RecipientViewer: ({ data }: { data: FormWithRecord<Recipient>[] }) => (
      <View>
        <Text>{JSON.stringify(data)}</Text>
      </View>
    ),
  };
});

afterEach(() => {
  jest.clearAllMocks();
});

describe('should match the snapshot', () => {
  it('edit mode', () => {
    const { toJSON } = render(
      <RecipientRegistrationScreenComponent
        mode="edit"
        data={[
          makeFormWithRecord(1),
          makeFormWithRecord(2),
          makeFormWithRecord(3),
        ]}
        onSubmit={() => {}}
        onCancel={() => {}}
        error={null}
        loading={false}
      />,
    );
    expect(toJSON()).toMatchSnapshot();
  });

  it('review mode', () => {
    const { toJSON } = render(
      <RecipientRegistrationScreenComponent
        mode="review"
        data={[
          makeFormWithRecord(1),
          makeFormWithRecord(2),
          makeFormWithRecord(3),
        ]}
        onSubmit={() => {}}
        onCancel={() => {}}
        error={null}
        loading={false}
      />,
    );
    expect(toJSON()).toMatchSnapshot();
  });

  it('error', () => {
    const { toJSON } = render(
      <RecipientRegistrationScreenComponent
        mode="edit"
        data={[
          makeFormWithRecord(1),
          makeFormWithRecord(2),
          makeFormWithRecord(3),
        ]}
        onSubmit={() => {}}
        onCancel={() => {}}
        error="error"
        loading={false}
      />,
    );
    expect(toJSON()).toMatchSnapshot();
  });

  it('loading', () => {
    const { toJSON } = render(
      <RecipientRegistrationScreenComponent
        mode="edit"
        data={[
          makeFormWithRecord(1),
          makeFormWithRecord(2),
          makeFormWithRecord(3),
        ]}
        onSubmit={() => {}}
        onCancel={() => {}}
        error={null}
        loading
      />,
    );
    expect(toJSON()).toMatchSnapshot();
  });
});

describe('onSubmit', () => {
  it('should transform the form values, then submit', async () => {
    const fromReactHookFormSpy = jest.spyOn(
      ReactHookFormTransformer,
      'fromReactHookForm',
    );

    const onSubmitSpy = jest.fn();

    const data = [
      makeFormWithRecord(1),
      makeFormWithRecord(2),
      makeFormWithRecord(3),
    ];

    const { getByTestId } = render(
      <RecipientRegistrationScreenComponent
        mode="edit"
        data={data}
        onSubmit={onSubmitSpy}
        onCancel={() => {}}
        error={null}
        loading={false}
      />,
    );

    const submitButton = getByTestId('recipient-registration-submit-button');
    fireEvent.press(submitButton);

    await waitFor(() =>
      expect(fromReactHookFormSpy).toHaveBeenCalledWith(
        data,
        ReactHookFormTransformer.toReactHookForm(data),
      ),
    );
    expect(onSubmitSpy).toHaveBeenCalledWith(data);
  });

  it('should submit without transforming values', () => {
    const fromReactHookFormSpy = jest.spyOn(
      ReactHookFormTransformer,
      'fromReactHookForm',
    );

    const onSubmitSpy = jest.fn();

    const data = [
      makeFormWithRecord(1),
      makeFormWithRecord(2),
      makeFormWithRecord(3),
    ];

    const { getByTestId } = render(
      <RecipientRegistrationScreenComponent
        mode="review"
        data={data}
        onSubmit={onSubmitSpy}
        onCancel={() => {}}
        error={null}
        loading={false}
      />,
    );

    const submitButton = getByTestId('recipient-registration-submit-button');
    fireEvent.press(submitButton);

    expect(fromReactHookFormSpy).not.toHaveBeenCalled();
    expect(onSubmitSpy).toHaveBeenCalledWith(data);
  });
});

it('should call the onCancel callback', () => {
  const onCancelSpy = jest.fn();

  const data = [
    makeFormWithRecord(1),
    makeFormWithRecord(2),
    makeFormWithRecord(3),
  ];

  const { getByTestId } = render(
    <RecipientRegistrationScreenComponent
      mode="edit"
      data={data}
      onSubmit={() => {}}
      onCancel={onCancelSpy}
      error={null}
      loading={false}
    />,
  );

  const cancelButton = getByTestId('recipient-registration-cancel-button');
  fireEvent.press(cancelButton);

  expect(onCancelSpy).toHaveBeenCalled();
});
