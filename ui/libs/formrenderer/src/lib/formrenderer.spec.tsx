import { render } from '@testing-library/react';

import Formrenderer from './formrenderer';

describe('Formrenderer', () => {
  it('should render successfully', () => {
    const { baseElement } = render(<Formrenderer />);
    expect(baseElement).toBeTruthy();
  });
});
