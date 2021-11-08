import React from 'react';
import App from '../App.tsx';
import { render, cleanup } from '@testing-library/react-native';

describe('App', () => {
    afterEach(cleanup);

    it('has one child', () => {
        const tree = render(<App />).toJSON();
        expect(tree.children.length).toBe(1);
    });

    it('renders correctly', () => {
        const tree = render(<App />).toJSON();
        expect(tree).toMatchSnapshot();
    });

});
