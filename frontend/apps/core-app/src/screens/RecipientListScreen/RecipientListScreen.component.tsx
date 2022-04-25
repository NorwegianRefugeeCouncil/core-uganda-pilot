import * as React from 'react';

import { RecordTableContext } from '../../components/RecordTable/RecordTableContext';
import { data } from '../../components/RecordTable/data.tmp';
import RecordTable from '../../components/RecordTable';
import { GlobalTableFilter } from '../../components/RecordTable/GlobalTableFilter';

import * as Styles from './RecipientListScreen.styles';

type Props = {
  onItemClick: (id: string) => void;
};

export const RecipientListScreenComponent: React.FC<Props> = ({
  onItemClick,
}) => {
  const tableContext = React.useContext(RecordTableContext);

  return (
    <Styles.Container>
      {tableContext?.tableInstance && (
        <GlobalTableFilter
          table={tableContext.tableInstance}
          globalFilter={tableContext.tableInstance.globalFilter}
          setGlobalFilter={tableContext.tableInstance.setGlobalFilter}
        />
      )}
      <RecordTable data={data} onItemClick={onItemClick} />
    </Styles.Container>
  );
};
