import React from 'react';

import { FormValue } from '../../reducers/Recorder/types';

import { SubRecord } from './SubRecord';

type Props = {
  records: FormValue[];
  select: (id: string) => void;
};

export const SubRecordList: React.FC<Props> = ({ records, select }) => {
  return (
    <div className="list-group bg-dark mb-3">
      {records.map((r) => (
        <SubRecord key={r.id} recordId={r.id} select={select} />
      ))}
    </div>
  );
};
