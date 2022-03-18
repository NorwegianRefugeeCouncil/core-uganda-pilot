import React from 'react';
import { Recipient } from 'core-api-client/src/types/client/Recipient';
import { FormDefinition } from 'core-api-client';

import { TableComponent } from './Table.component';

type Props = {
  records: Recipient[];
  form: FormDefinition;
  handleItemClick: (id: string) => void;
  searchTerm: string;
};

export const TableContainer: React.FC<Props> = ({
  records,
  handleItemClick,
  form,
  searchTerm,
}) => {
  return (
    <TableComponent
      records={records}
      handleItemClick={handleItemClick}
      form={form}
      searchTerm={searchTerm}
    />
  );
};
