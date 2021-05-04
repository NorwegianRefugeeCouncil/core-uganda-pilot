import { render } from '@testing-library/react';

import Formbuilder from './formbuilder';

describe('Formbuilder', () => {
  it('should render successfully', () => {
    const { baseElement } = render(<Formbuilder />);
    expect(baseElement).toBeTruthy();
  });
});
