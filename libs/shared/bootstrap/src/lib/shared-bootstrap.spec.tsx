import { render } from '@testing-library/react';

import SharedBootstrap from './shared-bootstrap';

describe('SharedBootstrap', () => {
  it('should render successfully', () => {
    const { baseElement } = render(<SharedBootstrap />);
    expect(baseElement).toBeTruthy();
  });
});
