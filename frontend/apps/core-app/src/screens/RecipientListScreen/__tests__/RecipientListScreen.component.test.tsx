import React from 'react';
import { FormType } from 'core-api-client';

import { render } from '../../../testUtils/render';
import { RecipientListScreenComponent } from '../RecipientListScreen.component';
import { makeField, makeForm } from '../../../testUtils/mockData';

const field1 = makeField(1, false, false, { text: {} });
const field2 = makeField(2, true, false, { multilineText: {} });
const field3 = makeField(3, false, false, { text: {} });
const field4 = makeField(4, false, false, { multilineText: {} });
const form1 = makeForm(1, FormType.DefaultFormType, [field1, field2]);
const form2 = makeForm(2, FormType.DefaultFormType, [field3, field4]);

describe('RecipientListScreenComponent', () => {
  describe('should match the snapshot', () => {
    it('with forms', () => {
      const { toJSON } = render(
        <RecipientListScreenComponent
          forms={[form1, form2]}
          isLoading={false}
        />,
      );
      expect(toJSON()).toMatchSnapshot();
    });
    it('loading', () => {
      const { toJSON } = render(
        <RecipientListScreenComponent forms={[]} isLoading />,
      );
      expect(toJSON()).toMatchSnapshot();
    });
  });
});
