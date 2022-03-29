import React, { FC } from 'react';
import { Folder, FormDefinition } from 'core-api-client';

import { FormRow } from '../FormRow';
import { FolderRow } from '../FolderRow';

import { AddFormButton } from './AddFormButton';
import { AddFolderButton } from './AddFolderButton';
import { EmptyState } from './EmptyState';

export type FolderBrowserProps = {
  databaseId: string;
  folderId: string | undefined;
  folders: Folder[];
  forms: FormDefinition[];
};

export const BrowserComponent: FC<FolderBrowserProps> = (props) => {
  const { databaseId, folderId, folders, forms } = props;

  const formEntries = forms.map((f) => <FormRow key={f.id} form={f} />);
  const folderEntries = folders.map((f) => <FolderRow key={f.id} folder={f} />);

  const isEmpty = forms.length === 0 && folders.length === 0;
  const isEmptyDatabase = folderId === undefined && isEmpty;
  const isEmptyFolder = folderId !== undefined && isEmpty;

  return (
    <>
      <div className="py-3">
        {isEmptyDatabase ? (
          <EmptyState databaseId={databaseId} folderId={folderId} />
        ) : (
          <>
            <AddFormButton databaseId={databaseId} folderId={folderId} />
            <AddFolderButton databaseId={databaseId} folderId={folderId} />
          </>
        )}
      </div>

      <div className="list-group shadow">
        {formEntries}
        {folderEntries}
        {isEmptyFolder ? (
          <div className="list-group-item py-4">This folder is empty</div>
        ) : (
          <></>
        )}
      </div>
    </>
  );
};
