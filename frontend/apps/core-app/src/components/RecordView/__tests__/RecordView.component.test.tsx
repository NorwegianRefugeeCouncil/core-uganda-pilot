import { FieldKind } from 'core-api-client';

import { render } from '../../../testUtils/render';
import {
  NormalisedBasicField,
  NormalisedFieldValue,
  NormalisedSubFormFieldValue,
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
      labels,
      items,
    }: {
      header: string;
      labels: string[];
      items: NormalisedSubFormFieldValue[][];
    }) => (
      <View>
        <Text>{header}</Text>
        <Text>{JSON.stringify(labels)}</Text>
        <Text>{JSON.stringify(items)}</Text>
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
    },
    {
      fieldType: FieldKind.SubForm,
      header: 'header-1',
      labels: ['label-1', 'label-2'],
      values: [
        [
          {
            fieldType: FieldKind.Text,
            value: 'test-2',
            formattedValue: 'format-test-2',
          },
        ],
      ],
    },
  ];

  const { toJSON } = render(<RecordViewComponent fieldValues={fieldValues} />);
  expect(toJSON()).toMatchSnapshot();
});
