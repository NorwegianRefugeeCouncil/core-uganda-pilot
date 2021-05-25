import { render } from '@testing-library/react';

import BeneficiaryList from './beneficiary-list';

describe('BeneficiaryList', () => {
  it('should render successfully', () => {
    const { baseElement } = render(<BeneficiaryList />);
    expect(baseElement).toBeTruthy();
  });
});
