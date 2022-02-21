import * as React from 'react';
import { Record } from 'core-api-client';

type Props = {
  record: Record;
};

export const RecordView: React.FC<Props> = ({ record }) => {
  return <>{record.id}</>;
};
