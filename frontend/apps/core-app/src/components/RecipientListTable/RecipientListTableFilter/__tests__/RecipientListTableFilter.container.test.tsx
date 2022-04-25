import React from 'react';

import { render } from '../../../../testUtils/render';
import { RecipientListTableFilterContainer } from '../RecipientListTableFilter.container';

describe('GlobalTableFilterContainer', () => {
  it('should match the snapshot', () => {
    const { toJSON } = render(
      <RecipientListTableFilterContainer
        table={{
          state: {},
          setGlobalFilter: jest.fn(),
        }}
      />,
    );
    expect(toJSON()).toMatchSnapshot();
  });
});
