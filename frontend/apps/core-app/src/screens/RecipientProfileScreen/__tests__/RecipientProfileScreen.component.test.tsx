import { FormDefinition, FormType, Record } from 'core-api-client';

import { render } from '../../../testUtils/render';
import { RecipientProfileScreenComponent } from '../RecipientProfileScreen.component';

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
  const { toJSON } = render(
    <RecipientProfileScreenComponent data={data} isLoading={false} />,
  );
  expect(toJSON()).toMatchSnapshot();
});
it('should match the snapshot while loading', () => {
  const { toJSON } = render(
    <RecipientProfileScreenComponent data={data} isLoading />,
  );
  expect(toJSON()).toMatchSnapshot();
});
it('should match the snapshot with error', () => {
  const { toJSON } = render(
    <RecipientProfileScreenComponent
      data={data}
      isLoading={false}
      error="error"
    />,
  );
  expect(toJSON()).toMatchSnapshot();
});
