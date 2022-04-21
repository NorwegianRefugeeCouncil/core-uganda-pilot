import React from 'react';

import { render } from '../../../testUtils/render';
import { RecordTableHeaderCell } from '../RecordTableHeaderCell';

it('should match the snapshot', () => {
  const { toJSON } = render(
    <RecordTableHeaderCell
      column={{
        getSortByToggleProps: () => ({ onClick: jest.fn() }),
        isSorted: true,
        isSortedDesc: false,
        render: jest.fn()
      }}
    />,
  );
  expect(toJSON()).toMatchSnapshot();
});
