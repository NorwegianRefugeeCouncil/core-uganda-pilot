import * as React from 'react';

import Table, { GlobalTableFilter } from '../../components/RecordTable';
import { RecordTableContext } from '../../components/RecordTable/RecordTableContext';
import {
  data,
  forms,
  formWithRecords,
  recordsForm0,
} from '../../components/RecordTable/data.tmp';

import * as Styles from './RecipientListScreen.styles';

type Props = {
  handleItemClick: (id: string) => void;
};

export const RecipientListScreenComponent: React.FC<Props> = ({
  handleItemClick,
}) => {
  const tableContext = React.useContext(RecordTableContext);

  return (
    <Styles.Container>
      {tableContext?.tableInstance && (
        <GlobalTableFilter table={tableContext.tableInstance} />
      )}
      <Table
        data={data}
        forms={forms}
        records={recordsForm0}
        handleItemClick={handleItemClick}
      />
    </Styles.Container>
  );
};
