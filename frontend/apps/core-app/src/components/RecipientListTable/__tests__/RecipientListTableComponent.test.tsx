import React from 'react';
import { act, fireEvent, waitFor } from '@testing-library/react-native';

import { render } from '../../../testUtils/render';
import { RecipientListTableComponent } from '../RecipientListTableComponent';

const mockPrepareRow = jest.fn();
const mockClickHandler = jest.fn();
const table = {
  rows: [
    {
      id: 'row1',
      values: {
        recordId: 'recordId',
        col1: 'sdf',
        col2: 'sadf',
      },
      cells: [
        {
          column: 'col1',
          render: jest.fn().mockReturnValue(<div>box</div>),

          value: 'val1',
          row: 'row1',
          getCellProps: jest.fn(),
        },
        {
          column: 'col2',
          render: jest.fn().mockReturnValue(<div>box</div>),

          value: 'val2',
          row: 'row1',
          getCellProps: jest.fn(),
        },
      ],
    },
  ],
  columns: [
    {
      id: 'col1',
      Header: 'col1',
      accessor: 'col1',
      hidden: false,
      render: jest.fn().mockReturnValue(<div>box</div>),
    },
    {
      id: 'col2',
      Header: 'col2',
      accessor: 'col2',
      hidden: true,
      render: jest.fn().mockReturnValue(<div>box</div>),
    },
  ],
  prepareRow: mockPrepareRow,
};

describe('RecipientListTableComponent', () => {
  describe('should match the snapshot', () => {
    it('data', () => {
      const { toJSON } = render(
        <RecipientListTableComponent
          table={table}
          onItemClick={jest.fn()}
          title="title"
          error={null}
          loading={false}
        />,
      );
      expect(toJSON()).toMatchSnapshot();
    });
    it('error', () => {
      const { toJSON } = render(
        <RecipientListTableComponent
          table={{ rows: [], columns: [], prepareRow: jest.fn() }}
          onItemClick={jest.fn()}
          title="title"
          error="error message"
          loading={false}
        />,
      );
      expect(toJSON()).toMatchSnapshot();
    });
    it('loading', () => {
      const { toJSON } = render(
        <RecipientListTableComponent
          table={{ rows: [], columns: [], prepareRow: jest.fn() }}
          onItemClick={jest.fn()}
          title="title"
          error={null}
          loading
        />,
      );
      expect(toJSON()).toMatchSnapshot();
    });
  });

  it('should render rows', () => {
    render(
      <RecipientListTableComponent
        table={table}
        onItemClick={mockClickHandler}
        title="title"
        error={null}
        loading
      />,
    );
    expect(mockPrepareRow).toHaveBeenCalledTimes(1);
  });

  it('should call row click handler', () => {
    const { debug, getByTestId } = render(
      <RecipientListTableComponent
        table={table}
        onItemClick={mockClickHandler}
        title="title"
        error={null}
        loading
      />,
    );

    debug();
    const row = getByTestId('recipient-list-table-row-row1');
    act(() => {
      fireEvent.press(row);
    });
    waitFor(() => expect(mockClickHandler).toHaveBeenCalledWith('row1'));
  });
});
