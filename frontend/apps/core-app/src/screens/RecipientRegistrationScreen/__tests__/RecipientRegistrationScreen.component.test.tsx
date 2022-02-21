import { render } from '../../../testUtils/render';
import { RecipientRegistrationScreen } from '../index';

it('should match the snapshot', () => {
  const { toJSON } = render(<RecipientRegistrationScreen />);
  expect(toJSON()).toMatchSnapshot();
});
