import { FC } from 'react';
import { NavLink, useRouteMatch } from 'react-router-dom';

import { Organization } from '../../types/types';

type Props = {
  organization: Organization;
};

export const OrganizationSideBar: FC<Props> = (props) => {
  const match = useRouteMatch();
  return (
    <div className="list-group list-group-darkula" style={{ width: '15rem' }}>
      <NavLink activeClassName="active" exact className="list-group-item list-group-item-action" to={`${match.url}`}>
        Overview
      </NavLink>

      <NavLink activeClassName="active" className="list-group-item list-group-item-action" to={`${match.url}/identity-providers`}>
        Identity Providers
      </NavLink>
    </div>
  );
};
