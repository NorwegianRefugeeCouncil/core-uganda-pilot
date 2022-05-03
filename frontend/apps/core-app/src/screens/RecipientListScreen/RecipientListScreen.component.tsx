import * as React from 'react';
import { Box, ScrollView, Skeleton } from 'native-base';
import { FormDefinition } from 'core-api-client';

import { RecipientListTable } from '../../components/RecipientListTable';
import { RecipientListTableFilter } from '../../components/RecipientListTable/RecipientListTableFilter';

type Props = {
  forms: FormDefinition[] | null;
  isLoading: boolean;
};

export const RecipientListScreenComponent: React.FC<Props> = ({
  forms,
  isLoading,
}) => {
  const [filter, setFilter] = React.useState('');

  return (
    <ScrollView bg="white">
      <Box bg="secondary.100" width="100%" my="16px" alignItems="center">
        <Box maxWidth={1180} width="100%">
          <Box mr="auto" mt="26px" mb="42px" maxWidth={580} width="100%">
            <RecipientListTableFilter
              globalFilter={filter}
              setGlobalFilter={setFilter}
            />
          </Box>
        </Box>
      </Box>

      <Box maxWidth={1180} mx="auto">
        {forms?.length &&
          forms.map((form) => (
            <RecipientListTable key={form.id} form={form} filter={filter} />
          ))}

        {isLoading && <Skeleton h="40" p="4" />}
      </Box>
    </ScrollView>
  );
};
