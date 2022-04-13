import { render } from '../../../testUtils/render';
import { SubformFieldValueComponent } from '../SubformFieldValue.component';

it('should match the snapshot', () => {
  const { toJSON } = render(
    <SubformFieldValueComponent
      header="header"
      columns={[
        {
          Header: 'header-0',
          accessor: 'sub-field-0',
        },
        {
          Header: 'header-1',
          accessor: 'sub-field-1',
        },
      ]}
      data={[
        {
          'sub-field-0': 'value-0-0',
          'sub-field-1': 'value-0-1',
        },
        {
          'sub-field-0': 'value-1-0',
          'sub-field-1': 'value-1-1',
        },
      ]}
    />,
  );
  expect(toJSON()).toMatchSnapshot();
});
