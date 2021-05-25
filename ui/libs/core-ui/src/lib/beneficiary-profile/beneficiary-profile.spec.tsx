import { render } from '@testing-library/react';

import BeneficiaryProfile from './beneficiary-profile';

describe('BeneficiaryProfile', () => {
  it('should render successfully', () => {
    const { baseElement } = render(<BeneficiaryProfile />);
    expect(baseElement).toBeTruthy();
  });
});
