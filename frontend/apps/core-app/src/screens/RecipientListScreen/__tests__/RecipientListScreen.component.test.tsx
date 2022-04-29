import React from 'react';

import { render } from '../../../testUtils/render';
import { RecipientListScreenComponent } from '../RecipientListScreen.component';
import { makeFormWithRecord } from '../../RecipientRegistrationScreen/__tests__/RecipientRegistrationScreen.component.test';

describe('RecipientListScreenComponent', () => {
  describe('should match the snapshot', () => {
    it('data', () => {
      const { toJSON } = render(
        <RecipientListScreenComponent
          data={[
            [makeFormWithRecord(1)],
            [makeFormWithRecord(2)],
            [makeFormWithRecord(3)],
          ]}
          forms={[]}
          onItemClick={jest.fn()}
          error={undefined}
          isLoading={false}
        />,
      );
      expect(toJSON()).toMatchSnapshot();
    });

    it('empty tables', () => {
      const { toJSON } = render(
        <RecipientListScreenComponent
          data={[]}
          forms={[]}
          onItemClick={jest.fn()}
          error={undefined}
          isLoading={false}
        />,
      );
      expect(toJSON()).toMatchSnapshot();
    });

    it('error', () => {
      const { toJSON } = render(
        <RecipientListScreenComponent
          data={null}
          forms={null}
          onItemClick={jest.fn()}
          error="error"
          isLoading={false}
        />,
      );
      expect(toJSON()).toMatchSnapshot();
    });

    it('loading', () => {
      const { toJSON } = render(
        <RecipientListScreenComponent
          data={null}
          forms={null}
          onItemClick={jest.fn()}
          error={undefined}
          isLoading
        />,
      );
      expect(toJSON()).toMatchSnapshot();
    });
  });
});
