import { render } from '@testing-library/react';

import BeneficiaryForm from './beneficiary-form';

describe('BeneficiaryForm', () => {
  it('should render successfully', () => {
    const { baseElement } = render(<BeneficiaryForm />);
    expect(baseElement).toBeTruthy();
  });
});
