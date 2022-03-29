import { FC } from 'react';
import { NavLink } from 'react-router-dom';

export const OrganizationSideBar: FC = () => {
  return (
    <div className="list-group list-group-darkula" style={{ width: '15rem' }}>
      <NavLink className="list-group-item list-group-item-action" to="">
        Overview
      </NavLink>

      <NavLink
        className="list-group-item list-group-item-action"
        to="identity-providers"
      >
        Identity Providers
      </NavLink>
    </div>
  );
};
