import { FieldKind } from 'core-api-client';

import { render } from '../../../testUtils/render';
import { FieldValueComponent } from '../FieldValue.component';

describe('Text', () => {
  it('should match the snapshot', () => {
    const { toJSON } = render(
      <FieldValueComponent
        item={{
          fieldType: FieldKind.Text,
          value: 'test',
          formattedValue: 'format-test',
          label: 'label',
        }}
      />,
    );
    expect(toJSON()).toMatchSnapshot();
  });
});

describe('Link', () => {
  it('should match the snapshot', () => {
    const { toJSON } = render(
      <FieldValueComponent
        item={{
          fieldType: FieldKind.Reference,
          value: 'test',
          formattedValue: 'format-test',
          label: 'label',
        }}
      />,
    );
    expect(toJSON()).toMatchSnapshot();
  });
});
