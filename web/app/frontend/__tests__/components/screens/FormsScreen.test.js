import React from 'react';
import { create } from 'react-test-renderer';

import { FormsScreen } from '../../../src/components/screens/FormsScreen';
import testIds from '../../../src/constants/testIds';
import { v4 as uuid } from 'uuid';
import { NativeBaseTestWrapper } from '../../../src/utils/NativeBaseTestWrapper';

let renderer;
describe(FormsScreen.name, () => {
    const fakeForms = Array(10)
        .fill(null)
        .map(() => ({ id: uuid(), name: 'snazzy name' }));
    const mockNavigation = {
        navigate: () => {},
    };
    beforeAll(() => {
        renderer = create(
            <NativeBaseTestWrapper>
                <FormsScreen
                    forms={fakeForms}
                    isLoading={false}
                    navigation={mockNavigation}
                />
            </NativeBaseTestWrapper>
        );
    });

    test('renders correctly', () => {
        const tree = renderer.toJSON();
        expect(tree).toMatchSnapshot();
    });

    test('renders a list of forms', () => {
        const forms = renderer.root.findAll(
            el => el.type === 'View' && el.props.testID === testIds.formListItem
        );
        expect(forms.length).toEqual(fakeForms.length);
    });
});
