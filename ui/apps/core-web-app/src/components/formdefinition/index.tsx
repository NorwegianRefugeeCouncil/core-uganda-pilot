import * as React from 'react';
import { useEffect, useMemo, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { getFormDefinition, selectFormDefinitionById } from '../../reducers/formdefinitions';
import { RootState } from '../../store/store';
import { FormDefinition, FormDefinitionVersion } from '@core/api-client';
import { Link, useParams } from 'react-router-dom';
import { Card, Container } from '@core/ui-toolkit';
import { BuilderContainer, replaceField } from '@core/formbuilder';
import { CardBody } from '@core/ui-toolkit';
import { CardHeader } from '@core/ui-toolkit';

export type Params = {
  id: string
}

export type Props = {
  formDefinition: FormDefinition
  currentVersionName?: string
  setCurrentVersionName: (name: string) => void
  currentVersion?: FormDefinitionVersion
}

const FormDefinitionComponent: React.FC<Props>
  = ({
       formDefinition,
       currentVersionName,
       setCurrentVersionName,
       currentVersion
     }, context) => {
  return (
    <>
      <Card>
        <CardBody>

          <div className={'input-group mb-3'}>
            <span style={{ width: '5rem' }} className='input-group-text'>Group</span>
            <input disabled type='text' className='form-control bg-light' value={formDefinition.spec.group} />
          </div>

          <div className={'input-group mb-3'}>
            <span style={{ width: '5rem' }} className='input-group-text'>Kind</span>
            <input disabled type='text' className='form-control bg-light' value={formDefinition.spec.names.kind} />
          </div>

          <div className={'input-group mb-3'}>
            <span style={{ width: '5rem' }} className='input-group-text'>Plural</span>
            <input disabled type='text' className='form-control bg-light' value={formDefinition.spec.names.plural} />
          </div>

          <div className={'input-group mb-3'}>
            <span style={{ width: '5rem' }} className='input-group-text'>Singular</span>
            <input disabled type='text' className='form-control bg-light' value={formDefinition.spec.names.singular} />
          </div>

          <hr />

          <div className='input-group mb-3'>
            <span className='input-group-text' id='basic-addon1'>Version</span>
            <select
              className={'form-select'}
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

          <div className={'mt-2'}>
            {renderVersion(currentVersion)}
          </div>
        </CardBody>
      </Card>
    </>);
};

const renderVersion = (version?: FormDefinitionVersion) => {
  if (!version) {
    return <div />;
  }
  return <div>

    <div className='form-check'>
      <input disabled className='form-check-input' type='checkbox' value='' id='flexCheckDefault'
             checked={version.served} />
      <label className='form-check-label text-dark' htmlFor='flexCheckDefault'>
        Served
      </label>
    </div>

    <div className='form-check'>
      <input disabled className='form-check-input' type='checkbox' value='' id='flexCheckDefault'
             checked={version.storage} />
      <label className='form-check-label text-dark' htmlFor='flexCheckDefault'>
        Storage
      </label>
    </div>

    <div className={'mt-3'}>
      <BuilderContainer />
    </div>
  </div>;
};


type FormDefinitionContainerProps = {}

const FormDefinitionContainer: React.FC<FormDefinitionContainerProps> = (props, context) => {

    // holds the id of the current formDefinition provided by BrowserRouter
    const { id } = useParams<Params>();

    // retrieves the formDefinition by name from the state
    const formDefinition = useSelector<RootState, FormDefinition | undefined>(state => selectFormDefinitionById(state, id));

    // holds the different versions available in the formDefinition
    const versions = useMemo(() => formDefinition?.spec?.versions, [formDefinition]);

    // holds the currently selected version name
    const [currentVersionName, setCurrentVersionName] = useState('');

    // holds the currently selected version
    const currentVersion = useMemo(() => {
      console.log(versions, currentVersionName);
      if (versions && currentVersionName) {
        return versions.find(v => v.name === currentVersionName);
      }
      return undefined;
    }, [versions, currentVersionName]);

    const dispatch = useDispatch();


    useEffect(() => {
      if (currentVersion) {
        dispatch(replaceField({
          path: '/root',
          field: currentVersion.schema.formSchema.root
        }));
      }
    }, [currentVersion]);


    // load form definitions on load
    useEffect(() => {
      if (id) {
        dispatch(getFormDefinition(id));
      }
    }, [id]);

    // select first version on load
    useEffect(() => {
      if (formDefinition?.spec?.versions?.length) {
        setCurrentVersionName(formDefinition.spec.versions[0].name);
      }
    }, [formDefinition]);



    if (formDefinition) {
      return <>

        <div className={'d-flex justify-content-end sticky-top bg-white shadow-sm p-2'}>
          <div className={'container'}>
            <div className={'row'}>
              <div className={'col text-end'}>
                <button className={'btn btn-primary text-light'}>Save</button>
                <button className={'btn btn-outline-primary ms-2'}>Cancel</button>
              </div>
            </div>
          </div>
        </div>

        <div className={'container px-0 mt-3'}>


          <div className='row'>
            <div className='col'>
              <FormDefinitionComponent
                formDefinition={formDefinition}
                currentVersion={currentVersion}
                currentVersionName={currentVersionName}
                setCurrentVersionName={setCurrentVersionName}
              />
            </div>
          </div>
        </div>
      </>;
    } else {
      return <div>Not found</div>;
    }
  }
;

export default FormDefinitionContainer;
