import {render} from '@testing-library/react-native';
import React from 'react';

import FormsScreen from '../../../src/components/screens/FormsScreen';
import testIds from '../../../src/testIds';
import * as hooks from "../../../src/utils/useApiClient";

const fakeForms = {response: {items: [{id: "formId"}]}}
const formsPromise = Promise.resolve(fakeForms)
const mock = jest.fn();
hooks.useApiClient = mock
mock.mockImplementation(() => ({
    listForms: jest.fn(() => formsPromise)
}))

describe(FormsScreen.name, () => {
    const mockNavigation = {
        navigate: () => {
        }
    }

    test('renders correctly', async () => {
        const {toJSON, findAllByTestId} = render(<FormsScreen route={""} navigation={mockNavigation}/>)
        await findAllByTestId(testIds.formListItem)
        expect(toJSON()).toMatchSnapshot();
    });

    test('renders a list of forms', async () => {
        const {findAllByTestId} = render(<FormsScreen route={""} navigation={mockNavigation}/>)
        expect((await findAllByTestId(testIds.formListItem)).length).toEqual(1);
    });

});
