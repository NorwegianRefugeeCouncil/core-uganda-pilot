import React, { FC } from 'react';
import { Database } from 'core-api-client';
import { Link } from 'react-router-dom';

export const DatabaseRow: FC<{ database: Database }> = ({ database }) => {
  return (
    <Link
      to={`/browse/databases/${database.id}`}
      key={database.id}
      className="list-group-item list-group-item-action py-4 fw-bold"
    >
      <i className="bi bi-box-seam me-3" />
      <span>{database.name}</span>
    </Link>
  );
};
