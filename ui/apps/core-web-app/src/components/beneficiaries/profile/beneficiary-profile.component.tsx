import * as React from 'react';
import { CardTitle } from '@core/ui-toolkit';
import { Link } from 'react-router-dom';
import * as uuid from 'uuid';

export type Props = {
  beneficiaryId: string
}

const renderLabel = (value: string) => {
  return <span className={'text-secondary mb-2 d-block'}>{value}</span>;
};

const renderValue = (value: string) => {
  return <span className={'text-dark fw-bold mb-2 d-block'}>{value}</span>;
};

const renderLink = (label: string, to: string = '#') => {
  return <Link to={to} className={'mb-2 d-block fw-bold'}>{label}</Link>;
};


export const BeneficiaryProfileComponent: React.FC<Props> = props => {

  const globalDetails = [
    {
      label: 'Name',
      value: 'John Doe'
    },
    {
      label: 'Status',
      value: 'Refugee'
    },
    {
      label: 'Gender',
      value: 'Male'
    },
    {
      label: 'ID Number',
      value: '546456'
    }
  ];

  const availableForms = [
    {
      label: 'ICLA Intake'
    },
    {
      label: 'ICLA to Shelter referral'
    },
    {
      label: 'Protection Intake Form'
    },
    {
      label: 'Shelter Case Opening'
    }
  ];

  const identificationDocuments = [
    {
      label: 'Passport'
    }, {
      label: 'Driving License'
    }
  ];

  return <div className={'p-2'}>
    <div className={'px-2 my-4'}>
      <h6 className={'text-uppercase text-primary mb-0'}>Beneficiary</h6>
      <h1 className={'display-5 fw-bold text-primary'}>John Doe</h1>
      <div className="input-group mb-3">
        <input type="text" className="form-control form-control-sm form-control-plaintext font-monospace bg-light text-secondary px-2" disabled value={uuid.v4()}/>
      </div>

    </div>

    <div className='row gx-0 border-bottom'>
      <div className='col-12 col-lg-6 mt-lg-0'>
        <div className={'py-4 px-2 px-lg-3 border-top'}>
          <CardTitle className={'text-uppercase text-primary fw-bold fs-6 mb-3'}>Global Details</CardTitle>
          {globalDetails.map(item => {
            return <div className={'row'}>
              <div className={'col-12 col-sm-3'}>
                {renderLabel(item.label)}
              </div>
              <div className={'col-12 col-sm-3'}>
                {renderValue(item.value)}
              </div>
            </div>;
          })}
        </div>
      </div>
      <div className='col-12 col-lg-6 mt-lg-0'>
        <div className={'py-4 px-2 px-lg-3 border-top border-left-0 border-lg-left'}>
          <CardTitle className={'text-uppercase text-primary fw-bold fs-6 mb-3'}>Forms</CardTitle>
          {availableForms.map(item => {
            return <div className={'row'}>
              <div className={'col-12'}>
                {renderLink(item.label)}
              </div>
            </div>;
          })}
        </div>
      </div>
      <div className='col-12 col-lg-6 mt-lg-0'>
        <div className={'py-4 px-2 px-lg-3 border-top border-left-0'}>
          <CardTitle className={'text-uppercase text-primary fw-bold fs-6 mb-3'}>Household</CardTitle>
          <div className={'row'}>
            <div className={'col-6'}>
              {renderLink('Household', `/households/${uuid.v4()}`)}
            </div>
          </div>
        </div>
      </div>
      <div className='col-12 col-lg-6 mt-lg-0'>
        <div className={'py-4 px-2 px-lg-3 border-top border-left-0 border-lg-left'}>
          <CardTitle className={'text-uppercase text-primary fw-bold fs-6 mb-3'}>Identification Documents</CardTitle>
          {identificationDocuments.map(item => {
            return <div className={'row'}>
              <div className={'col-9'}>
                {renderLink(item.label, `/identificationdocuments/${uuid.v4()}`)}
              </div>
            </div>;
          })}
          <div className={'row mt-3'}>
            <div className={'col-12'}>
              <a href={'#'}>Add</a>
            </div>
          </div>
        </div>
      </div>
    </div>

  </div>;


};
