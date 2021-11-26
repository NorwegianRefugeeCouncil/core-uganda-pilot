import React from 'react';
import { FormsScreenContainer } from '../../../src/components/screen_containers/FormsScreenContainer';
import { act, create } from 'react-test-renderer';
import { NativeBaseTestWrapper } from '../../../src/utils/NativeBaseTestWrapper';

let renderer;
describe(FormsScreenContainer.name, () => {
    const mockNavigation = jest.fn();
    const mockRoute = {};
    const mockListForms = jest.fn();
    jest.mock('../../../src/utils/useApiClient', () => {
        return { useApiClient: () => ({ listForms: mockListForms }) };
    });
    beforeAll(() => {
        act(() => {
            renderer = create(
                <NativeBaseTestWrapper>
                    <FormsScreenContainer
                        navigation={mockNavigation}
                        route={mockRoute}
                    />
                </NativeBaseTestWrapper>
            );
        });
    });

    test('fetches forms', () => {
        expect(mockListForms).toHaveBeenCalled();
    });
});
