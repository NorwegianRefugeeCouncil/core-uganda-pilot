import { FormType, FormWithRecord } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';

import { render } from '../../../testUtils/render';
import { RecipientProfileScreenComponent } from '../RecipientProfileScreen.component';

jest.mock('../../../components/Recipient/RecipientViewer', () => {
  const { View, Text } = jest.requireActual('react-native');
  return {
    RecipientViewer: ({ data }: { data: FormWithRecord<Recipient>[] }) => (
      <View>
        <Text>{JSON.stringify(data)}</Text>
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

describe('should match the snapshot', () => {
  it('default', () => {
    const { toJSON } = render(
      <RecipientProfileScreenComponent
        data={data}
        isLoading={false}
        error={null}
      />,
    );
    expect(toJSON()).toMatchSnapshot();
  });

  it('loading', () => {
    const { toJSON } = render(
      <RecipientProfileScreenComponent data={data} isLoading error={null} />,
    );
    expect(toJSON()).toMatchSnapshot();
  });

  it('error', () => {
    const { toJSON } = render(
      <RecipientProfileScreenComponent
        data={data}
        isLoading={false}
        error="error"
      />,
    );
    expect(toJSON()).toMatchSnapshot();
  });
});
