import { fireEvent } from '@testing-library/react-native';
import { Linking } from 'react-native';

import { render } from '../../../testUtils/render';
import { LargeNavHeaderComponent } from '../LargeNavHeader.component';
import * as logout from '../logout';

it('should match the snapshot', () => {
  const { toJSON } = render(
    <LargeNavHeaderComponent
      route={{
        key: '',
        name: '',
        params: {},
      }}
    />,
  );

  expect(toJSON()).toMatchSnapshot();
});

it('should render the root route header', () => {
  const { getByTestId } = render(
    <LargeNavHeaderComponent
      route={{
        key: '',
        name: '',
        params: {},
      }}
    />,
  );

  expect(getByTestId('nav-header')).toHaveTextContent('Beneficiaries');
});

it('should render the recipients list route header', () => {
  const { getByTestId } = render(
    <LargeNavHeaderComponent
      route={{
        key: '',
        name: 'recipientsRoot',
        params: { screen: 'recipientsList' },
      }}
    />,
  );

  expect(getByTestId('nav-header')).toHaveTextContent('Beneficiaries');
});

it('should render the recipients registration route header', () => {
  const { getByTestId } = render(
    <LargeNavHeaderComponent
      route={{
        key: '',
        name: 'recipientsRoot',
        params: { screen: 'recipientsRegistration' },
      }}
    />,
  );

  expect(getByTestId('nav-header')).toHaveTextContent(
    'Beneficiary Registration',
  );
});

it('should call the logout function', () => {
  jest
    .spyOn(Linking, 'openURL')
    .mockImplementationOnce(() => Promise.resolve());

  const logoutSpy = jest.spyOn(logout, 'logout').mockImplementation(() => {});

  const { getByText } = render(
    <LargeNavHeaderComponent
      route={{
        key: '',
        name: 'recipientsRoot',
        params: { screen: 'recipientsList' },
      }}
    />,
  );

  const logoutButton = getByText('Log out');
  fireEvent.press(logoutButton);

  expect(logoutSpy).toHaveBeenCalled();
});
