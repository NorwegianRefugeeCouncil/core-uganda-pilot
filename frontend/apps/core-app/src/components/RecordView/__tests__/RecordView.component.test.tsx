import { Column } from 'react-table';
import { FieldKind } from 'core-api-client';

import { render } from '../../../testUtils/render';
import {
  NormalisedBasicField,
  NormalisedFieldValue,
} from '../RecordView.types';
import { RecordViewComponent } from '../RecordView.component';

jest.mock('../FieldValue.component', () => {
  const { View, Text } = jest.requireActual('react-native');
  return {
    FieldValueComponent: ({ item }: { item: NormalisedBasicField }) => (
      <View>
        <Text>{JSON.stringify(item)}</Text>
      </View>
    ),
  };
});

jest.mock('../SubformFieldValue.component', () => {
  const { View, Text } = jest.requireActual('react-native');
  return {
    SubformFieldValueComponent: ({
      header,
      data,
      columns,
    }: {
      header: string;
      data: Record<string, string>[];
      columns: Column<Record<string, string>>[];
    }) => (
      <View>
        <Text>{header}</Text>
        <Text>{JSON.stringify(data)}</Text>
        <Text>{JSON.stringify(columns)}</Text>
      </View>
    ),
  };
});

it('should match the snapshot', () => {
  const fieldValues: NormalisedFieldValue[] = [
    {
      label: 'label-1',
      fieldType: FieldKind.Text,
      value: 'test-1',
      formattedValue: 'format-test-1',
      key: false,
    },
    {
      fieldType: FieldKind.SubForm,
      header: 'header-1',
      columns: [
        { Header: 'label-1', accessor: 'label-1' },
        { Header: 'label-2', accessor: 'label-2' },
      ],
      data: [
        {
          'label-1': 'value-1',
          'label-2': 'value-2',
        },
        {
          'label-1': 'value-3',
          'label-2': 'value-4',
        },
      ],
      key: false,
    },
  ];

  const { toJSON } = render(<RecordViewComponent fieldValues={fieldValues} />);
  expect(toJSON()).toMatchSnapshot();
});
