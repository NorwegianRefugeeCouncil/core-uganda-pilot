import { render } from '../../../testUtils/render';
import { RecipientRegistrationScreenComponent } from '../RecipientRegistrationScreen.component';

it('should match the snapshot', () => {
  const { toJSON } = render(
    <RecipientRegistrationScreenComponent
      route={{ name: 'RecipientRegistration', params: {}, key: 'key' }}
    />,
  );
  expect(toJSON()).toMatchSnapshot();
});
