import { Outlet } from 'react-router-dom';
import React, { FC } from 'react';

import { NavBar } from './navbar/NavBar';

const AuthenticatedApp: FC = () => {
  return (
    <div className="d-flex flex-column vh-100 vw-100 bg-dark">
      <NavBar />
      <Outlet />
    </div>
  );
};
export default AuthenticatedApp;
