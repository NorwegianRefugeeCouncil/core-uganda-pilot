import { FormDefinition, FormType } from 'core-api-client';

import { render } from '../../../testUtils/render';
import { RecordEditorComponent } from '../RecordEditor.component';

jest.mock('../FieldInput', () => {
  const { View, Text } = jest.requireActual('react-native');
  return {
    FieldInput: ({
      form,
      field,
    }: {
      form: FormDefinition;
      field: FormDefinition['fields'][number];
    }) => (
      <View>
        <Text>{JSON.stringify(form)}</Text>
        <Text>{JSON.stringify(field)}</Text>
      </View>
    ),
  };
});

it('should match the snapshot', () => {
  const form: FormDefinition = {
    id: 'form-1',
    code: '',
    name: 'form 1',
    databaseId: 'database-id-1',
    folderId: 'folder-id-1',
    formType: FormType.DefaultFormType,
    fields: [
      {
        id: 'field-1',
        name: 'field 1',
        description: 'description 1',
        code: '',
        required: false,
        key: false,
        fieldType: {
          text: {},
        },
      },
      {
        id: 'field-2',
        name: 'field 2',
        description: 'description 2',
        code: '',
        required: true,
        key: true,
        fieldType: {
          multilineText: {},
        },
      },
    ],
  };

  const { toJSON } = render(<RecordEditorComponent form={form} />);
  expect(toJSON()).toMatchSnapshot();
});
