import React from 'react';
import { act, fireEvent, waitFor } from '@testing-library/react-native';

import { render } from '../../../testUtils/render';
import { RecipientListTableHeaderCell } from '../RecipientListTableHeaderCell';

describe('RecipientTableHeaderCell', () => {
  describe('should match the snapshot', () => {
    it('unsorted', () => {
      const { toJSON } = render(
        <RecipientListTableHeaderCell
          column={{
            getSortByToggleProps: () => ({ onClick: jest.fn() }),
            isSorted: false,
            isSortedDesc: false,
            render: jest.fn(),
            toggleSortBy: jest.fn(),
          }}
        />,
      );
      expect(toJSON()).toMatchSnapshot();
    });
    it('sorted, descending', () => {
      const { toJSON } = render(
        <RecipientListTableHeaderCell
          column={{
            getSortByToggleProps: () => ({ onClick: jest.fn() }),
            isSorted: true,
            isSortedDesc: true,
            render: jest.fn(),
            toggleSortBy: jest.fn(),
          }}
        />,
      );
      expect(toJSON()).toMatchSnapshot();
    });
    it('sorted, ascending', () => {
      const { toJSON } = render(
        <RecipientListTableHeaderCell
          column={{
            getSortByToggleProps: () => ({ onClick: jest.fn() }),
            isSorted: true,
            isSortedDesc: false,
            render: jest.fn(),
            toggleSortBy: jest.fn(),
          }}
        />,
      );
      expect(toJSON()).toMatchSnapshot();
    });
  });

  it('should call press handler', () => {
    const toggleHandler = jest.fn();
    const { getByTestId } = render(
      <RecipientListTableHeaderCell
        column={{
          getSortByToggleProps: () => ({ onClick: jest.fn() }),
          isSorted: true,
          isSortedDesc: false,
          render: jest.fn(),
          toggleSortBy: toggleHandler,
        }}
      />,
    );
    const button = getByTestId('recipient-table-sort-button');
    act(() => {
      fireEvent.press(button);
    });
    waitFor(() => {
      expect(toggleHandler).toHaveBeenCalled();
    });
  });
});
