import { render } from '../../../testUtils/render';
import { RecipientProfileScreenComponent } from '../RecipientProfileScreen.component';

it('should match the snapshot', () => {
  const { toJSON } = render(
    <RecipientProfileScreenComponent
      route={{ name: 'RecipientProfile', params: { id: 'id' }, key: 'key' }}
    />,
  );
  expect(toJSON()).toMatchSnapshot();
});
