import { BrowserRouter, Route, Routes } from 'react-router-dom';
import React from 'react';

import { RecordEditor } from './features/recorder';
import { FolderEditor } from './features/folders/FolderEditor';
import { DatabaseEditor } from './features/databases/DatabaseEditor';
import { FormBrowserContainer } from './features/browser/FormBrowser';
import { DatabasesContainer } from './features/browser/Databases';
import { RecordBrowser } from './features/browser/RecordBrowser';
import { FormerContainer } from './features/former/FormerContainer';
import AuthenticatedApp from './components/AuthenticatedApp';
import { DatabaseBrowser, FolderBrowser } from './features/browser/Browser';

type Props = {
  baseUrl: string;
};

export const Router: React.FC<Props> = ({ baseUrl }) => {
  return (
    <BrowserRouter basename={baseUrl}>
      <Routes>
        <Route path="/" element={<AuthenticatedApp />}>
          <Route path="/edit/forms/:formId/record" element={<RecordEditor />} />
          <Route path="/edit/forms" element={<FormerContainer />} />
          <Route path="/add/folders" element={<FolderEditor />} />
          <Route path="/edit/databases" element={<DatabaseEditor />} />
          <Route
            path="/browse/databases/:databaseId"
            element={<DatabaseBrowser />}
          />
          <Route
            path="/browse/databases/:databaseId/folders/:folderId"
            element={<FolderBrowser />}
          />
          <Route
            path="/browse/forms/:formId"
            element={<FormBrowserContainer />}
          />
          <Route path="/browse/databases" element={<DatabasesContainer />} />
          <Route path="/browse/records/:recordId" element={<RecordBrowser />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
};
