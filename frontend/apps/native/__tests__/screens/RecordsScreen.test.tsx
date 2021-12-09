import React from 'react';
import { fireEvent, render, waitFor } from '@testing-library/react-native';

import RecordsScreen from '../../src/components/screens/RecordsScreen';
import testIds from '../../src/constants/testIds';
import routes from '../../src/constants/routes';

const mockFormId = 'a form id';
const mockDatabaseId = 'a database id';
const mockOwnerId = 'an owner id';
const mockValues = {};

const mockNavigate = jest.fn();

const mockRoute = {
  params: {
    formId: mockFormId,
  },
};

const mockNavigation = {
  navigate: mockNavigate,
};

const mockRemoteRecord = {
  databaseId: mockDatabaseId,
  formId: mockFormId,
  ownerId: mockOwnerId,
  values: mockValues,
};

const n_remoteRecords = 10;
const n_localRecords = 10;

const mockRemoteRecords = Array(n_remoteRecords)
  .fill(mockRemoteRecord)
  .map((r, i) => ({ ...r, id: `remote record #${i}` }));

const mockLocalRecord = '';
const mockLocalRecords = Array(n_localRecords)
  .fill(mockLocalRecord)
  .map((_, i) => `local record #${i}`);

const mockState: any = {
  formsById: {},
};

mockState.formsById[mockFormId] = {
  records: mockRemoteRecords,
  localRecords: mockLocalRecords,
};

const mockProps: any = {
  isLoading: false,
  navigation: mockNavigation,
  route: mockRoute,
  state: mockState,
};

describe(RecordsScreen.name, () => {
  afterEach(() => {
    mockNavigate.mockClear();
  });

  test('renders correctly', async () => {
    const { toJSON } = await waitFor(() => render(<RecordsScreen {...mockProps} />));
    expect(toJSON()).toMatchSnapshot();
  });

  test('renders a list of remote records', async () => {
    const { findAllByTestId } = await waitFor(() => render(<RecordsScreen {...mockProps} />));
    const records = await findAllByTestId(testIds.remoteRecord);
    expect(records.length).toEqual(n_remoteRecords);
  });

  test('renders a list of local records', async () => {
    const { findAllByTestId } = await waitFor(() => render(<RecordsScreen {...mockProps} />));
    const records = await findAllByTestId(testIds.localRecord);
    expect(records.length).toEqual(n_localRecords);
  });

  test('shows a message if loading', async () => {
    const { findByText } = await waitFor(() => render(<RecordsScreen {...mockProps} isLoading />));
    expect(await findByText('Loading...')).toBeTruthy();
  });

  test("navigates to a remote record's view page on press", async () => {
    const { findAllByA11yLabel } = await waitFor(() => render(<RecordsScreen {...mockProps} />));

    const records = await findAllByA11yLabel('view record');

    await fireEvent.press(records[0]);

    const expectedNavigateArgs = [routes.viewRecord.name, { recordId: mockRemoteRecords[0].id, formId: mockFormId }];

    await waitFor(() => expect(mockNavigate).toHaveBeenCalledWith(...expectedNavigateArgs));
  });

  test("navigates to a local record's add record page on press", async () => {
    const { findAllByA11yLabel } = await waitFor(() => render(<RecordsScreen {...mockProps} />));

    const localRecords = await findAllByA11yLabel('add record');

    await fireEvent.press(localRecords[0]);

    const expectedNavigateArgs = [routes.addRecord.name, { recordId: mockLocalRecords[0], formId: mockFormId }];

    await waitFor(() => expect(mockNavigate).toHaveBeenCalledWith(...expectedNavigateArgs));
  });
});
