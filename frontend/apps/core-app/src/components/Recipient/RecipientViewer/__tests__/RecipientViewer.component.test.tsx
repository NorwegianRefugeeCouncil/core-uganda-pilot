import { FormDefinition, FormType, Record } from 'core-api-client';

import { render } from '../../../../testUtils/render';
import { RecipientViewerComponent } from '../RecipientViewer.component';

jest.mock('../../../components/RecordView', () => {
  const { View, Text } = jest.requireActual('react-native');
  return {
    RecordView: ({
      form,
      record,
    }: {
      form: FormDefinition;
      record: Record;
    }) => (
      <View>
        <Text>{JSON.stringify(form)}</Text>
        <Text>{JSON.stringify(record)}</Text>
      </View>
    ),
  };
});

const data = [
  {
    form: {
      id: 'id',
      databaseId: 'databaseId',
      formType: FormType.RecipientFormType,
      name: 'name',
      code: 'code',
      folderId: 'folderId',
      fields: [],
    },
    record: {
      id: 'id',
      formId: 'formId',
      databaseId: 'databaseId',
      values: [],
      ownerId: undefined,
    },
  },
];

it('should match the snapshot', () => {
  const { toJSON } = render(<RecipientViewerComponent data={data} />);
  expect(toJSON()).toMatchSnapshot();
});

it('should throw an error with less than 2 forms', () => {
  expect(() => {
    render(<RecipientViewerComponent data={[data[0]]} />);
  }).toThrow();
});
