import { render } from '../../../testUtils/render';
import { RecipientProfileScreen } from '../index';

it('should match the snapshot', () => {
  const { toJSON } = render(<RecipientProfileScreen />);
  expect(toJSON()).toMatchSnapshot();
});
