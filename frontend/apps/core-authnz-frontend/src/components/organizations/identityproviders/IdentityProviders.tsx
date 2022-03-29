import { FC } from 'react';
import { Link, useMatch } from 'react-router-dom';

import { useIdentityProviders } from '../../../hooks/hooks';
import { SectionTitle } from '../../sectiontitle/SectionTitle';

export const IdentityProviders: FC = () => {
  const match = useMatch('organizations/:organizationId/identity-providers');
  const idps = useIdentityProviders(match?.params.organizationId);

  return (
    <div>
      <SectionTitle className="text-light" title="Identity Providers">
        <Link className="btn btn-success btn-sm" to="add">
          Add Identity Provider
        </Link>
      </SectionTitle>

      <div className="list-group list-group-darkula">
        {idps.map((idp) => (
          <Link
            key={idp.id}
            className="list-group-item list-group-item-action"
            to={`${match?.pathname}/${idp.id}`}
          >
            {idp.name}{' '}
            <span className="badge bg-dark font-monospace">
              {idp.emailDomain}
            </span>
          </Link>
        ))}
        {idps.length === 0 ? (
          <div className="disabled list-group-item">No Identity Provider</div>
        ) : (
          <></>
        )}
      </div>
    </div>
  );
};
