import * as React from 'react';
import { Box, Skeleton, Text } from 'native-base';
import { FormDefinition, FormWithRecord } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';

import { RecipientListTableContext } from '../../components/RecipientListTable/RecipientListTableContext';
import { RecipientListTable } from '../../components/RecipientListTable';
import { RecipientListTableFilter } from '../../components/RecipientListTable/RecipientListTableFilter';

import * as Styles from './RecipientListScreen.styles';

type Props = {
  onItemClick: (id: string) => void;
  data: FormWithRecord<Recipient>[][] | null;
  form: FormDefinition | undefined;
  isLoading: boolean;
  error?: string;
};

export const RecipientListScreenComponent: React.FC<Props> = ({
  data,
  form,
  onItemClick,
  isLoading,
  error,
}) => {
  const tableContext = React.useContext(RecipientListTableContext);

  return (
    <Styles.Container bg="white" height="100%">
      {isLoading && <Skeleton h="20" p="4" />}
      <Box bg="secondary.100" width="100%" my="16px">
        <Box mx="130px" mt="26px" mb="42px" maxWidth="580px">
          {tableContext?.tableInstance && (
            <RecipientListTableFilter
              table={tableContext.tableInstance}
              globalFilter={tableContext.tableInstance.globalFilter}
              setGlobalFilter={tableContext.tableInstance.setGlobalFilter}
            />
          )}
        </Box>
      </Box>
      <Box px="130px">
        {data && form && (
          <Box>
            <Text variant="heading" level="2">
              {form.name}
            </Text>
            <RecipientListTable
              data={data}
              form={form}
              onItemClick={onItemClick}
            />
          </Box>
        )}
        {error && (
          <Text variant="heading" color="signalDanger">
            {error}
          </Text>
        )}
      </Box>
    </Styles.Container>
  );
};
