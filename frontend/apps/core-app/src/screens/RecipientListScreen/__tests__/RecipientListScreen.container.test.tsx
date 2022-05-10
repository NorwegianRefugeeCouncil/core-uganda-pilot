import React from 'react';
import { FormDefinition, FormType } from 'core-api-client';
import { waitFor } from '@testing-library/react-native';

import { render } from '../../../testUtils/render';
import { RecipientListScreenContainer } from '../RecipientListScreen.container';
import { mockNavigationProp } from '../../../testUtils/navigation';
import { makeField, makeForm } from '../../../testUtils/mockData';

const mockUseRecipientForms = jest.fn();

jest.mock('../../../contexts/RecipientForms', () => {
  return {
    useRecipientForms: () => mockUseRecipientForms(),
  };
});

jest.mock('../RecipientListScreen.component', () => {
  const { View, Text } = jest.requireActual('react-native');
  return {
    RecipientListScreenComponent: ({
      forms,
      isLoading,
    }: {
      forms: FormDefinition[];
      isLoading: boolean;
    }) => (
      <View>
        <Text testID="mock-data">{JSON.stringify(forms)}</Text>
        <Text>loading - {isLoading.toString()}</Text>
      </View>
    ),
  };
});

describe('RecipientListScreenContainer', () => {
  it('should call useRecipientForms', () => {
    render(
      <RecipientListScreenContainer
        navigation={mockNavigationProp}
        route={{ key: 'key', name: 'recipientsList' }}
      />,
    );
    expect(mockUseRecipientForms).toHaveBeenCalledTimes(1);
  });

  it('should render data', async () => {
    const f1 = makeField(1, true, false, { text: {} });
    const f2 = makeField(2, false, false, { text: {} });
    const f3 = makeField(3, true, false, { text: {} });
    const f4 = makeField(4, false, false, { text: {} });
    const form1 = makeForm(1, FormType.RecipientFormType, [f1, f2]);
    const form2 = makeForm(1, FormType.RecipientFormType, [f3, f4]);

    mockUseRecipientForms.mockReturnValue([form1, form2]);
    const { getByTestId, getByText } = render(
      <RecipientListScreenContainer
        navigation={mockNavigationProp}
        route={{ key: 'key', name: 'recipientsList' }}
      />,
    );

    await waitFor(() => {
      expect(
        JSON.parse(getByTestId('mock-data').children[0].toString()),
      ).toEqual([form1, form2]);
      expect(getByText('loading - false')).toBeTruthy();
    });
  });

  it('should render loading state', async () => {
    mockUseRecipientForms.mockReturnValue([]);
    const { getByTestId, getByText } = render(
      <RecipientListScreenContainer
        navigation={mockNavigationProp}
        route={{ key: 'key', name: 'recipientsList' }}
      />,
    );

    await waitFor(() => {
      expect(
        JSON.parse(getByTestId('mock-data').children[0].toString()),
      ).toEqual([]);
      expect(getByText('loading - true')).toBeTruthy();
    });
  });
});
