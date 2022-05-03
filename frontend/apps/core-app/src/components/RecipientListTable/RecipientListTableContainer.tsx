import React, { useEffect } from 'react';
import { FormDefinition } from 'core-api-client';
import { useGlobalFilter, useSortBy, useTable } from 'react-table';
import { useNavigation } from '@react-navigation/native';

import { useAPICall } from '../../hooks/useAPICall';
import { formsClient } from '../../clients/formsClient';
import { routes } from '../../constants/routes';

import { RecipientListTableComponent } from './RecipientListTableComponent';
import { createTableColumns } from './createTableColumns';
import { mapRecordsToRecordTableData } from './mapRecordsToRecordTableData';
import { RecipientListTableEntry, SortedFilteredTable } from './types';

type Props = {
  form: FormDefinition;
  // onItemClick: (id: string) => void;
  filter: string;
};

export const RecipientListTableContainer: React.FC<Props> = ({
  form,
  // onItemClick,
  filter,
}) => {
  const navigation = useNavigation();
  const [_, recipientState] = useAPICall(
    formsClient.Recipient.list,
    [
      {
        formId: form.id,
        databaseId: form.databaseId,
      },
    ],
    true,
  );

  const memoizedData = React.useMemo(
    () => mapRecordsToRecordTableData(recipientState.data || []),
    [JSON.stringify(recipientState.data)],
  );
  const memoizedColumns = React.useMemo(
    () => createTableColumns(form),
    [JSON.stringify(form)],
  );

  const table: SortedFilteredTable<RecipientListTableEntry> = useTable(
    {
      data: memoizedData,
      columns: memoizedColumns,
    },
    useGlobalFilter,
    useSortBy,
  );

  useEffect(() => {
    table.setGlobalFilter(filter);
  }, [filter]);

  const handleItemClick = (id: string) => {
    navigation.navigate(routes.recipientsProfile.name, {
      recordId: id,
      formId: form.id,
      databaseId: form.databaseId,
    });
  };

  return (
    <RecipientListTableComponent
      onItemClick={handleItemClick}
      title={form.name}
      table={table}
      error={recipientState.error}
      loading={recipientState.loading}
    />
  );
};
