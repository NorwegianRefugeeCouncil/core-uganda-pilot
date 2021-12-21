import { fireEvent } from '@testing-library/react-native';

import { render } from '../../testUtils/render';

import { LoginScreenComponent } from './LoginScreen.component';

it('should match the snapshot', () => {
  const login = jest.fn();
  const { toJSON } = render(<LoginScreenComponent onLogin={login} />);
  expect(toJSON()).toMatchSnapshot();
});

it('should call login on click', () => {
  const login = jest.fn();
  const { getByText } = render(<LoginScreenComponent onLogin={login} />);
  fireEvent.press(getByText('Login'));
  expect(login).toHaveBeenCalled();
});
