import React from 'react';

import { render } from '../../../testUtils/render';
import { RecipientListTableHeaderCell } from '../RecipientListTableHeaderCell';

it('should match the snapshot', () => {
  const { toJSON } = render(
    <RecipientListTableHeaderCell
      column={{
        getSortByToggleProps: () => ({ onClick: jest.fn() }),
        isSorted: true,
        isSortedDesc: false,
        render: jest.fn(),
      }}
    />,
  );
  expect(toJSON()).toMatchSnapshot();
});
