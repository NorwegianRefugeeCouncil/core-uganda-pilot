import { FieldKind } from 'core-api-client';

import { render } from '../../../testUtils/render';
import { SubformFieldValueComponent } from '../SubformFieldValue.component';

it('should match the snapshot', () => {
  const { toJSON } = render(
    <SubformFieldValueComponent
      header="header"
      labels={['label-1', 'label-2']}
      items={[
        [
          {
            fieldType: FieldKind.Text,
            value: 'test-1',
            formattedValue: 'format-test-1',
          },
          {
            fieldType: FieldKind.Reference,
            value: 'test-1',
            formattedValue: 'format-test-1',
          },
        ],
        [
          {
            fieldType: FieldKind.Text,
            value: 'test-2',
            formattedValue: 'format-test-2',
          },
          {
            fieldType: FieldKind.Reference,
            value: 'test-2',
            formattedValue: 'format-test-2',
          },
        ],
      ]}
    />,
  );
  expect(toJSON()).toMatchSnapshot();
});
