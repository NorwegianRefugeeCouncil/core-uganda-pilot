import * as React from 'react';
import Renderer, {ReactTestRendererJSON} from 'react-test-renderer';
import App from '../App.tsx';

describe('App', () => {

    it('<App /> has 1 child', () => {
      const tree = Renderer.create(<App />).toJSON() as ReactTestRendererJSON;
      expect(tree?.children?.length).toBe(1);
    });

    it('renders correctly', () => {
        const tree = Renderer
            .create(<App/>)
            .toJSON();
        expect(tree).toMatchSnapshot();
    });

});
