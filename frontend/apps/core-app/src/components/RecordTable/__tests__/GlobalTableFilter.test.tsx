/**
 * @jest-environment jsdom
 */

import React from 'react';
import { fireEvent, waitFor } from '@testing-library/react';

import { render, renderWeb } from '../../../testUtils/render';
import { GlobalTableFilter } from '../GlobalTableFilter';

jest.mock('react-native/Libraries/Utilities/Platform', () => ({
  OS: 'web',
  select: () => null,
}));

describe('GlobalTableFilter', () => {
  it('should match the snapshot', () => {
    const { toJSON } = render(
      <GlobalTableFilter table={{ state: {}, preGlobalFilteredRows: [] }} />,
    );
    expect(toJSON()).toMatchSnapshot();
  });

  it('should call onChange handler', async () => {
    const setGlobalFilterSpy = jest.fn();

    const { getByPlaceholderText, debug } = renderWeb(
      <GlobalTableFilter
        table={{
          state: {},
          preGlobalFilteredRows: [],
          setGlobalFilter: setGlobalFilterSpy,
        }}
      />,
    );
    debug();

    const input = getByPlaceholderText('Search');
    fireEvent.input(input, 'newFilter');

    await waitFor(() =>
      expect(setGlobalFilterSpy).toHaveBeenCalledWith('newFilter'),
    );
  });
});
