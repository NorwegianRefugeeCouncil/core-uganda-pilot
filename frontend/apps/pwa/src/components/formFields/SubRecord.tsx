import React from 'react';

type Props = { recordId: string; select: (recordId: string) => void };

export const SubRecord: React.FC<Props> = ({ recordId, select }) => {
  return (
    <a
      href="/#"
      key={recordId}
      onClick={(e) => {
        e.preventDefault();
        select(recordId);
      }}
      className="list-group-item list-group-item-action bg-dark border-secondary text-secondary"
    >
      View Record
    </a>
  );
};
