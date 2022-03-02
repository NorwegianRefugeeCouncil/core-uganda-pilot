import React from 'react';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-icons/font/bootstrap-icons.css';
import 'bootstrap/dist/js/bootstrap.bundle.min';
import { BrowserRouter, Route, Routes } from 'react-router-dom';

import AuthenticatedApp from './AuthenticatedApp';
import { Organizations } from './organizations/Organizations';
import { OrganizationEditor } from './organizations/OrganizationEditor';
import { OrganizationPortal } from './organizations/OrganizationPortal';
import { Clients } from './clients/Clients';
import { ClientEditor } from './clients/ClientEditor';
import { IdentityProviderEditor } from './organizations/identityproviders/IdentityProviderEditor';
import { IdentityProviders } from './organizations/identityproviders/IdentityProviders';
import { OrganizationOverview } from './organizations/OrganizationOverview';

type Props = {
  baseUrl: string;
};

export const Router: React.FC<Props> = ({ baseUrl }) => (
  <BrowserRouter basename={baseUrl}>
    <Routes>
      <Route path="/" element={<AuthenticatedApp />}>
        <Route path="organizations" element={<Organizations />} />
        <Route path="organizations/add" element={<OrganizationEditor />} />
        <Route
          path="organizations/:organizationId"
          element={<OrganizationPortal />}
        >
          <Route path="identity-providers" element={<IdentityProviders />} />
          <Route
            path="identity-providers/add"
            element={<IdentityProviderEditor />}
          />
          <Route
            path="identity-providers/:id"
            element={<IdentityProviderEditor />}
          />
          <Route index element={<OrganizationOverview />} />
        </Route>
        <Route path="clients" element={<Clients />} />
        <Route path="clients/add" element={<ClientEditor />} />
        <Route path="clients/:clientId" element={<ClientEditor />} />
      </Route>
    </Routes>
  </BrowserRouter>
);
