import * as React from 'react';
import { ListGroup, ListGroupItem } from '@core/ui-toolkit';
import { Route, Switch } from 'react-router-dom';
import { FormsComponent } from './forms/forms.component';
import { Link, useRouteMatch } from 'react-router-dom';

export const SettingsComponent: React.FC = props => {


  let { path, url } = useRouteMatch();

  return <Switch>
    <Route exact path={path}>
      <ListGroup>
        <ListGroupItem>
          <Link className={'nav-item'} to={`${url}/forms`}>Forms</Link>
        </ListGroupItem>
      </ListGroup>
    </Route>
    <Route path={`${path}/forms`}>
      <FormsComponent />
    </Route>
  </Switch>;

};
