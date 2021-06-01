import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { listFormDefinitions, selectAllFormDefinitions } from '../../../reducers/formdefinitions';
import { useEffect } from 'react';
import { useRouteMatch, Link, Switch, Route } from 'react-router-dom';
import FormDefinitionContainer from '../../formdefinition';

export const FormsComponent: React.FC = props => {
  const dispatch = useDispatch();
  const allFormDefinitions = useSelector(selectAllFormDefinitions);

  useEffect(() => {
    dispatch(listFormDefinitions());
  }, []);

  let { path, url } = useRouteMatch();

  return (
    <Switch>
      <Route exact path={`${path}`}>

        {allFormDefinitions.map(f => {
          return <div className={'d-flex flex-row align-items-center  border-bottom'}>
            <div className={'p-3'}>
              <i className='bi bi-pencil' />
            </div>
            <Link to={`${url}/${f.metadata.name}`}>{f.metadata.name}</Link>
          </div>;
        })}

      </Route>
      <Route path={`${path}/:id`}>
        <div className={'px-2 mb-5'}>
          <FormDefinitionContainer />
        </div>
      </Route>
    </Switch>
  );
};
