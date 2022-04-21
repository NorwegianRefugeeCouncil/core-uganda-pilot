/**
 * @jest-environment jsdom
 */

import React from 'react';
import { fireEvent, waitFor } from '@testing-library/react-native';
import { getByText } from '@testing-library/react';

import { render } from '../../../testUtils/render';
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

  it.only('should call onChange handler', async () => {
    const setGlobalFilterSpy = jest.fn();

    const { getByPlaceholderText, debug } = render(
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
    fireEvent.changeText(input, 'newFilter');

    await waitFor(() =>
      //   expect(fromReactHookFormSpy).toHaveBeenCalledWith(
      //     data,
      //     ReactHookFormTransformer.toReactHookForm(data),
      //   ),
      expect(setGlobalFilterSpy).toHaveBeenCalledWith('newFilter'),
    );
  });
});
