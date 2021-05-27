import * as React from 'react';
import { useEffect, useMemo, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { getFormDefinition, selectFormDefinitionById } from '../../reducers/formdefinitions';
import { RootState } from '../../store/store';
import { FormDefinition, FormDefinitionVersion } from '@core/api-client';
import { useParams } from 'react-router-dom';
import { Card, Container } from '@core/ui-toolkit';
import { Field, FieldType } from '@core/formbuilder';

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
  return (<Card>
    <Card.Header>
      Form : {formDefinition.metadata.name}
    </Card.Header>
    <Card.Body>
      <code>
        <pre>
          Kind: {formDefinition.spec.names.kind}<br />
          Group: {formDefinition.spec.group}<br />
          Plural: {formDefinition.spec.names.plural}<br />
          Singular: {formDefinition.spec.names.singular}<br />
        </pre>
      </code>
      <select
        value={currentVersionName}
        onChange={ev => setCurrentVersionName(ev.target.value)}>
        {formDefinition.spec.versions.map(v => {
          return <option
            key={v.name}
            value={v.name}
          >{v.name}</option>;
        })}
      </select>
      <div className={'mt-2'}>
        {renderVersion(currentVersion)}
      </div>
    </Card.Body>
  </Card>);
};

const renderVersion = (version?: FormDefinitionVersion) => {
  if (!version) {
    return <div />;
  }
  return <div>
    <h5>
      Served: {version.served ? 'yes' : 'no'}
    </h5>
    <h5>
      Stored: {version.storage ? 'yes' : 'no'}
    </h5>
    {version.schema.formSchema.root.children?.map(c => {
      return <Field key={c.key} type={c.type as FieldType} options={c}/>;
    })}
    Version: {JSON.stringify(version)}
  </div>;
};


const FormDefinitionContainer: React.FC = (props, context) => {

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
    return <Container size='fluid'>
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
    </Container>;
  } else {
    return <div>Not found</div>;
  }
};

export default FormDefinitionContainer;
