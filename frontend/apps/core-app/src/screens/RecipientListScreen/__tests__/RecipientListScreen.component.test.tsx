import { render } from '../../../testUtils/render';
import { RecipientListScreen } from '../index';

it('should match the snapshot', () => {
  const { toJSON } = render(<RecipientListScreen />);
  expect(toJSON()).toMatchSnapshot();
});
