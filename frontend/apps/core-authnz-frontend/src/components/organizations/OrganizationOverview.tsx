import { FC } from 'react';

import { Organization } from '../../types/types';

type Props = {
  organization: Organization;
};

export const OrganizationOverview: FC<Props> = (props) => {
  const { organization } = props;
  return (
    <div className="card bg-dark text-white border-secondary pb-2">
      <div className="card-body">
        <div className="form-group mb-2">
          <label className="form-label">Organization ID</label>
          <input className="form-control form-control-darkula" type="text" disabled value={organization.id} />
        </div>

        <div className="form-group">
          <label className="form-label">Organization Name</label>
          <input className="form-control form-control-darkula" type="text" value={organization.name} />
        </div>
      </div>
    </div>
  );
};
