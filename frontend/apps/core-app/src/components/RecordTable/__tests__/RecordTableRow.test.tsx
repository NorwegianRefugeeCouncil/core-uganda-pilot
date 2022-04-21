import React from 'react';

import { render } from '../../../testUtils/render';
import { RecordTableRow } from '../RecordTableRow';

it('should match the snapshot', () => {
  const { toJSON } = render(
    <RecordTableRow
      onRowClick={jest.fn()}
      row={{
        id: 'id',
        cells: [{ column: { width: '100' }, render: jest.fn() }],
      }}
    />,
  );
  expect(toJSON()).toMatchSnapshot();
});
