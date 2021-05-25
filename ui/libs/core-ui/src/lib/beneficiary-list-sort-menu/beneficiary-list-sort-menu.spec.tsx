import { render } from '@testing-library/react';

import BeneficiaryListSortMenu from './beneficiary-list-sort-menu';

describe('BeneficiaryListSortMenu', () => {
  it('should render successfully', () => {
    const { baseElement } = render(<BeneficiaryListSortMenu />);
    expect(baseElement).toBeTruthy();
  });
});
