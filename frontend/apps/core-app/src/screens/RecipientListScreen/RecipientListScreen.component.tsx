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
      <Box bg="secondary.100" width="100%" my="2xs" alignItems="center">
        <Box maxWidth="1180px" width="100%">
          <Box mr="auto" mt="sm" mb="xl" maxWidth="580px" width="100%">
            <RecipientListTableFilter filter={filter} setFilter={setFilter} />
          </Box>
        </Box>
      </Box>

      <Box maxWidth="1180px" mx="auto" width="100%">
        {!isLoading &&
          forms?.map((form) => (
            <RecipientListTable key={form.id} form={form} filter={filter} />
          ))}
        {isLoading && <Skeleton h="xl" />}
      </Box>
    </ScrollView>
  );
};
