import { act, create } from 'react-test-renderer';
import React from 'react';

import { AddRecordScreen } from '../../../src/components/screens/AddRecordScreen';

const fakeFormName = 'snazzy name';
const fakeForm = { id: 'blabla', name: fakeFormName };
const mockNavigation = {
    navigate: () => {
    },
};
const mockRoute = { params: { formId: 'blabla', recordId: '' } };
const mockState = {};
const mockDispatch = jest.fn();
const mockGetForm = jest.fn(() => Promise.resolve({ response: fakeForm }));
jest.mock('../../../src/utils/useApiClient', () => {
    return { useApiClient: () => ({ getForm: mockGetForm }) };
});

describe(AddRecordScreen.name, () => {
    test('renders correctly', async () => {
        let root;
        act(() => {
                root = create(<AddRecordScreen control={} formState={} hasLocalData isConnected isLoading isWeb
                                               onSubmit={}/>);
            }
        );
        expect(root.toJSON()).toMatchSnapshot();
    });

    test('fetches a form', async () => {
        let root;
        act(() => {
                root = create(<AddRecordScreen control={} formState={} hasLocalData isConnected isLoading isWeb
                                               onSubmit={}/>);
            }
        );
        expect(mockGetForm).toHaveBeenCalled();
        expect(root.find(fakeFormName)).toBeTruthy();
    });
});
