import React from 'react';

import { AddFormButton } from './AddFormButton';
import { AddFolderButton } from './AddFolderButton';
import { DatabaseOrFolderId } from './types';

export const EmptyState: React.FC<DatabaseOrFolderId> = ({
  databaseId,
  folderId,
}) => {
  return (
    <div className="jumbotron">
      <h1 className="display-4">Welcome to your database!</h1>
      <p className="lead">
        Your database is empty right now. Start by adding a form.
      </p>
      <hr className="my-4" />
      <p>Design a form to start collecting data.</p>
      <p className="lead">
        <AddFormButton databaseId={databaseId} folderId={folderId} />
        <AddFolderButton databaseId={databaseId} folderId={folderId} />
      </p>
    </div>
  );
};
