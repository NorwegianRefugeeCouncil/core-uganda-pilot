import React from 'react';
import { Box } from 'native-base';
import { act, fireEvent, waitFor } from '@testing-library/react-native';

import { render } from '../../../testUtils/render';
import { RecipientListTableComponent } from '../RecipientListTableComponent';

const prepareRowMock = jest.fn();
const clickHandlerMock = jest.fn();
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
          render: jest.fn().mockImplementation((a) => {
            <Box>{a}</Box>;
          }),
          value: 'val1',
          row: 'row1',
          getCellProps: jest.fn(),
        },
        {
          column: 'col2',
          render: jest.fn().mockImplementation((a) => {
            <Box>{a}</Box>;
          }),
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
      render: jest.fn().mockImplementation((a) => {
        <Box>{a}</Box>;
      }),
    },
    {
      id: 'col2',
      Header: 'col2',
      accessor: 'col2',
      hidden: true,
      render: jest.fn().mockImplementation((a) => {
        <Box>{a}</Box>;
      }),
    },
  ],
  prepareRow: prepareRowMock,
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
        onItemClick={clickHandlerMock}
        title="title"
        error={null}
        loading
      />,
    );
    expect(prepareRowMock).toHaveBeenCalledTimes(1);
  });

  it('should call row click handler', () => {
    const { debug, getByTestId } = render(
      <RecipientListTableComponent
        table={table}
        onItemClick={clickHandlerMock}
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
    waitFor(() => expect(clickHandlerMock).toHaveBeenCalledWith('row1'));
  });
});
