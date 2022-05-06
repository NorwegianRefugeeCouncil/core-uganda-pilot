import React from 'react';
import { FormType } from 'core-api-client';

import { render } from '../../../testUtils/render';
import { RecipientListTableContainer } from '../RecipientListTableContainer';
import { makeField, makeForm } from '../../../testUtils/mockData';

describe('RecipientListTableContainer', () => {
  const f1 = makeField(1, true, false, { text: {} });
  const f2 = makeField(2, false, false, { text: {} });
  const form = makeForm(1, FormType.RecipientFormType, [f1, f2]);
  it('should kjsdf', () => {
    render(<RecipientListTableContainer form={form} filter="filter" />);
    expect(toJSON()).toMatchSnapshot();
  });
});
