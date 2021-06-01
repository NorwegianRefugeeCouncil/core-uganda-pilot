import * as React from 'react';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { setFormDefinition, StateSlice } from '../../reducers';
import { useDispatch, useSelector } from 'react-redux';
import { FormDefinition } from '@core/api-client';
import { Card, CardBody } from '@core/ui-toolkit';
import { VersionBuilder } from './version.builder';


type BuilderProps = {
  formDefinition: FormDefinition
}

export const Builder: React.FC<BuilderProps> = (props, context) => {

  const givenFormDefinition = props.formDefinition;

  if (!givenFormDefinition) {
    return <div />;
  }

  const dispatch = useDispatch();

  const formDefinition = useSelector<StateSlice, FormDefinition>(state => state.formBuilder.formDefinition);

  const [currentVersionName, setCurrentVersionName] = useState(formDefinition?.spec?.versions[0].name);

  useEffect(() => {
    dispatch(setFormDefinition({ formDefinition: givenFormDefinition }));
    setCurrentVersionName(givenFormDefinition.spec.versions[0].name);
  }, [givenFormDefinition]);

  const currentVersion = useMemo(() => {
    return formDefinition?.spec?.versions?.find(v => v.name === currentVersionName);
  }, [formDefinition, currentVersionName]);

  useCallback(args => {
    console.log(args);
  }, [currentVersionName]);

  return <>


    <Card className={'shadow'}>
      <CardBody>


        <div className={'input-group mb-3'}>
          <span style={{ width: '5rem' }} className='input-group-text'>Group</span>
          <input disabled type='text' className='form-control bg-light font-monospace'
                 value={formDefinition.spec.group} />
        </div>

        <div className={'input-group mb-3'}>
          <span style={{ width: '5rem' }} className='input-group-text'>Kind</span>
          <input disabled type='text' className='form-control bg-light font-monospace'
                 value={formDefinition.spec.names.kind} />
        </div>

        <div className={'input-group mb-3'}>
          <span style={{ width: '5rem' }} className='input-group-text'>Plural</span>
          <input disabled type='text' className='form-control bg-light font-monospace'
                 value={formDefinition.spec.names.plural} />
        </div>

        <div className={'input-group mb-3'}>
          <span style={{ width: '5rem' }} className='input-group-text'>Singular</span>
          <input disabled type='text' className='form-control bg-light  font-monospace'
                 value={formDefinition.spec.names.singular} />
        </div>

        <label className={'form-label'}>Select version</label>
        <div className='input-group mb-3'>
          <span className='input-group-text' id='basic-addon1'>Version</span>
          <select
            className={'form-select  font-monospace'}
            value={currentVersionName}
            onChange={ev => setCurrentVersionName(ev.target.value)}>
            {formDefinition.spec.versions.map(v => {
              return <option
                key={v.name}
                value={v.name}
              >{v.name}</option>;
            })}
          </select>
        </div>


        {/* Is the version served ? */}
        <div className='form-check'>
          <input disabled
                 className='form-check-input'
                 type='checkbox'
                 value=''
                 checked={currentVersion.served}
          />

          <label className='form-check-label text-dark'>
            Served
          </label>

        </div>

        {/* Is it a storage version ? */}
        <div className='form-check'>

          <input disabled
                 className='form-check-input'
                 type='checkbox'
                 value=''
                 checked={currentVersion.storage} />

          <label className='form-check-label text-dark'>
            Storage
          </label>

        </div>

      </CardBody>
    </Card>

    <div className={'mt-2'}>
      {VersionBuilder({
        path: 'spec.versions[0]',
        version: currentVersion
      })}
    </div>


  </>;
};


