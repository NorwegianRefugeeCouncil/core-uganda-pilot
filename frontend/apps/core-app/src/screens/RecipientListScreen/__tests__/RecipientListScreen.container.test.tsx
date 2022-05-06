import React from 'react';

import { render } from '../../../testUtils/render';
import { RecipientListScreenContainer } from '../RecipientListScreen.container';
import * as hooks from '../../../hooks/useAPICall';

const mockNavigate = jest.fn();
// const useRecipientFormsMock = jest.fn();

jest.mock('@react-navigation/native', () => {
  const actualNav = jest.requireActual('@react-navigation/native');
  return {
    ...actualNav,
    useNavigation: () => ({
      navigate: mockNavigate,
    }),
  };
});

// jest.mock('../../../contexts/RecipientForms', () => {
//   return {
//     useRecipientForms: useRecipientFormsMock,
//   };
// });

describe('RecipientListScreenContainer', () => {
  const useAPICallSpy = jest.spyOn(hooks, 'useAPICall');

  afterEach(() => {
    useAPICallSpy.mockReset();
  });

  it('should match the snapshot', () => {
    const { toJSON } = render(<RecipientListScreenContainer />);
    expect(toJSON()).toMatchSnapshot();
  });

  // it('should call useRecipientForms', () => {
  //   render(<RecipientListScreenContainer />);
  //   expect(useRecipientFormsMock).toHaveBeenCalledTimes(1);
  // });
});
