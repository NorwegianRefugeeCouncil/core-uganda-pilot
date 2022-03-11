import { FormDefinition, FormType, Record } from 'core-api-client';

import { render } from '../../../testUtils/render';
import { buildDefaultRecord } from '../buildDefaultRecord';
import { RecipientRegistrationScreenComponent } from '../RecipientRegistrationScreen.component';

jest.mock('../../../components/RecordEditor', () => {
  const { View, Text } = jest.requireActual('react-native');
  return {
    RecordEditor: ({
      form,
      record,
      onChange,
    }: {
      form: FormDefinition;
      record: Record;
      onChange: () => void;
    }) => (
      <View>
        <Text>{JSON.stringify(form)}</Text>
        <Text>{JSON.stringify(record)}</Text>
      </View>
    ),
  };
});

it('should match the snapshot', () => {
  const makeForm = (i: number): FormDefinition => ({
    id: `form-id-${i}`,
    code: 'form-code',
    databaseId: 'database-id',
    folderId: 'folder-id',
    name: `form-name-${i}`,
    formType: FormType.DefaultFormType,
    fields: [
      {
        id: `field-id-${i}`,
        name: `field-name-${i}`,
        code: '',
        description: '',
        required: false,
        key: false,
        fieldType: { text: {} },
      },
    ],
  });

  const { toJSON } = render(
    <RecipientRegistrationScreenComponent
      forms={[makeForm(1), makeForm(2)]}
      records={[
        buildDefaultRecord(makeForm(1)),
        buildDefaultRecord(makeForm(2)),
      ]}
      onSubmit={() => {}}
      onCancel={() => {}}
    />,
  );
  expect(toJSON()).toMatchSnapshot();
});
