import React, { FC } from 'react';
import { useMatch } from 'react-router-dom';

import { BrowserContainer } from './Browser.container';

export const FolderBrowserContainer: FC = () => {
  const match = useMatch('/browse/databases/:databaseId/folders/:folderId');
  const folderId = match?.params.folderId;
  const databaseId = match?.params.databaseId;

  if (!databaseId || !folderId) {
    return <></>;
  }

  return <BrowserContainer databaseId={databaseId} folderId={folderId} />;
};
