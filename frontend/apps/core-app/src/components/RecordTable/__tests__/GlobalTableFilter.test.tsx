import React from 'react';

import { render } from '../../../testUtils/render';
import { GlobalTableFilter } from '../GlobalTableFilter';

it('should match the snapshot', () => {
  const { toJSON } = render(
    <GlobalTableFilter table={{ state: {}, preGlobalFilteredRows: [] }} />,
  );
  expect(toJSON()).toMatchSnapshot();
});
