import { FormDefinition, FormType } from 'core-api-client';
import { FormProvider, useForm } from 'react-hook-form';

import { render } from '../../../testUtils/render';
import { withFormContext } from '../../../testUtils/withFormContext';
import { RecordEditorContainer } from '../RecordEditor.container';

jest.mock('../RecordEditor.component', () => {
  const { View, Text } = jest.requireActual('react-native');
  return {
    RecordEditorComponent: ({ form }: { form: FormDefinition }) => (
      <View>
        <Text>{JSON.stringify(form)}</Text>
      </View>
    ),
  };
});

it('should throw an error without a form context', () => {
  const form: FormDefinition = {
    id: 'form-1',
    code: '',
    name: 'form 1',
    databaseId: 'database-id-1',
    folderId: 'folder-id-1',
    formType: FormType.DefaultFormType,
    fields: [],
  };

  expect(() => render(<RecordEditorContainer form={form} />)).toThrowError();
});

it('should not throw an error with a form context', () => {
  const form: FormDefinition = {
    id: 'form-1',
    code: '',
    name: 'form 1',
    databaseId: 'database-id-1',
    folderId: 'folder-id-1',
    formType: FormType.DefaultFormType,
    fields: [],
  };

  const RecordEditorContainerWithContext = withFormContext(
    RecordEditorContainer,
  );

  expect(() =>
    render(<RecordEditorContainerWithContext form={form} />),
  ).not.toThrowError();
});
