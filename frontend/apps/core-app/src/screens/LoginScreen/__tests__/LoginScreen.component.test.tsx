import { fireEvent } from '@testing-library/react-native';

import { render } from '../../../testUtils/render';
import { LoginScreenComponent } from '../LoginScreen.component';

it('should match the snapshot', () => {
  const login = jest.fn();
  const { toJSON } = render(
    <LoginScreenComponent onLogin={login} isLoading={false} />,
  );
  expect(toJSON()).toMatchSnapshot();
});

it('should call login on click', () => {
  const login = jest.fn();
  const { getByTestId } = render(
    <LoginScreenComponent onLogin={login} isLoading={false} />,
  );
  fireEvent.press(getByTestId('login-button'));
  expect(login).toHaveBeenCalled();
});

it('should not call login on click when loading', () => {
  const login = jest.fn();
  const { getByTestId } = render(
    <LoginScreenComponent onLogin={login} isLoading />,
  );
  fireEvent.press(getByTestId('login-button'));
  expect(login).not.toHaveBeenCalled();
});
