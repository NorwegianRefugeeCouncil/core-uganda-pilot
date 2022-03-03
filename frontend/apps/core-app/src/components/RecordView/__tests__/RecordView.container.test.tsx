import {
  FieldDefinition,
  FormType,
  FieldValue,
  FieldKind,
} from 'core-api-client';

import { render } from '../../../testUtils/render';
import * as nfv from '../normaliseFieldValues';
import { RecordViewContainer } from '../RecordView.container';
import { NormalisedFieldValue } from '../RecordView.types';

jest.mock('../RecordView.component', () => {
  const { View, Text } = jest.requireActual('react-native');
  return {
    RecordViewComponent: ({
      fieldValues,
    }: {
      fieldValues: NormalisedFieldValue[];
    }) => (
      <View>
        <Text testID="fieldValues">{JSON.stringify(fieldValues)}</Text>
      </View>
    ),
  };
});

const makeForm = (fields: FieldDefinition[]) => ({
  id: 'form-id',
  code: 'form-code',
  databaseId: 'database-id',
  folderId: 'folder-id',
  name: 'form-name',
  formType: FormType.DefaultFormType,
  fields,
});

const makeRecord = (values: FieldValue[]) => ({
  id: 'record-id',
  databaseId: 'database-id',
  formId: 'form-id',
  ownerId: undefined,
  values,
});

it('should', () => {
  const fieldValues: NormalisedFieldValue[] = [
    {
      label: 'label-1',
      fieldType: FieldKind.Text,
      value: 'test-1',
      formattedValue: 'format-test-1',
    },
  ];
  const normaliseFieldValuesSpy = jest
    .spyOn(nfv, 'normaliseFieldValues')
    .mockImplementation(() => fieldValues);

  const form = makeForm([]);
  const record = makeRecord([]);

  const { getByTestId } = render(
    <RecordViewContainer form={form} record={record} />,
  );

  expect(normaliseFieldValuesSpy).toHaveBeenCalledWith(form, record);

  expect(getByTestId('fieldValues')).toHaveTextContent(
    JSON.stringify(fieldValues),
  );
});
