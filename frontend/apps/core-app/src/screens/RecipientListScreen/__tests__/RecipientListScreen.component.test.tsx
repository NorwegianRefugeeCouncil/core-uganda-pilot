import React from 'react';

import { render } from '../../../testUtils/render';
import { RecipientListScreenComponent } from '../RecipientListScreen.component';

it('should match the snapshot', () => {
  const { toJSON } = render(
    <RecipientListScreenComponent
      route={{ name: 'RecipientList', params: {}, key: 'key' }}
      handleItemClick={() => jest.fn()}
    />,
  );
  expect(toJSON()).toMatchSnapshot();
});
