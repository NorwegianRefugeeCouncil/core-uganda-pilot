import React, { useEffect } from 'react';
import { Outlet } from 'react-router-dom';

import { useAppDispatch } from '../app/hooks';
import { fetchDatabases } from '../reducers/database';
import { fetchForms } from '../reducers/form';
import { fetchFolders } from '../reducers/folder';
import { NavBarContainer } from '../features/navbar/navbar';

const AuthenticatedApp: React.FC = () => {
  const dispatch = useAppDispatch();

  useEffect(() => {
    dispatch(fetchDatabases());
    dispatch(fetchForms());
    dispatch(fetchFolders());
  }, [dispatch]);
  return (
    <div className="App">
      <div
        style={{ maxHeight: '100vh', maxWidth: '100vw' }}
        className="d-flex flex-column wh-100 vh-100"
      >
        <NavBarContainer />
        <Outlet />
      </div>
    </div>
  );
};

export default AuthenticatedApp;
