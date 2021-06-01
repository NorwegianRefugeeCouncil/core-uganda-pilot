import * as React from 'react';
import { Link } from 'react-router-dom';
import * as uuid from 'uuid';
import { Button } from '@core/ui-toolkit';

const renderLabel = (value: string) => {
  return <span className={'text-secondary mb-2 d-block'}>{value}</span>;
};

const renderValue = (value: string) => {
  return <span className={'text-dark fw-bold mb-2 d-block'}>{value}</span>;
};

const renderLink = (label: string, to: string = '#') => {
  return <Link to={to} className={'mb-2 d-block fw-bold'}>{label}</Link>;
};

export const IdentificationDocumentComponent: React.FC = props => {

  return <div className={'p-2'}>
    <div className={'px-2 my-4'}>
      <h6 className={'text-uppercase text-primary mb-0'}>Identification Document</h6>
      <h1 className={'display-5 fw-bold text-primary'}>John Doe</h1>
      <div className="input-group mb-3">
        <input type="text" className="form-control form-control-sm form-control-plaintext font-monospace bg-light text-secondary px-2" disabled value={uuid.v4()}/>
      </div>

    </div>

    <div className='row gx-0 border-bottom border-top'>
      <div className='col-12 col-lg-6 mt-lg-0'>
        <div className={'px-2 py-4'}>

          {renderLabel('Beneficiary')}
          {renderLink('John Doe', `/beneficiaries/${uuid.v4()}`)}

          {renderLabel('Type')}
          {renderValue('Passport')}

          {renderLabel('Code')}
          {renderValue('AB0293203ABA')}

          <Button size={'sm'} className={'text-white mt-4 me-2'}>Edit</Button>
          <Button size={'sm'} className={'text-white mt-4 me-2'}>Delete</Button>

        </div>
      </div>
    </div>

  </div>;

};
