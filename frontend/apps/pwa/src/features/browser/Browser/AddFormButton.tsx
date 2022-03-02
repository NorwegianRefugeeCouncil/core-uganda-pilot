import { Link } from 'react-router-dom';
import React from 'react';

import { DatabaseOrFolderId } from './types';

function addFormURL(databaseId: string, folderId: string | undefined) {
  let addFormURLString = `/edit/forms?databaseId=${databaseId}`;
  if (folderId) {
    addFormURLString += `&folderId=${folderId}`;
  }
  return addFormURLString;
}

export const AddFormButton: React.FC<DatabaseOrFolderId> = ({
  databaseId,
  folderId,
}) => {
  return (
    <Link
      className="btn btn-primary me-2"
      to={addFormURL(databaseId, folderId)}
    >
      Create a Form
    </Link>
  );
};
