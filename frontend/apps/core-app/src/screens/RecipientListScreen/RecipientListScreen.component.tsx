import * as React from 'react';
import { Skeleton, Text } from 'native-base';
import { FormWithRecord } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';

import { RecipientListTableContext } from '../../components/RecipientListTable/RecipientListTableContext';
import { RecipientListTable } from '../../components/RecipientListTable';
import { RecipientListTableFilter } from '../../components/RecipientListTable/RecipientListTableFilter';

import * as Styles from './RecipientListScreen.styles';

type Props = {
  onItemClick: (id: string) => void;
  data: FormWithRecord<Recipient>[][] | null;
  isLoading: boolean;
  error?: string;
};

export const RecipientListScreenComponent: React.FC<Props> = ({
  data,
  onItemClick,
  isLoading,
  error,
}) => {
  const tableContext = React.useContext(RecipientListTableContext);

  return (
    <Styles.Container>
      {isLoading && <Skeleton h="20" p="4" />}
      {error && (
        <Text variant="heading" color="signalDanger">
          {error}
        </Text>
      )}

      {tableContext?.tableInstance && (
        <RecipientListTableFilter
          table={tableContext.tableInstance}
          globalFilter={tableContext.tableInstance.globalFilter}
          setGlobalFilter={tableContext.tableInstance.setGlobalFilter}
        />
      )}
      {data && <RecipientListTable data={data} onItemClick={onItemClick} />}
    </Styles.Container>
  );
};
