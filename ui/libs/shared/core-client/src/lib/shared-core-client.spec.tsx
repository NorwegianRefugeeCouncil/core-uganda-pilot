import { render } from '@testing-library/react';

import SharedCoreClient from './shared-core-client';

describe('SharedCoreClient', () => {
  it('should render successfully', () => {
    const { baseElement } = render(<SharedCoreClient />);
    expect(baseElement).toBeTruthy();
  });
});
