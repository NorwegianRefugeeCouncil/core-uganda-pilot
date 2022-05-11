import React from 'react';
import { FormType } from 'core-api-client';
import { act, fireEvent } from '@testing-library/react-native';

import { render } from '../../../testUtils/render';
import { RecipientListTableContainer } from '../RecipientListTableContainer';
import { makeField, makeForm, makeRecord } from '../../../testUtils/mockData';
import * as hooks from '../../../hooks/useAPICall';
import * as mapRecordsUtils from '../mapRecordsToRecipientTableData';
import * as createColumnsUtils from '../createTableColumns';
import { formsClient } from '../../../clients/formsClient';
import { routes } from '../../../constants/routes';

const mockNavigate = jest.fn();
const mockSetGlobalFilter = jest.fn();

jest.mock('@react-navigation/native', () => {
  const actualNav = jest.requireActual('@react-navigation/native');
  return {
    ...actualNav,
    useNavigation: () => ({
      navigate: mockNavigate,
    }),
  };
});

jest.mock('react-table', () => {
  const actualTable = jest.requireActual('react-table');
  return {
    ...actualTable,
    useTable: () => ({
      setGlobalFilter: mockSetGlobalFilter,
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
      prepareRow: jest.fn(),
    }),
  };
});

describe('RecipientListTableContainer', () => {
  const mapRecordsToRecipientTableDataSpy = jest.spyOn(
    mapRecordsUtils,
    'mapRecordsToRecipientTableData',
  );
  const createTableColumnsSpy = jest.spyOn(
    createColumnsUtils,
    'createTableColumns',
  );

  const f1 = makeField(1, true, false, { text: {} });
  const f2 = makeField(2, false, false, { text: {} });
  const form = makeForm(1, FormType.RecipientFormType, [f1, f2]);

  const record1 = makeRecord(1, form);
  const record2 = makeRecord(2, form);

  const data = [
    [
      { form, record: record1 },
      { form, record: record2 },
    ],
  ];

  it('should call api and prepare data', () => {
    const useAPICallSpy = jest
      .spyOn(hooks, 'useAPICall')
      .mockImplementation((_, __, ___) => [
        () => Promise.resolve(),
        { data, loading: false, error: null },
      ]);
    render(<RecipientListTableContainer form={form} filter="filter" />);

    expect(useAPICallSpy).toHaveBeenCalledWith(
      formsClient.Recipient.list,
      [{ formId: form.id, databaseId: form.databaseId }],
      true,
    );
    expect(mapRecordsToRecipientTableDataSpy).toHaveBeenCalled();
    expect(createTableColumnsSpy).toHaveBeenCalled();
    expect(mockSetGlobalFilter).toHaveBeenCalledWith('filter');
  });

  it('should call navigate on click handler', () => {
    jest
      .spyOn(hooks, 'useAPICall')
      .mockImplementation((_, __, ___) => [
        () => Promise.resolve(),
        { data, loading: false, error: null },
      ]);
    const { getByTestId } = render(
      <RecipientListTableContainer form={form} filter="filter" />,
    );

    const row = getByTestId('recipient-list-table-row-row1');

    act(() => {
      fireEvent.press(row);
    });

    expect(mockNavigate).toHaveBeenCalledWith(routes.recipientsProfile.name, {
      recordId: 'recordId',
      formId: form.id,
      databaseId: form.databaseId,
    });
  });
});
