import * as React from 'react';

import RecordTable, { GlobalTableFilter } from '../../components/RecordTable';
import { RecordTableContext } from '../../components/RecordTable/RecordTableContext';
import { data } from '../../components/RecordTable/data.tmp';

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
        <GlobalTableFilter table={tableContext.tableInstance} />
      )}
      <RecordTable data={data} onItemClick={onItemClick} />
    </Styles.Container>
  );
};
