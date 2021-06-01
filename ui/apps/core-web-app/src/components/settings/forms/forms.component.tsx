import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { listFormDefinitions, selectAllFormDefinitions } from '../../../reducers/formdefinitions';
import { useEffect } from 'react';
import { ListGroup } from '@core/ui-toolkit';
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
        <ListGroup>
          {allFormDefinitions.map(f => {
            return <Link className={'list-group-item'} to={`${url}/${f.metadata.name}`}>
              {f.metadata.name}
            </Link>;
          })}
        </ListGroup>
      </Route>
      <Route path={`${path}/:id`}>
        <div className={'px-2 mb-5'}>
          <FormDefinitionContainer />
        </div>
      </Route>
    </Switch>
  );
};
