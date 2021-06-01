import * as React from 'react';
import { Link } from 'react-router-dom';

export const HomeComponent: React.FC = props => {
  return (
    <div className={'d-flex p-3 border-bottom'}>
      <div className={'me-3'}>
        <i className={'bi bi-pencil'} />
      </div>
      <Link to={'/'}>Uganda Intake Form</Link>
    </div>
  );
};
