import * as React from 'react';
import { Card, CardBody, CardTitle } from '@core/ui-toolkit';

export type Props = {
  beneficiaryId: string
}

const renderLabel = (value: string) => {
  return <span className={'text-secondary mb-2 d-block'}>{value}</span>;
};

const renderValue = (value: string) => {
  return <span className={'text-dark fw-bold mb-2 d-block'}>{value}</span>;
};

const renderLink = (label: string) => {
  return <a href='#' className={'mb-2 d-block fw-bold'}>{label}</a>;
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

  return <div className={'p-2 overflow-hidden'}>
    <div className={'px-2 my-4'}>
      <h6 className={'text-uppercase text-primary mb-0'}>Beneficiary</h6>
      <h1 className={'display-5 fw-bold text-primary'}>John Doe</h1>
    </div>

    <div className='row border-0 border-lg-top border-bottom'>
      <div className='col-12 col-lg-6'>
        <div className={'p-4 border-top border-lg-top-0'}>
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
        <div className={'p-4 border-top border-lg-top-0 border-left-0 border-lg-left'}>
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
    </div>

  </div>;


};
