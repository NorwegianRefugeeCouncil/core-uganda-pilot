import React from 'react';
import { render, waitFor } from '@testing-library/react-native';

import { RecordsScreenContainer } from '../../src/components/screen_containers/RecordsScreenContainer';
import { Record } from '../../../client/lib/esm';
import { RECORD_ACTIONS } from '../../src/reducers/recordsReducers';

const mockFormId = 'a form id';
const mockDatabaseId = 'a database id';
const mockRoute = {
  params: {
    formId: mockFormId,
    databaseId: mockDatabaseId,
  },
};
const mockNavigation = {};
const mockState: any = {
  formsById: {},
};

const mockDispatch = jest.fn();

const mockRemoteRecord: Record = {
  id: '',
  databaseId: mockDatabaseId,
  formId: mockFormId,
  values: {},
  ownerId: '',
};

const mockRemoteRecords = Array(10)
  .fill(mockRemoteRecord)
  .map((r, i) => ({ ...r, id: i.toString() }));

const mockLocalRecord = '';
const mockLocalRecords = Array(10)
  .fill(mockLocalRecord)
  .map((_, i) => `local record #${i}`);

const mockRemoteData = {
  response: {
    items: mockRemoteRecords,
  },
};

mockState.formsById[mockFormId] = {
  records: mockRemoteRecords,
  localRecords: mockLocalRecords,
};

const mockProps: any = {
  route: mockRoute,
  navigation: mockNavigation,
  state: mockState,
  dispatch: mockDispatch,
};

const mockListRecords = jest.fn().mockResolvedValue(mockRemoteData);
jest.mock('../../src/utils/clients', () => () => ({
  listRecords: (object: any) => mockListRecords(object),
}));

describe(RecordsScreenContainer.name, () => {
  afterEach(() => {
    mockListRecords.mockClear();
  });

  test('renders correctly', async () => {
    const { toJSON } = await waitFor(() => render(<RecordsScreenContainer {...mockProps} />));
    expect(toJSON()).toMatchSnapshot();
  });

  test('fetches records and dispatches the appropriate action', async () => {
    await waitFor(() => render(<RecordsScreenContainer {...mockProps} />));

    const expectedListRecordsArg = {
      formId: mockFormId,
      databaseId: mockDatabaseId,
    };
    expect(mockListRecords).toHaveBeenCalledWith(expectedListRecordsArg);

    const expectedDispatchAction = {
      type: RECORD_ACTIONS.GET_RECORDS,
      payload: {
        formId: mockFormId,
        records: mockRemoteRecords,
      },
    };
    expect(mockDispatch).toHaveBeenCalledWith(expectedDispatchAction);
  });

  test('does not attempt to fetch if missing formId on route', async () => {
    const incompleteRoute = {
      ...mockRoute,
      params: {
        ...mockRoute.params,
        formId: null,
      },
    };
    await waitFor(() => render(<RecordsScreenContainer {...mockProps} route={incompleteRoute} />));

    expect(mockListRecords).not.toHaveBeenCalled();
  });

  test('does not attempt to fetch if missing databaseId on route', async () => {
    const incompleteRoute = {
      ...mockRoute,
      params: {
        ...mockRoute.params,
        databaseId: null,
      },
    };
    await waitFor(() => render(<RecordsScreenContainer {...mockProps} route={incompleteRoute} />));

    expect(mockListRecords).not.toHaveBeenCalled();
  });
});
