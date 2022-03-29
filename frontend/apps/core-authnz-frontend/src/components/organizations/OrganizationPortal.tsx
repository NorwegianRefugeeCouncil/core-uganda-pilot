import { FC } from 'react';
import { Outlet, useMatch } from 'react-router-dom';

import { useOrganization } from '../../hooks/hooks';

import { OrganizationSideBar } from './OrganizationSideBar';

export const OrganizationPortal: FC = () => {
  const match = useMatch('organizations/:organizationId');
  const organization = useOrganization(match?.params.organizationId);

  if (!organization) {
    return <></>;
  }

  return (
    <div className="flex-grow-1 d-flex flex-column">
      <div className="py-2 ps-4 bg-darkula text-white">
        <h5 className="p-0 m-2">{organization.name}</h5>
      </div>
      <div className="d-flex flex-row flex-grow-1 mt-4 px-4">
        <div className="">
          <OrganizationSideBar />
        </div>
        <div className="flex-grow-1 ps-4 pe-2">
          <Outlet />
        </div>
      </div>
    </div>
  );
};
