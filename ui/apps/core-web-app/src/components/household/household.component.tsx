import * as React from 'react';
import * as uuid from 'uuid';
import { Link } from 'react-router-dom';
import { Button } from '@core/ui-toolkit';

const renderLink = (label: string, to: string = '#') => {
  return <Link to={to} className={'mb-2 d-block fw-bold'}>{label}</Link>;
};


const renderLabel = (value: string) => {
  return <span className={'text-secondary mb-2 d-block'}>{value}</span>;
};

const renderValue = (value: string) => {
  return <span className={'text-dark fw-bold mb-2 d-block'}>{value}</span>;
};

export const HouseholdComponent: React.FC = props => {
  return <div className={'p-2'}>
    <div className={'px-2 my-4'}>
      <h1 className={'display-5 fw-bold text-primary'}>Household</h1>
      <div className="input-group mb-3">
        <input type="text" className="form-control form-control-sm form-control-plaintext font-monospace bg-light text-secondary px-2" disabled value={uuid.v4()}/>
      </div>

    </div>

    <div className='row gx-0 border-bottom border-top'>
      <div className='col-12 col-lg-6 mt-lg-0'>
        <div className={'px-2 py-4'}>


          {renderLabel('Head of Household')}
          {renderLink('Bobby P. Peterson', `/beneficiaries/${uuid.v4()}`)}

          {renderLabel('Members')}
          {renderLink('John Doe', `/beneficiaries/${uuid.v4()}`)}
          {renderLink('Mary Poppins', `/beneficiaries/${uuid.v4()}`)}
          {renderLink('Bo Diddley', `/beneficiaries/${uuid.v4()}`)}
          {renderLink('Steeve Ray Vaughn', `/beneficiaries/${uuid.v4()}`)}

          <Button size={'sm'} className={'text-white mt-4'}>Edit</Button>

        </div>
      </div>
    </div>

  </div>;
};
