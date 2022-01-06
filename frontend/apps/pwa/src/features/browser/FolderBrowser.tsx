import React, { FC, Fragment, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Folder, FormDefinition } from 'core-api-client';

import { useAppDispatch, useAppSelector } from '../../app/hooks';
import {
  databaseGlobalSelectors,
  fetchDatabases,
} from '../../reducers/database';
import { fetchFolders, folderGlobalSelectors } from '../../reducers/folder';
import { fetchForms } from '../../reducers/form';
import { selectChildFolders, selectChildForms } from '../../reducers/browser';

import { FormRow } from './FormRow';
import { FolderRow } from './FolderRow';

export type FolderBrowserProps = {
  databaseId: string;
  folderId: string | undefined;
  folders: Folder[];
  forms: FormDefinition[];
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

export const FolderBrowser: FC<FolderBrowserProps> = (props) => {
  const { databaseId, folderId, folders, forms } = props;

  const formEntries = forms.map((f) => <FormRow key={f.id} form={f} />);
  const folderEntries = folders.map((f) => <FolderRow key={f.id} folder={f} />);

  const isEmpty = forms.length === 0 && folders.length === 0;
  const isEmptyDatabase = folderId === undefined && isEmpty;
  const isEmptyFolder = folderId !== undefined && isEmpty;

  return (
    <>
      <div className="py-3">
        {isEmptyDatabase
          ? emptyDatabase(databaseId, folderId)
          : addButtons(databaseId, folderId)}
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

function addButtons(databaseId: string, folderId: string | undefined) {
  return (
    <>
      {addFormButton(databaseId, folderId)}
      {addFolderButton(databaseId, folderId)}
    </>
  );
}

function addFormButton(databaseId: string, folderId: string | undefined) {
  return (
    <Link
      className="btn btn-primary me-2"
      to={addFormURL(databaseId, folderId)}
    >
      Create a Form
    </Link>
  );
}

function addFormURL(databaseId: string, folderId: string | undefined) {
  let addFormURL = `/edit/forms?databaseId=${databaseId}`;
  if (folderId) {
    addFormURL += `&folderId=${folderId}`;
  }
  return addFormURL;
}

function addFolderButton(databaseId: string, folderId: string | undefined) {
  return (
    <Link className="btn btn-primary" to={addFolderURL(databaseId, folderId)}>
      Create a Folder
    </Link>
  );
}

function addFolderURL(databaseId: string, folderId: string | undefined) {
  let addFolderURL = `/add/folders?databaseId=${databaseId}`;
  if (folderId) {
    addFolderURL += `&parentId=${folderId}`;
  }
  return addFolderURL;
}

function emptyDatabase(databaseId: string, folderId: string | undefined) {
  return (
    <div className="jumbotron">
      <h1 className="display-4">Welcome to your database!</h1>
      <p className="lead">
        Your database is empty right now. Start by adding a form.
      </p>
      <hr className="my-4" />
      <p>Design a form to start collecting data.</p>
      <p className="lead">
        {addFormButton(databaseId, folderId)}
        {addFolderButton(databaseId, folderId)}
      </p>
    </div>
  );
}

export type FolderBrowserContainerProps = {
  databaseId?: string;
  folderId?: string;
};

export const FolderBrowserContainer: FC<FolderBrowserContainerProps> = (
  props,
) => {
  const dispatch = useAppDispatch();

  useEffect(() => {
    dispatch(fetchDatabases());
    dispatch(fetchFolders());
    dispatch(fetchForms());
  }, [dispatch]);

  const childFolders = useAppSelector(
    selectChildFolders(props.folderId ? props.folderId : props.databaseId),
  );
  const childForms = useAppSelector(
    selectChildForms({ dbId: props.databaseId, folderId: props.folderId }),
  );

  const databaseId = useAppSelector((state) => {
    if (props.folderId) {
      const folder = folderGlobalSelectors.selectById(state, props.folderId);
      if (folder) {
        return folder.databaseId;
      }
    } else if (props.databaseId) {
      return props.databaseId;
    }
  });

  const database = useAppSelector((state) => {
    if (databaseId) {
      return databaseGlobalSelectors.selectById(state, databaseId);
    }
  });

  if (!props.databaseId && !props.folderId) {
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
      url: `/browse/folders/${childFolder.id}`,
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
            <FolderBrowser
              databaseId={database.id}
              folderId={props.folderId}
              folders={childFolders}
              forms={childForms}
            />
          </div>
        </div>
      </div>
    </div>
  );
};
