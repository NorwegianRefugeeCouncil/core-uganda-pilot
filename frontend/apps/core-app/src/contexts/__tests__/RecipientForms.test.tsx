import { Text, View } from 'react-native';
import { renderHook } from '@testing-library/react-hooks';
import { waitFor } from '@testing-library/react-native';

import { RecipientFormsProvider, useRecipientForms } from '../RecipientForms';
import { formsClient } from '../../clients/formsClient';
import { makeForm } from '../../testUtils/mockData';
import { render } from '../../testUtils/render';

it('should return the correct value', async () => {
  const forms = [makeForm()];

  jest
    .spyOn(formsClient.Recipient, 'getRecipientForms')
    .mockResolvedValue(forms);

  const { result, waitForNextUpdate } = renderHook(() => useRecipientForms(), {
    wrapper: RecipientFormsProvider,
    initialProps: {
      children: <View />,
    },
  });

  expect(result.current).toEqual([]);

  await waitForNextUpdate();

  expect(result.current).toEqual(forms);
});

it('should display the children', async () => {
  const forms = [makeForm()];

  jest
    .spyOn(formsClient.Recipient, 'getRecipientForms')
    .mockResolvedValue(forms);

  const { getByText } = render(
    <RecipientFormsProvider>
      <View>
        <Text>Success</Text>
      </View>
    </RecipientFormsProvider>,
  );

  await waitFor(() => expect(getByText('Success')).toBeTruthy());
});

it('should display an error if the API call fails', async () => {
  jest
    .spyOn(formsClient.Recipient, 'getRecipientForms')
    .mockRejectedValue(new Error('API call failed'));

  const { getByText } = render(
    <RecipientFormsProvider>
      <View>
        <Text>Success</Text>
      </View>
    </RecipientFormsProvider>,
  );

  await waitFor(() =>
    expect(
      getByText('Error loading recipient forms: API call failed'),
    ).toBeTruthy(),
  );
});
