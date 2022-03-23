import * as React from 'react';
import { FormType } from 'core-api-client';

import Table, { GlobalFilter } from '../../components/Table';
import { TableContext } from '../../components/Table/useTableContext';

import * as Styles from './RecipientListScreen.styles';

type Props = {
  handleItemClick: (id: string) => void;
};

export const RecipientListScreenComponent: React.FC<Props> = ({
  handleItemClick,
}) => {
  const tableContext = React.useContext(TableContext);

  return (
    <Styles.Container>
      {tableContext?.tableInstance && (
        <GlobalFilter table={tableContext.tableInstance} />
      )}
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
    </Styles.Container>
  );
};
