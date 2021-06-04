import * as React from 'react';
import * as uuid from 'uuid';

export const UgandaAssessmentComponent: React.FC = props => {
  return <div>
    <div className={'px-2 my-4'}>
      <h6 className={'text-uppercase text-primary mb-0'}>Service</h6>
      <h1 className={'display-5 fw-bold text-primary'}>Uganda Assessment</h1>
      <div className='input-group mb-3'>
        <input type='text'
               className='form-control form-control-sm form-control-plaintext font-monospace bg-light text-secondary px-2'
               disabled value={uuid.v4()} />
      </div>
    </div>

    <div className='px-2'>
      <div className='row'>
        <div className='col'>
          <h6 className={'text-uppercase text-primary mb-3'}>Required Information</h6>

          <div className='input-group mb-3'>
            <input type='text' className='form-control' value='beneficiary' disabled />
            <span style={{ width: '15rem' }} className='input-group-text'>Beneficiary</span>
          </div>

          <div className='input-group mb-3'>
            <input type='text' className='form-control' value='age' disabled />
            <span style={{ width: '15rem' }} className='input-group-text'>Number</span>
          </div>

          <div className='input-group mb-3'>
            <input type='text' className='form-control' value='needsLegalAssistance' disabled />
            <span style={{ width: '15rem' }} className='input-group-text'>Boolean</span>
          </div>

          <div className='input-group mb-3'>
            <input type='text' className='form-control' value='isProtectionCase' disabled />
            <span style={{ width: '15rem' }} className='input-group-text'>Boolean</span>
          </div>

          <div className='input-group mb-3'>
            <input type='text' className='form-control' value='hasDisability' disabled />
            <span style={{ width: '15rem' }} className='input-group-text'>Boolean</span>
          </div>

          <h6 className={'text-uppercase text-primary mb-3'}>State</h6>
          <small className='text-muted mb-3 d-block'>Uganda Assessment has no state</small>

          <button className='btn btn-primary text-white '>Fill Uganda Assessment</button>


        </div>
      </div>
    </div>

  </div>;
};
