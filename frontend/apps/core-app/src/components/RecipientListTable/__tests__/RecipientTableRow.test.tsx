import React from 'react';

import { render } from '../../../testUtils/render';
import { RecipientListTableRow } from '../RecipientListTableRow';

it('should match the snapshot', () => {
  const { toJSON } = render(
    <RecipientListTableRow
      onRowClick={jest.fn()}
      row={{
        id: 'id',
        cells: [{ column: { width: '100' }, render: jest.fn() }],
      }}
    />,
  );
  expect(toJSON()).toMatchSnapshot();
});
