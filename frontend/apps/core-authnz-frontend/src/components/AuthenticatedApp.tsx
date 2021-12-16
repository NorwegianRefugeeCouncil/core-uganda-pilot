import { Route, Switch } from 'react-router-dom';
import React, { FC } from 'react';

import { NavBar } from './navbar/NavBar';
import { OrganizationEditor } from './organizations/OrganizationEditor';
import { OrganizationPortal } from './organizations/OrganizationPortal';
import { Organizations } from './organizations/Organizations';
import { ClientEditor } from './clients/ClientEditor';
import { Clients } from './clients/Clients';

const AuthenticatedApp: FC = () => {
  return (
    <div className="d-flex flex-column vh-100 vw-100 bg-dark">
      <NavBar />
      <Switch>
        <Route path="/organizations/add" component={OrganizationEditor} />
        <Route path="/organizations/:organizationId" component={OrganizationPortal} />
        <Route path="/organizations" component={Organizations} />
        <Route path="/clients/add" component={ClientEditor} />
        <Route path="/clients/:clientId" component={ClientEditor} />
        <Route path="/clients" component={Clients} />
      </Switch>
    </div>
  );
};
export default AuthenticatedApp;
