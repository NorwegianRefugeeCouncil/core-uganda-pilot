import * as React from 'react';
import { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { getFormDefinition, selectFormDefinitionById } from '../../reducers/formdefinitions';
import { RootState } from '../../store/store';
import { FormDefinition, FormDefinitionVersion } from '@core/api-client';
import { useParams } from 'react-router-dom';
import { Builder } from '@core/formbuilder';

export type Params = {
  id: string
}

export type Props = {
  formDefinition: FormDefinition
  currentVersionName?: string
  setCurrentVersionName: (name: string) => void
  currentVersion?: FormDefinitionVersion
}

type FormDefinitionContainerProps = {}

const FormDefinitionContainer: React.FC<FormDefinitionContainerProps> = (props, context) => {

    const dispatch = useDispatch();

    // holds the id of the current formDefinition provided by BrowserRouter
    const { id } = useParams<Params>();

    // retrieves the formDefinition by name from the state
    const formDefinition = useSelector<RootState, FormDefinition | undefined>(state => selectFormDefinitionById(state, id));

    // load form definitions on load
    useEffect(() => {
      if (id) {
        dispatch(getFormDefinition(id));
      }
    }, [id]);

    return <Builder formDefinition={formDefinition} />;
  }
;

export default FormDefinitionContainer;
