import React from 'react';
import Renderer from 'react-test-renderer';
import App from '../App.tsx';

describe('App', () => {
    
    it('<App /> has 1 child', () => {
      const tree = Renderer.create(<App />).toJSON();
      expect(tree.children.length).toBe(1);
    });

    it('renders correctly', () => {
        const tree = Renderer
          .create(<App />)
          .toJSON();
        expect(tree).toMatchSnapshot();
      });

  });