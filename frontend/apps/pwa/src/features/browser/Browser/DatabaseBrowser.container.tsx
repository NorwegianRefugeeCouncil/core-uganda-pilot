import React, { FC } from 'react';
import { useMatch } from 'react-router-dom';

import { BrowserContainer } from './Browser.container';

export const DatabaseBrowserContainer: FC = () => {
  const match = useMatch('/browse/databases/:databaseId');
  const databaseId = match?.params.databaseId;

  if (!databaseId) {
    return <></>;
  }

  return <BrowserContainer databaseId={databaseId} folderId={undefined} />;
};
