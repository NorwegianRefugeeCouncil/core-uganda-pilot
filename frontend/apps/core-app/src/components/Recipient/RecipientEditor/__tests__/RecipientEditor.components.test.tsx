import {
  FormDefinition,
  FormType,
  FormWithRecord,
  Record,
} from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';

import { render } from '../../../../testUtils/render';
import { buildDefaultRecord } from '../../../../utils/buildDefaultRecord';
import { RecipientEditorComponent } from '../RecipientEditor.component';

jest.mock('../../../RecordEditor', () => {
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

const makeFormWithRecord = (i: number): FormWithRecord<Recipient> => {
  const form = {
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
  };

  return {
    form,
    record: buildDefaultRecord(form),
  };
};

it('should match the snapshot', () => {
  const { toJSON } = render(
    <RecipientEditorComponent
      data={[
        makeFormWithRecord(1),
        makeFormWithRecord(2),
        makeFormWithRecord(3),
      ]}
    />,
  );
  expect(toJSON()).toMatchSnapshot();
});

it('should throw an error with less than 2 forms', () => {
  expect(() => {
    render(<RecipientEditorComponent data={[makeFormWithRecord(1)]} />);
  }).toThrow();
});
