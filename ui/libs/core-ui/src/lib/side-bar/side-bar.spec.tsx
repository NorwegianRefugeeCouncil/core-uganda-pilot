import { render } from '@testing-library/react';

import SideBar from './side-bar';

describe('SideBar', () => {
  it('should render successfully', () => {
    const { baseElement } = render(<SideBar />);
    expect(baseElement).toBeTruthy();
  });
});
