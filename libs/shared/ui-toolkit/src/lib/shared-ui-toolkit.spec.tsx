import { render } from '@testing-library/react';

import SharedUiToolkit from './shared-ui-toolkit';

describe('SharedUiToolkit', () => {
  it('should render successfully', () => {
    const { baseElement } = render(<SharedUiToolkit />);
    expect(baseElement).toBeTruthy();
  });
});
