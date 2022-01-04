import { Folder, FormDefinition } from 'core-api-client';

import { RootState } from '../app/store';

import { selectByFolderOrDBId } from './form';
import { selectByParentId, selectDatabaseRootFolders } from './folder';
import { databaseGlobalSelectors } from './database';

export const selectChildFolders =
  (dbOrFolderId?: string) =>
  (state: RootState): Folder[] => {
    if (!dbOrFolderId) {
      return [];
    }
    const db = databaseGlobalSelectors.selectById(state, dbOrFolderId);
    if (db) {
      return selectDatabaseRootFolders(state, db.id);
    }
    return selectByParentId(state, dbOrFolderId);
  };

export const selectChildForms =
  (ids: { dbId?: string; folderId?: string }) =>
  (state: RootState): FormDefinition[] => {
    return selectByFolderOrDBId(state, ids);
  };
