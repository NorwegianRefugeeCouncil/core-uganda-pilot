import { FC } from 'react';
import { useMatch } from 'react-router-dom';

import { useOrganization } from '../../hooks/hooks';

export const OrganizationOverview: FC = () => {
  const match = useMatch('organizations/:organizationId');
  const organization = useOrganization(match?.params.organizationId);

  if (!organization) {
    return <></>;
  }

  return (
    <div className="card bg-dark text-white border-secondary pb-2">
      <div className="card-body">
        <div className="form-group mb-2">
          <label className="form-label">Organization ID</label>
          <input
            className="form-control form-control-darkula"
            type="text"
            disabled
            value={organization.id}
          />
        </div>

        <div className="form-group">
          <label className="form-label">Organization Name</label>
          <input
            className="form-control form-control-darkula"
            type="text"
            value={organization.name}
          />
        </div>
      </div>
    </div>
  );
};
