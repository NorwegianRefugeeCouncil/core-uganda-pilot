import { Link } from 'react-router-dom';
import React from 'react';

import { DatabaseOrFolderId } from './types';

function addFolderURL(databaseId: string, folderId: string | undefined) {
  let addFolderURLString = `/add/folders?databaseId=${databaseId}`;
  if (folderId) {
    addFolderURLString += `&parentId=${folderId}`;
  }
  return addFolderURLString;
}

export const AddFolderButton: React.FC<DatabaseOrFolderId> = ({
  databaseId,
  folderId,
}) => {
  return (
    <Link className="btn btn-primary" to={addFolderURL(databaseId, folderId)}>
      Create a Folder
    </Link>
  );
};
