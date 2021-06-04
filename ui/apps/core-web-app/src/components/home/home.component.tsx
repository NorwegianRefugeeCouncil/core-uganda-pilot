import * as React from 'react';
import { CardBody } from '@core/ui-toolkit';
import { Link } from 'react-router-dom';

export const HomeComponent: React.FC = props => {
  return (
    <div className='container pt-3'>
      <div className='row'>
        <div className='col-6'>
          <Link to='/register-beneficiary' className='card mb-3'>
            <CardBody>
              <div className='d-flex flex-row'>
                <i className='bi bi-person pe-2' />
                <span>Register new Beneficiary</span>
              </div>
            </CardBody>
          </Link>
        </div>
        <div className='col-6'>
          <Link to='/register-household' className='card mb-3'>
            <CardBody>
              <div className='d-flex flex-row'>
                <i className='bi bi-house pe-2' />
                <span>Register new Household</span>
              </div>
            </CardBody>
          </Link>
        </div>
        <div className='col-6'>
          <Link to='/uganda-assessment' className='card mb-3'>
            <CardBody>
              <div className='d-flex flex-row'>
                <i className='bi bi-alarm pe-2' />
                <span>Fill Uganda Assessment</span>
              </div>
            </CardBody>
          </Link>
        </div>
        <div className='col-6'>
          <Link to='/servicecatalog' className='card mb-3'>
            <CardBody>
              <div className='d-flex flex-row'>
                <i className='bi bi-alarm pe-2' />
                <span>Service Catalog</span>
              </div>
            </CardBody>
          </Link>
        </div>
      </div>
    </div>
  );
  // return (
  //   <div className={'d-flex p-3 border-bottom'}>
  //     <div className={'me-3'}>
  //       <i className={'bi bi-pencil'} />
  //     </div>
  //     <Link to={'/'}>Uganda Intake Form</Link>
  //   </div>
  // );
};
