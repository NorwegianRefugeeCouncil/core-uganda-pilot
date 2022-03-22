import * as React from 'react';
import { Text, VStack } from 'native-base';
import { RouteProp } from '@react-navigation/native';
import { FormType } from 'core-api-client';

import { RootParamList } from '../../navigation/types';
import { Table } from '../../components/Table';
import { DefaultColumnFilter } from '../../components/Table/SearchInput';

import * as Styles from './RecipientListScreen.styles';

type Props = {
  route: RouteProp<RootParamList, 'RecipientList'>;
  handleItemClick: (id: string) => void;
};

export const RecipientListScreenComponent: React.FC<Props> = ({
  route,
  handleItemClick,
}) => {

  return (
    <Styles.Container>
      <Text variant="display">{route.name}</Text>
      <DefaultColumnFilter />
      <VStack space={2} width="sm">
        <Table
          records={[
            {
              id: 'id1',
              formId: 'formId',
              ownerId: undefined,
              databaseId: 'dbid',
              values: [
                {
                  value: 'value1',
                  fieldId: 'field1',
                },
                { value: 'value2', fieldId: 'field2' },
              ],
            },
            {
              id: 'id2',
              formId: 'formId',
              ownerId: undefined,
              databaseId: 'dbid',
              values: [
                {
                  value: 'value3',
                  fieldId: 'field1',
                },
                { value: 'value4', fieldId: 'field2' },
              ],
            },
          ]}
          form={{
            id: 'form1',
            name: 'name',
            databaseId: 'dbid',
            formType: FormType.RecipientFormType,
            folderId: '',
            fields: [
              {
                id: 'field1',
                name: 'fieldName',
                fieldType: { text: {} },
                key: false,
                code: '',
                required: false,
                description: '',
              },
              {
                id: 'field2',
                name: 'fieldName2',
                fieldType: { text: {} },
                key: false,
                code: '',
                required: false,
                description: '',
              },
            ],
            code: '',
          }}
          handleItemClick={handleItemClick}
        />
      </VStack>
    </Styles.Container>
  );
};
