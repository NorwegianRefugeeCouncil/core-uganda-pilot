import React, { FC, useEffect } from 'react';

import { useAppDispatch, useAppSelector } from '../../../app/hooks';
import {
  databaseGlobalSelectors,
  fetchDatabases,
} from '../../../reducers/database';
import { fetchFolders } from '../../../reducers/folder';
import { fetchForms } from '../../../reducers/form';
import {
  selectChildFolders,
  selectChildForms,
} from '../../../reducers/browser';

import { BrowserComponent } from './Browser.component';

type FolderBrowserContainerProps = {
  databaseId: string;
  folderId?: string;
};

type MenuEntry = {
  id: string;
  label: string;
  icon?: string;
  muted?: boolean;
  url: string;
  type: 'folder' | 'form';
};

type MenuEntries = MenuEntry[];

export const BrowserContainer: FC<FolderBrowserContainerProps> = ({
  folderId,
  databaseId,
}) => {
  const dispatch = useAppDispatch();

  useEffect(() => {
    dispatch(fetchDatabases());
    dispatch(fetchFolders());
    dispatch(fetchForms());
  }, [dispatch]);

  const childFolders = useAppSelector(
    selectChildFolders(folderId || databaseId),
  );
  const childForms = useAppSelector(
    selectChildForms({ dbId: databaseId, folderId }),
  );

  const database = useAppSelector((state) => {
    if (databaseId) {
      return databaseGlobalSelectors.selectById(state, databaseId);
    }
  });

  if (!databaseId && !folderId) {
    return <></>;
  }

  if (!database) {
    return <></>;
  }

  const menuEntries: MenuEntries = [];

  for (const childFolder of childFolders) {
    menuEntries.push({
      id: childFolder.id,
      label: childFolder.name,
      muted: true,
      url: `/browse/databases/${databaseId}/folders/${childFolder.id}`,
      type: 'folder',
    });
  }

  for (const childForm of childForms) {
    menuEntries.push({
      id: childForm.id,
      label: childForm.name,
      muted: true,
      url: `/browse/forms/${childForm.id}`,
      type: 'form',
    });
  }

  return (
    <div className="flex-grow-1 bg-light">
      <div className="container">
        <div className="row">
          <div className="col">
            <BrowserComponent
              databaseId={databaseId}
              folderId={folderId}
              folders={childFolders}
              forms={childForms}
            />
          </div>
        </div>
      </div>
    </div>
  );
};
