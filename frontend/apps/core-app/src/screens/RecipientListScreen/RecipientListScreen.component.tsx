import * as React from 'react';
import { Badge, Box, ScrollView, Skeleton, Text } from 'native-base';
import { FormDefinition, FormWithRecord } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';

import { RecipientListTableContext } from '../../components/RecipientListTable/RecipientListTableContext';
import { RecipientListTable } from '../../components/RecipientListTable';
import { RecipientListTableFilter } from '../../components/RecipientListTable/RecipientListTableFilter';

type Props = {
  onItemClick: (id: string) => void;
  data: FormWithRecord<Recipient>[][] | null;
  forms: FormDefinition[] | null;
  isLoading: boolean;
  error?: string;
};

export const RecipientListScreenComponent: React.FC<Props> = ({
  data,
  forms,
  onItemClick,
  isLoading,
  error,
}) => {
  const tableContext = React.useContext(RecipientListTableContext);

  return (
    <ScrollView bg="white">
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
      <Box maxWidth="1180px">
        {data &&
          forms &&
          forms.map((form) => (
            <Box key={form.id} mb="63px">
              <Box flexDirection="row" alignItems="center">
                <Text variant="heading" level="2">
                  {form.name}
                </Text>
                <Badge
                  bg="secondary.500"
                  height="5"
                  width="7"
                  borderRadius="4px"
                  mx="12px"
                >
                  <Text variant="heading" level="5" color="white">
                    {data.length}
                  </Text>
                </Badge>
              </Box>
              <RecipientListTable
                data={data}
                form={form}
                onItemClick={onItemClick}
              />
            </Box>
          ))}
        {error && (
          <Text variant="heading" color="signalDanger">
            {error}
          </Text>
        )}
      </Box>
    </ScrollView>
  );
};
