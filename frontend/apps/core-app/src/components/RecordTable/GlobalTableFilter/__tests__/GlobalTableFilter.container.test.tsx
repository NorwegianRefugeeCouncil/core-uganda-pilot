import React from 'react';

import { render } from '../../../../testUtils/render';
import { GlobalTableFilterContainer } from '../GlobalTableFilter.container';

describe('GlobalTableFilterContainer', () => {
  it('should match the snapshot', () => {
    const { toJSON } = render(
      <GlobalTableFilterContainer
        table={{
          state: {},
          setGlobalFilter: jest.fn(),
        }}
      />,
    );
    expect(toJSON()).toMatchSnapshot();
  });
});
