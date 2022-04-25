import * as React from 'react';

import { RecipientListTableContext } from '../../components/RecipientListTable/RecipientListTableContext';
import { data } from '../../components/RecipientListTable/data.tmp';
import { RecipientListTable } from '../../components/RecipientListTable';
import { RecipientListTableFilter } from '../../components/RecipientListTable/RecipientListTableFilter';

import * as Styles from './RecipientListScreen.styles';

type Props = {
  onItemClick: (id: string) => void;
};

export const RecipientListScreenComponent: React.FC<Props> = ({
  onItemClick,
}) => {
  const tableContext = React.useContext(RecipientListTableContext);

  return (
    <Styles.Container>
      {tableContext?.tableInstance && (
        <RecipientListTableFilter
          table={tableContext.tableInstance}
          globalFilter={tableContext.tableInstance.globalFilter}
          setGlobalFilter={tableContext.tableInstance.setGlobalFilter}
        />
      )}
      <RecipientListTable data={data} onItemClick={onItemClick} />
    </Styles.Container>
  );
};
