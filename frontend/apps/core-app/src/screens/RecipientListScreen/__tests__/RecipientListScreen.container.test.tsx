import React from 'react';

import { render } from '../../../testUtils/render';
import { RecipientListScreenContainer } from '../RecipientListScreen.container';
import { mockNavigationProp } from '../../../testUtils/navigation';

const mockUseRecipientForms = jest.fn();

jest.mock('../../../contexts/RecipientForms', () => {
  return {
    useRecipientForms: () => mockUseRecipientForms(),
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
});
