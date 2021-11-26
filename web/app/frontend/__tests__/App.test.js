import {create} from 'react-test-renderer';
import React from 'react';

import App from '../App.tsx';
import { testInset } from '../src/utils/NativeBaseTestWrapper';

let renderer

describe('App', () => {
    beforeAll(() => {
        renderer = create(<App inset={testInset}/>)
    })

    it('renders correctly', () => {
        const tree = renderer.toJSON();
        expect(tree).toMatchSnapshot();
    });

    it('has one child', () => {
        const tree = renderer.toJSON();
        console.log(tree)
        expect(tree.children?.length).toBe(1);
    });
});
