import * as React from "react";
import NavigationBar from '../../src/components/NavigationBar'
import { render, cleanup } from '@testing-library/react-native';

describe('NavigationBar', () => {
    afterEach(cleanup);

    it('renders correctly', () => {
        const tree = render(<NavigationBar options={{ title: "test" }} />).toJSON();
        expect(tree).toMatchSnapshot();
    })

})
